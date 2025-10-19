package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abhijeet/ci-pipeline-dashboard/internal/db"
)

// GitHubWebhookPayload matches the workflow_run payload from GitHub
type GitHubWebhookPayload struct {
	Workflow struct {
		Name string `json:"name"`
	} `json:"workflow"`
	Repository struct {
		Name string `json:"name"`
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

// WebhookHandler receives GitHub webhook requests
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify GitHub signature
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	signature := r.Header.Get("X-Hub-Signature-256")
	if secret == "" || signature == "" || !verifySignature(secret, body, signature) {
		log.Println("❌ Invalid or missing GitHub signature")
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Debug: log raw payload
	log.Println("Webhook payload raw:", string(body))

	// Parse JSON payload
	var payload GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("❌ Failed to parse JSON:", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate workflow run ID
	if payload.WorkflowRun.ID == 0 {
		log.Println("❌ WorkflowRun.ID is 0, skipping insert")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Webhook received, but no valid workflow ID"))
		return
	}

	// Parse timestamps with fallback to current time
	var createdAt, updatedAt time.Time
	if payload.WorkflowRun.CreatedAt != "" {
		t, err := time.Parse(time.RFC3339, payload.WorkflowRun.CreatedAt)
		if err != nil {
			log.Println("⚠️ Invalid created_at, using now:", err)
			createdAt = time.Now()
		} else {
			createdAt = t
		}
	} else {
		createdAt = time.Now()
	}

	if payload.WorkflowRun.UpdatedAt != "" {
		t, err := time.Parse(time.RFC3339, payload.WorkflowRun.UpdatedAt)
		if err != nil {
			log.Println("⚠️ Invalid updated_at, using now:", err)
			updatedAt = time.Now()
		} else {
			updatedAt = t
		}
	} else {
		updatedAt = time.Now()
	}

	// Compute duration in seconds
	durationSec := int(updatedAt.Sub(createdAt).Seconds())

	// Map payload to Pipeline struct
	pipeline := db.Pipeline{
		ID:           payload.WorkflowRun.ID,
		WorkflowName: payload.Workflow.Name,
		RepoName:     payload.Repository.Name,
		Status:       payload.WorkflowRun.Status,
		Conclusion:   payload.WorkflowRun.Conclusion,
		Actor:        payload.Sender.Login,
		RunNumber:    payload.WorkflowRun.RunNumber,
		HeadBranch:   payload.WorkflowRun.HeadBranch,
		HeadSha:      payload.WorkflowRun.HeadSha,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		DurationSec:  durationSec,
	}

	// Insert or update pipeline in Supabase
	if err := db.UpsertPipeline(pipeline); err != nil {
		http.Error(w, "Failed to insert/update pipeline", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Pipeline processed successfully:", pipeline.ID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}
