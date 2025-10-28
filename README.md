# CI Pipeline Dashboard

A full-stack web application that visualizes and tracks GitHub Actions CI/CD workflows in real time.

---

(Frontend) Accessible live Dashboard:

  https://ci-pipeline-dashboard.vercel.app/

(Backend) Accessible live:

  For list view : https://ci-pipeline-dashboard.onrender.com/api/pipelines
  
  For detail view : https://ci-pipeline-dashboard.onrender.com/api/pipelines/details?id=<pipeline-id>   
                   (e.g.  https://ci-pipeline-dashboard.onrender.com/api/pipelines/details?id=78d72bc3-e683-45fb-84c0-c645ae29ff90)

## Overview

The CI Pipeline Dashboard connects directly to GitHub via webhooks.  
When a workflow runs (triggered by a push or pull request), GitHub sends data to the "Go backend", which processes and stores it in "Supabase".  
The "React frontend" then fetches this data to display a clear visualization of the pipelines — showing commits, authors, job steps, and statuses.

## Core features:

- Real-time pipeline tracking from GitHub Actions  
- Detailed view of commits, authors, and workflow steps  
- Supabase database integration for structured pipeline data  
- Responsive, modern dashboard UI built with React + TypeScript  


## Architecture

 GitHub Actions → Webhook → Go Backend → Supabase → React Dashboard


## High-Level Design

1. 'GitHub Actions' triggers a webhook on workflow events.  
2. 'Go Backend' receives the webhook, parses payloads, and upserts data into Supabase.  
3. 'Supabase' stores structured data for pipelines, commits, and job steps.  
4. 'React Frontend' queries the backend REST API to display pipeline runs and details.


## Tech Stack

    Frontend : React + Vite + TypeScript 
    Backend  : Go (Golang) 
    Database : Supabase (PostgreSQL)
    Deployment : Render (Backend), Vercel (Frontend)
    Version Control :  GitHub 

## Database Schema

### Tables & Relationships:

Entity Relationship Diagram 

    git_info (id) ----- 1:N --> pipelines (git_info_id) ----- 1:N --> jobs_steps (pipelines_id)
                                  
                 

### Table Details


 1. git_info

| Column          | Type      |      Description               

- id              | UUID      |    Primary key 
- repo_name       | text      |    Repository name
- commit_sha      | text      |    Commit SHA for the workflow 
- commit_message  | text      |    Commit message 
- author_name     | text      |    Commit author 
- branch          | text      |    Branch name 
- commited_at     | timestampz|    commit time

 2. pipelines

| Column        | Type       | Description 

- id            | UUID       | Primary key 
- git_info_id   | UUID       | References git_info.id
- run_id        | int8       | run id of the workflow
- workflow_name | text       | Name of the workflow 
- status        | text       | completed or not
- conclusion    | text       | success or failure 
- started_at    | timestampz | Trigger time of the workflow 
- completed_at  | timestampz | Last updated time 

 3. jobs_steps

| Column        | Type         | Description 

- id            | UUID         | Primary key 
- pipeline_id   | UUID         | References pipelines.id 
- name          | text         | Job name 
- status        | text         | job status (queued/in_progress/completed) |
- started_at    | timestampz   | Start time 
- completed_at  | timestampz   | End time 
- job_id        | int8         | id of the job 
- duration_sec  | int4         | duration of the job
- attempt       | int4         | number of times the job id ran

---

## API Endpoints

### Get All Pipelines

   1. GET /api/pipelines

      Response Example:

      json
[
  {
    "id": "b2d3d24e-849b-4672-8894-b8ec15711179",
    "runId": 18746686599,
    "status": "completed",
    "branch": "main",
    "commitSha": "c8374b57e47e356eb01d044a72d1fc05233a8994",
    "startedAt": "2025-10-23T11:21:01Z",
    "duration": "14s"
  },
  {
    "id": "b3c5dcbf-1180-401c-975c-9f5128305793",
    "runId": 18753685849,
    "status": "completed",
    "branch": "main",
    "commitSha": "c8374b57e47e356eb01d044a72d1fc05233a8994",
    "startedAt": "2025-10-23T15:32:57Z",
    "duration": "12s"
  }
]


    2. GET GET /api/pipelines/details?id=<pipeline-id>

       Response Example:

       json
 {
  "author": {
    	"email": "",
    	"name": "abhijeetkr0123"
  },
  "branch": "main",
  "commitMessage": "",
  "commitSha": "c8374b57e47e356eb01d044a72d1fc05233a8994",
  "duration": "13s",
  "jobs": [
    {
      "id": "53691507943",
      "name": "build",
      "status": "queued",
      "startedAt": "2025-10-26T14:01:45Z",
      "completedAt": "0001-01-01T00:00:00Z",
      "duration": "0s",
      "steps": [
        {
          "name": "build",
          "status": "queued",
          "duration": "0s"
        },
        {
          "name": "build",
          "status": "in_progress",
          "duration": "0s"
        },
        {
          "name": "build",
          "status": "completed",
          "duration": "4s"
        }
      ]
    }
  ],
  "repository_url": "Todo_app",
  "runId": 18818980724,
  "startedAt": "2025-10-26T14:01:45Z",
  "status": "completed"
}

    3. POST /webhook

     Purpose: Receives GitHub Action workflow payloads.



## Setup Instructions

### Backend (Go)

**Prerequisites
- Go 1.22+
- Supabase project created  
- '.env' file with your credentials

**.env example:

- SUPABASE_URL=<your_supabase_url>
- SUPABASE_KEY=<your_supabase_service_role_key>
- PORT=8080

**Run locally:

-  go run cmd/server/main.go

   Backend starts at `http://localhost:8080`


### Frontend (React + Vite)

**Prerequisites
- Node.js 18+
- Backend running (locally or Render)

**.env example:

-  VITE_API_BASE_URL = http://localhost:8080

### Webhook Secret (Security)

To ensure only GitHub can trigger your backend, the '/webhook' endpoint verifies an HMAC SHA-256 signature sent with each request.

When creating the webhook in your GitHub repository settings:

1. Go to Settings -> Webhooks -> Add Webhook
2. Set the Payload URL to your backend endpoint:  
    https://<your-backend-service>.onrender.com/webhook
3. Set Content type to 'application/json'.
4. Add a secret — this should match the 'WEBHOOK_SECRET' in your .env file


**Run locally:

-  npm install
-  npm run dev

   App runs at `http://localhost:5173


## Deployment

"Backend":  

- Deployed on [Render](https://render.com/)  
- Receives GitHub webhook payloads at:
  
  https://ci-pipeline-dashboard.onrender.com/webhook
  (Webhook events are verified using an HMAC-SHA256 signature and secret key.)

"Frontend":  

- Deployed on [Vercel](https://vercel.com/)  
- Make sure to update the '.env' variable:
  
  VITE_API_BASE_URL = https://ci-pipeline-dashboard.onrender.com

- Accessible live Dashboard:

  https://ci-pipeline-dashboard.vercel.app/
  (The frontend automatically connects to the backend API deployed on Render.)

"Database":  

- Hosted on 'Supabase', with public and service role keys configured securely.


## Tradeoffs & Future Improvements

-  Backend Language ->  Chose Go for speed & concurrency , Slightly more setup vs Node , can add middleware for auth/logging
-  Database  ->  Supabase (managed Postgres) , Dependent on external API , can be migrated to self-hosted Postgres if scaling
-  Frontend  ->  Vite + React , Fast build but needs manual proxy config , Integrate CI/CD preview builds 
-  Integration ->  GitHub webhooks , Only tracks specific repo actions , Expand to multiple repos per user
-  Testing ->  Deferred for MVP , No coverage yet , Add Go tests for webhook + Supabase insertions 

---

## Author

- Abhijeet Kumar  
  [GitHub: abhijeetkr0123](https://github.com/abhijeetkr0123)
