package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	Actor      struct {
		Login string `json:"login"`
	} `json:"actor"`
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

	// Parse JSON payload
	var payload GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("❌ Failed to parse JSON:", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Map payload to Pipeline struct
	pipeline := db.Pipeline{
		WorkflowName: payload.Workflow.Name,
		RepoName:     payload.Repository.Name,
		Status:       payload.Status,
		Conclusion:   payload.Conclusion,
		Actor:        payload.Actor.Login,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	// Insert into Supabase
	if err := db.InsertPipeline(pipeline); err != nil {
		http.Error(w, "Failed to insert pipeline", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}
