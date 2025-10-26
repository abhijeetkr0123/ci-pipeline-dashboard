package db

import (
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	supabase "github.com/nedpals/supabase-go"
)

// Client is the global Supabase client
var Client *supabase.Client

// InitDB initializes the Supabase client
func InitDB() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("❌ Missing SUPABASE_URL or SUPABASE_KEY in environment variables")
	}

	Client = supabase.CreateClient(supabaseURL, supabaseKey)
	log.Println("✅ Connected to Supabase successfully")
}

// UpsertGitInfo inserts or fetches GitInfo for a commit
func UpsertGitInfo(g GitInfo) (string, error) {
	// Check if commit already exists
	existing := []GitInfo{}
	err := Client.DB.
		From("git_info").
		Select("*").
		Eq("commit_sha", g.CommitSHA).
		Execute(&existing)
	if err != nil {
		log.Println("❌ Error checking existing GitInfo:", err)
		return "", err
	}

	if len(existing) > 0 {
		return existing[0].ID, nil
	}

	// Insert new GitInfo
	id := uuid.New().String()
	g.ID = id
	err = Client.DB.
		From("git_info").
		Insert(g).
		Execute(nil)
	if err != nil {
		log.Println("❌ Error inserting GitInfo:", err)
		return "", err
	}

	log.Printf("✅ GitInfo inserted successfully: %s", id)
	return id, nil
}

// UpsertPipeline inserts a new pipeline or updates an existing one
// UpsertPipeline inserts a new pipeline or updates an existing one
func UpsertPipeline(p Pipeline) (string, error) {
	existing := []Pipeline{}

	// Convert RunID to string for Supabase Eq()
	runIDStr := strconv.FormatInt(p.RunID, 10)

	// Check if pipeline already exists
	err := Client.DB.
		From("pipelines").
		Select("*").
		Eq("run_id", runIDStr).
		Execute(&existing)
	if err != nil {
		log.Println("❌ Error checking existing pipeline:", err)
		return "", err
	}

	if len(existing) > 0 {
		// Update existing pipeline
		updateData := map[string]interface{}{
			"status":       p.Status,
			"conclusion":   p.Conclusion,
			"started_at":   p.StartedAt,
			"completed_at": p.CompletedAt,
		}

		err = Client.DB.
			From("pipelines").
			Update(updateData).
			Eq("id", existing[0].ID). // UUID
			Execute(nil)
		if err != nil {
			log.Println("❌ Error updating pipeline:", err)
			return "", err
		}

		log.Printf("✅ Pipeline updated successfully: %s", existing[0].ID)
		return existing[0].ID, nil
	}

	// Insert new pipeline
	p.ID = uuid.New().String()
	err = Client.DB.
		From("pipelines").
		Insert(p).
		Execute(nil)
	if err != nil {
		log.Println("❌ Error inserting pipeline:", err)
		return "", err
	}

	log.Printf("✅ Pipeline inserted successfully: %s", p.ID)
	return p.ID, nil
}

// UpsertJobSteps inserts multiple job steps for a pipeline
func UpsertJobSteps(jobs []JobStep) error {
	if len(jobs) == 0 {
		return nil
	}

	// Assign UUIDs if not set
	for i := range jobs {
		if jobs[i].ID == "" {
			jobs[i].ID = uuid.New().String()
		}
	}

	err := Client.DB.
		From("jobs_steps").
		Insert(jobs).
		Execute(nil)
	if err != nil {
		log.Println("❌ Error inserting JobSteps:", err)
		return err
	}

	log.Printf("✅ Inserted %d job steps", len(jobs))
	return nil
}
