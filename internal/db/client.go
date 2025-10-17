package db

import (
	"log"
	"os"

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
