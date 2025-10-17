package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abhijeet/ci-pipeline-dashboard/internal/db"
	"github.com/abhijeet/ci-pipeline-dashboard/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Initialize Supabase client
	db.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/webhook", handlers.WebhookHandler)

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
