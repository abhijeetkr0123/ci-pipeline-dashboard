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

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (for development)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

	// Use default mux
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", handlers.WebhookHandler)
	mux.HandleFunc("/api/pipelines", handlers.GetPipelinesHandler)
	mux.HandleFunc("/api/pipelines/details", handlers.GetPipelineDetailsHandler)

	// Wrap mux with CORS middleware
	handler := corsMiddleware(mux)

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
