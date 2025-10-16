package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/abhijeet/ci-pipeline-dashboard/internal/handlers"
 // adjust module path
)

func main() {
    r := chi.NewRouter()

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Server is running!")
    })

    r.Post("/webhook", handlers.WebhookHandler) // your POST endpoint

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
