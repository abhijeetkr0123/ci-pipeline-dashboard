package db

import "time"

// Pipeline represents a workflow run
type Pipeline struct {
	ID          string    `json:"id"`            // UUID
	RunID       int64     `json:"run_id"`        // GitHub run ID
	Workflow    string    `json:"workflow_name"` // matches pipelines.workflow_name
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	GitInfoID   string    `json:"git_info_id"` // UUID FK to git_info.id
	CreatedAt   time.Time `json:"created_at"`
}

// GitInfo represents git-related info
type GitInfo struct {
	ID            string    `json:"id"`        // UUID
	RepoName      string    `json:"repo_name"` // matches git_info.repo_name
	CommitSHA     string    `json:"commit_sha"`
	Branch        string    `json:"branch"`
	AuthorName    string    `json:"author_name"`
	AuthorEmail   string    `json:"author_email"`
	CommitMessage string    `json:"commit_message"`
	CommittedAt   time.Time `json:"committed_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// JobStep represents a job or step in a workflow run
type JobStep struct {
	ID          string    `json:"id"`          // UUID or string if you use GitHub job ID as string
	PipelineID  string    `json:"pipeline_id"` // UUID FK to pipelines.id
	JobID       string    `json:"job_id"`      // GitHub job ID as string
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "job" or "step"
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	DurationSec int       `json:"duration_sec"`
	Attempt     int       `json:"attempt"` // increments for reruns
}
