package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/abhijeet/ci-pipeline-dashboard/internal/db"
)

// ==============================
// GET /api/pipelines (List view)
// ==============================
func GetPipelinesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var pipelines []db.Pipeline

	err := db.Client.DB.
		From("pipelines").
		Select("id,run_id,status,started_at,completed_at,git_info_id").
		Execute(&pipelines)
	if err != nil {
		log.Println("❌ Error fetching pipelines:", err)
		http.Error(w, "Failed to fetch pipelines", http.StatusInternalServerError)
		return
	}

	log.Printf("Fetched pipelines: %+v\n", pipelines)

	type PipelineListItem struct {
		RunID     int64  `json:"runId"`
		Status    string `json:"status"`
		Branch    string `json:"branch"`
		CommitSHA string `json:"commitSha"`
		StartedAt string `json:"startedAt"`
		Duration  string `json:"duration"`
	}

	var response []PipelineListItem

	for _, p := range pipelines {
		var gitList []db.GitInfo
		if p.GitInfoID != "" {
			err := db.Client.DB.
				From("git_info").
				Select("branch,commit_sha").
				Eq("id", p.GitInfoID).
				Execute(&gitList)
			if err != nil {
				log.Println("⚠️ Error fetching git_info:", err)
			}
		}

		branch, commitSHA := "", ""
		if len(gitList) > 0 {
			branch = gitList[0].Branch
			commitSHA = gitList[0].CommitSHA
		}

		duration := ""
		if !p.StartedAt.IsZero() && !p.CompletedAt.IsZero() {
			duration = p.CompletedAt.Sub(p.StartedAt).Truncate(time.Second).String()
		}

		response = append(response, PipelineListItem{
			RunID:     p.RunID,
			Status:    p.Status,
			Branch:    branch,
			CommitSHA: commitSHA,
			StartedAt: p.StartedAt.Format(time.RFC3339),
			Duration:  duration,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// =================================
// GET /api/pipelines/details?id=<uuid>
// =================================
func GetPipelineDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pipelineID := r.URL.Query().Get("id")
	if pipelineID == "" {
		http.Error(w, "Missing pipeline id", http.StatusBadRequest)
		return
	}

	var pipelineList []db.Pipeline
	err := db.Client.DB.
		From("pipelines").
		Select("*").
		Eq("run_id", pipelineID).
		Execute(&pipelineList)
	if err != nil || len(pipelineList) == 0 {
		log.Println("❌ Error fetching pipeline details:", err)
		http.Error(w, "Pipeline not found", http.StatusNotFound)
		return
	}
	pipeline := pipelineList[0]

	var gitList []db.GitInfo
	if pipeline.GitInfoID != "" {
		err = db.Client.DB.
			From("git_info").
			Select("branch,commit_sha,commit_message,author_name,author_email,committed_at,repo_name").
			Eq("id", pipeline.GitInfoID).
			Execute(&gitList)
		if err != nil {
			log.Println("⚠️ Error fetching git_info:", err)
		}
	}

	var jobs []db.JobStep
	err = db.Client.DB.
		From("jobs_steps").
		Select("id, pipeline_id, job_id::text as job_id, name, type, status, conclusion, started_at, completed_at, duration_sec, attempt").
		Eq("pipeline_id", pipeline.ID).
		Execute(&jobs)
	if err != nil {
		log.Println("⚠️ Error fetching jobs_steps:", err)
	}

	if jobs == nil {
		jobs = []db.JobStep{}
	}
	// Prepare frontend-friendly jobs structure
	type Step struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Duration string `json:"duration"`
	}

	type Job struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Status      string `json:"status"`
		StartedAt   string `json:"startedAt"`
		CompletedAt string `json:"completedAt"`
		Duration    string `json:"duration"`
		Steps       []Step `json:"steps"`
	}

	jobMap := map[string][]Step{}
	for _, j := range jobs {
		step := Step{
			Name:     j.Name,
			Status:   j.Status,
			Duration: time.Duration(j.DurationSec * int(time.Second)).String(),
		}
		jobMap[j.JobID] = append(jobMap[j.JobID], step)
	}

	var jobResponse []Job
	for _, j := range jobs {
		exists := false
		for _, jr := range jobResponse {
			if jr.ID == j.JobID {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		jobResponse = append(jobResponse, Job{
			ID:          j.JobID,
			Name:        j.Name,
			Status:      j.Status,
			StartedAt:   j.StartedAt.Format(time.RFC3339),
			CompletedAt: j.CompletedAt.Format(time.RFC3339),
			Duration:    time.Duration(j.DurationSec * int(time.Second)).String(),
			Steps:       jobMap[j.JobID],
		})
	}

	// Combine response
	response := map[string]interface{}{
		"pipeline": map[string]interface{}{
			"runId":  pipeline.RunID,
			"status": pipeline.Status,
			"branch": func() string {
				if len(gitList) > 0 {
					return gitList[0].Branch
				}
				return ""
			}(),
			"commitSha": func() string {
				if len(gitList) > 0 {
					return gitList[0].CommitSHA
				}
				return ""
			}(),
			"startedAt": pipeline.StartedAt.Format(time.RFC3339),
			"duration": func() string {
				if !pipeline.StartedAt.IsZero() && !pipeline.CompletedAt.IsZero() {
					return pipeline.CompletedAt.Sub(pipeline.StartedAt).Truncate(time.Second).String()
				}
				return ""
			}(),
		},
		"git_info": func() map[string]string {
			if len(gitList) > 0 {
				g := gitList[0]
				return map[string]string{
					"commitMessage": g.CommitMessage,
					"authorName":    g.AuthorName,
					"authorEmail":   g.AuthorEmail,
					"committedAt":   g.CommittedAt.Format(time.RFC3339),
					"repoName":      g.RepoName,
				}
			}
			return nil
		}(),
		"jobs": jobResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
