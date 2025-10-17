package db

import (
	"log"
)

// Pipeline struct matches your Supabase table columns
type Pipeline struct {
	WorkflowName string `json:"workflow_name"`
	RepoName     string `json:"repo_name"`
	Status       string `json:"status"`
	Conclusion   string `json:"conclusion"`
	Actor        string `json:"actor"`
	CreatedAt    string `json:"created_at"`
}

// InsertPipeline inserts a pipeline into Supabase
func InsertPipeline(p Pipeline) error {
	// Wrap pipeline in a slice for Insert
	data := []Pipeline{p}

	// Use a variable to capture returned rows (can ignore if you don't need them)
	var result []Pipeline

	// Execute insert
	err := Client.DB.From("pipelines").Insert(data).Execute(&result)
	if err != nil {
		log.Println("❌ Error inserting pipeline:", err)
		return err
	}

	log.Println("✅ Pipeline inserted successfully")
	return nil
}
