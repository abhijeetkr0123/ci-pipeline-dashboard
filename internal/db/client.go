package db

import (
	"log"
	"os"
	"strconv"

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

// UpsertPipeline inserts a new pipeline or updates existing one
func UpsertPipeline(p Pipeline) error {
	existing := []Pipeline{}

	idStr := strconv.FormatInt(p.ID, 10)

	// Check if pipeline exists
	err := Client.DB.From("pipelines").
		Select("*").
		Eq("id", idStr).
		Execute(&existing)
	if err != nil {
		log.Println("❌ Error checking existing pipeline:", err)
		return err
	}

	if len(existing) > 0 {
		// Pipeline exists → update status and other fields
		updateData := map[string]interface{}{
			"status":       p.Status,
			"conclusion":   p.Conclusion,
			"updated_at":   p.UpdatedAt,
			"duration_sec": p.DurationSec,
		}

		err = Client.DB.From("pipelines").
			Update(updateData).
			Eq("id", idStr).
			Execute(nil)
		if err != nil {
			log.Println("❌ Error updating pipeline:", err)
			return err
		}

		log.Println("✅ Pipeline updated successfully:", p.ID)
		return nil
	}

	// Pipeline doesn't exist → insert new
	err = Client.DB.From("pipelines").Insert(p).Execute(nil)
	if err != nil {
		log.Println("❌ Error inserting pipeline:", err)
		return err
	}

	log.Println("✅ Pipeline inserted successfully:", p.ID)
	return nil
}
