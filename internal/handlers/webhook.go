package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/abhijeet/ci-pipeline-dashboard/internal/db"
)

// GitHubWebhookPayload matches the workflow_run payload from GitHub
type GitHubWebhookPayload struct {
	Workflow struct {
		Name string `json:"name"`
	} `json:"workflow"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
	WorkflowRun struct {
		ID         int64  `json:"id"`
		HeadBranch string `json:"head_branch"`
		HeadSha    string `json:"head_sha"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
		RunNumber  int    `json:"run_number"`
	} `json:"workflow_run"`
}

// verifySignature verifies GitHub webhook HMAC SHA256 signature
func verifySignature(secret string, body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// parseTimestamps safely parses GitHub timestamps; fallback to now when invalid
func parseTimestamps(created, updated string) (time.Time, time.Time) {
	var createdAt, updatedAt time.Time
	if t, err := time.Parse(time.RFC3339, created); err == nil {
		createdAt = t
	} else {
		createdAt = time.Now()
	}
	if t, err := time.Parse(time.RFC3339, updated); err == nil {
		updatedAt = t
	} else {
		updatedAt = time.Now()
	}
	return createdAt, updatedAt
}

// fetchJobsWithAttempts fetches jobs for the workflow run and handles rerun attempts
func fetchJobsWithAttempts(owner, repo, runID, pipelineID string) ([]db.JobStep, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("missing GITHUB_TOKEN")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%s/jobs", owner, repo, runID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyB, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API error status %d: %s", resp.StatusCode, string(bodyB))
	}

	var parsed struct {
		Jobs []struct {
			ID          int64  `json:"id"`
			Name        string `json:"name"`
			Status      string `json:"status"`
			Conclusion  string `json:"conclusion"`
			StartedAt   string `json:"started_at"`
			CompletedAt string `json:"completed_at"`
		} `json:"jobs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode jobs json: %w", err)
	}

	var out []db.JobStep
	for _, j := range parsed.Jobs {
		var startedAt, completedAt time.Time
		if t, err := time.Parse(time.RFC3339, j.StartedAt); err == nil {
			startedAt = t
		}
		if t, err := time.Parse(time.RFC3339, j.CompletedAt); err == nil {
			completedAt = t
		}
		duration := 0
		if !startedAt.IsZero() && !completedAt.IsZero() {
			duration = int(completedAt.Sub(startedAt).Seconds())
		}

		// --- convert job ID to string for DB ---
		jobIDStr := strconv.FormatInt(j.ID, 10)

		// --- check for previous attempts ---
		existingJobs := []db.JobStep{}
		_ = db.Client.DB.From("jobs_steps").
			Select("*").
			Eq("pipeline_id", pipelineID).
			Eq("job_id", jobIDStr).
			Execute(&existingJobs)

		attempt := 1
		if len(existingJobs) > 0 {
			attempt = existingJobs[0].Attempt + 1
		}

		out = append(out, db.JobStep{
			JobID:       jobIDStr, // string now
			Name:        j.Name,
			Type:        "job",
			Status:      j.Status,
			Conclusion:  j.Conclusion,
			StartedAt:   startedAt,
			CompletedAt: completedAt,
			DurationSec: duration,
			PipelineID:  pipelineID,
			Attempt:     attempt,
		})
	}

	return out, nil
}

// WebhookHandler handles GitHub workflow_run webhook events
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	// verify signature
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	signature := r.Header.Get("X-Hub-Signature-256")
	if secret == "" || signature == "" || !verifySignature(secret, body, signature) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	var payload GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}

	if payload.WorkflowRun.ID == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no workflow run id"))
		return
	}

	createdAt, updatedAt := parseTimestamps(payload.WorkflowRun.CreatedAt, payload.WorkflowRun.UpdatedAt)

	// --- upsert git info first ---
	git := db.GitInfo{
		RepoName:   payload.Repository.Name,
		Branch:     payload.WorkflowRun.HeadBranch,
		CommitSHA:  payload.WorkflowRun.HeadSha,
		AuthorName: payload.Sender.Login,
	}

	gitID, err := db.UpsertGitInfo(git)
	if err != nil {
		log.Println("❌ UpsertGitInfo warning:", err)
	}

	// --- upsert pipeline ---
	pipeline := db.Pipeline{
		Workflow:    payload.Workflow.Name,
		RunID:       payload.WorkflowRun.ID,
		Status:      payload.WorkflowRun.Status,
		Conclusion:  payload.WorkflowRun.Conclusion,
		StartedAt:   createdAt,
		CompletedAt: updatedAt,
		GitInfoID:   gitID,
	}

	pipelineID, err := db.UpsertPipeline(pipeline)
	if err != nil {
		log.Println("❌ UpsertPipeline error:", err)
		http.Error(w, "failed to upsert pipeline", http.StatusInternalServerError)
		return
	}
	log.Printf("✅ Pipeline upserted successfully: %s", pipelineID)

	// --- fetch jobs with attempt handling ---
	jobs, err := fetchJobsWithAttempts(
		payload.Repository.Owner.Login,
		payload.Repository.Name,
		strconv.FormatInt(payload.WorkflowRun.ID, 10),
		pipelineID,
	)
	if err != nil {
		log.Println("⚠️ fetchJobs warning:", err)
	} else {
		if err := db.UpsertJobSteps(jobs); err != nil {
			log.Println("⚠️ UpsertJobSteps warning:", err)
		} else {
			log.Printf("✅ Upserted %d job steps for pipeline %s", len(jobs), pipelineID)
		}
	}

	log.Printf("✅ Processed workflow run %d -> pipeline id %s", payload.WorkflowRun.ID, pipelineID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("webhook processed"))
}
