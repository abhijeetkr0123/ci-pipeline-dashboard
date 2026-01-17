The repository is a full-stack application (Go backend + React frontend) for visualizing GitHub Actions pipelines.

**Architecture:**
- **Backend (Go):** Receives GitHub webhooks, stores data in Supabase, and serves a REST API.
- **Frontend (React/Vite):** Consumes the backend API to display a dashboard of pipeline runs.
- **Database:** Supabase (PostgreSQL).

**Key Findings:**
- Clean separation of concerns (handlers, db, services, components).
- Frontend uses Material UI.
- Backend uses `net/http` and `supabase-go`.
- **Potential Issue:** `frontend/src/services/pipelineService.ts` hardcodes the production URL, potentially ignoring the `VITE_API_BASE_URL` environment variable and the Vite proxy configuration for local development.

**Next Steps:**
- Be aware of the hardcoded URL if debugging frontend-backend connection issues.
- Follow the existing pattern of using `internal` packages in Go and `components`/`services` in React.