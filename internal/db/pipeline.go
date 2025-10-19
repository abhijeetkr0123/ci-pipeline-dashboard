package db

import (
	"time"
)

// Pipeline struct matches your Supabase table columns
type Pipeline struct {
	ID           int64     `json:"id"`
	WorkflowName string    `json:"workflow_name"`
	RepoName     string    `json:"repo_name"`
	Actor        string    `json:"actor"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RunNumber    int       `json:"run_number"`
	HeadBranch   string    `json:"head_branch"`
	HeadSha      string    `json:"head_sha"`
	DurationSec  int       `json:"duration_sec"`
}

// InsertPipeline inserts a pipeline into Supabase
