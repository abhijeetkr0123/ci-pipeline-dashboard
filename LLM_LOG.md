1. Key Prompts :

- Help me design the backend for a CI Pipeline Dashboard in Go using Supabase and GitHub webhooks.

- Explain how to structure Go files like main.go, client.go, webhook.go, and pipeline.go for this project.

- Help me connect my Go backend deployed on Render to Supabase.

- Guide me to set up the frontend detail view using React and TypeScript to fetch data from /api/pipelines/?id=<pipeline-id>.

- Help me fix the Grid layout inside pipelineDetail.tsx that’s throwing JSX element errors.

- How to deploy my React frontend (Vite + TypeScript) on Vercel with the backend already on Render?

- My GitHub Actions CI/CD is failing with ‘go.mod not found’ — how do I fix this when go.mod is in /backend?


2. Modifications

Examples where I accepted, rejected, or modified LLM suggestions:

Accepted:
LLM suggested separating backend Go files (config.go, client.go, webhook.go) for better modularity.
- This structure improved clarity and debugging.

Modified:
LLM recommended using Material UI Grid for layout in the React frontend.
- I modified it to use Box for more consistent spacing and fewer type issues.

Rejected:
Initially, LLM suggested the schema to make parent table pipelines and its child table git_info.
- I rejected this after discovering that one commit message can have multiple workflow runs so i made parent table git_info and corrected relationships.


3. Debugging

Instances where LLM was wrong/unhelpful and how I fixed them:

Issue: Frontend build (npm run build) didn’t generate dist/ after switching Grid → Box.
Fix:   Ran npx vite build directly and confirmed successful output, learning how Vite bundling works independently from TypeScript builds.

Issue: API data not displaying in detail view due to missing /api/pipelines/details?id=<pipeline-id>
Fix:   Adjusted pipelineDetail.tsx and added some console to debug where the code is being broken and verified with network & console tab.


4. Learning

Concepts I researched beyond what the LLM provided:

GitHub Webhooks:
Understood how events like workflow_run and push trigger webhook payloads and how to process them to store pipeline data.
There are various types of webhooks explored from the docs that you shared.

Contract between backend and frontend:
Understood the importance of payload json through which interation happens.

Vercel Environment & Privacy:
Learned how Vercel restricts app visibility for Hobby plan users and how to share access with collaborators safely.


