# CI Pipeline Dashboard

A minimal app that receives GitHub Actions workflow run events via webhooks, stores them in a database, and displays them in a dashboard.

### Backend
- Built with Go (Chi router)
- PostgreSQL database

### Setup
```bash
go run ./cmd/server/main.go