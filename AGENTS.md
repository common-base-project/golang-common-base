# AGENTS Guide (Monorepo)

This file gives AI coding agents the minimum context to work safely and efficiently in this repository.

## Workspace Boundaries

- `backend-admin/`: Go backend (Gin + Gorm), auth via Casdoor OIDC, authorization via Casbin.
- `admin-frontend/`: Vue3 admin frontend imported from v3-admin-vite. It has its own project-specific guide: [admin-frontend/AGENTS.md](admin-frontend/AGENTS.md).
- `deploy/`: Docker compose deployment assets for local integration and Casdoor production deployment.

Use the root guide for cross-project tasks; follow subproject guides when working inside a specific folder.

## Canonical Docs

- Monorepo overview: [README.md](README.md)
- Backend details: [backend-admin/README.md](backend-admin/README.md)
- Deployment overview: [deploy/README.md](deploy/README.md)
- Casdoor production deployment: [deploy/casdoor/README.md](deploy/casdoor/README.md)

Prefer linking to these docs instead of duplicating long instructions in code comments.

## Common Commands

### Backend (`backend-admin/`)

- Compile/check: `go test ./...`
- Local build/start options are in [backend-admin/Makefile](backend-admin/Makefile)

### Frontend (`admin-frontend/`)

- Follow command set in [admin-frontend/AGENTS.md](admin-frontend/AGENTS.md)

### Deployment (`deploy/`)

- Local integration stack: `docker compose up -d --build`
- Casdoor production stack: `docker compose --env-file .env.casdoor -f docker-compose.casdoor-prod.yaml up -d`

## Project Conventions

- Keep changes scoped: only edit files relevant to the user request.
- Do not refactor unrelated modules while fixing a focused issue.
- Prefer minimal patches compatible with current style.
- For frontend tasks, preserve `admin-frontend` conventions from its local guide.

## Security and Config Safety

- Never hardcode secrets in code or docs; use env files/templates.
- Treat `conf/*.yaml` and `deploy/.env*` as environment-specific and sensitive.
- For Casdoor OIDC, `issuer` must match the externally reachable HTTPS domain used in deployment.

## Known Pitfalls

- Shell search tools may differ by environment (for example `rg` can be missing); use alternatives when needed.
- The repository contains legacy backend artifacts and newer monorepo deployment assets; prefer files under `deploy/` for current deployment workflows.
- When changing authentication/authorization behavior, validate both token verification (Casdoor/OIDC) and route authorization (Casbin), not only one side.

## Chronicle-based Improvement

Session history suggests recurrent friction around deployment/auth configuration consistency. After larger changes, run a short retrospective and refine this file over time using `/chronicle improve`.
