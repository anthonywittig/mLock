# mLock – Agent context

High-level context for AI agents working in this repo.

## What this project is

mLock is a full-stack app for managing locks, units, properties, climate controls, and users (likely vacation-rental / property management). Backend is Go on AWS Lambda; frontend is React + TypeScript.

## Layout

- **`backend/`** – Go monorepo. Shared code and config live here; Lambdas are under `backend/lambdas/`.
- **`backend/lambdas/`**
  - **`apis/`** – API Gateway Lambdas (e.g. `devices`, `units`, `users`, `properties`, `signin`, `climate-controls`, `webhooks`). Each has its own `main.go` and often subpackages (e.g. `devices/lockcodes`).
  - **`jobs/`** – Scheduled/background Lambdas (e.g. `manage-climate-controls`, `pollschedules`).
  - **`db/migrations/`** – DB migration Lambda.
  - **`helpers/`** – Lambda bootstrap (e.g. `StartAPILambda`).
  - **`shared/`** – Shared types, DynamoDB repos, Ezlo/Home Assistant/Hostaway integrations, lock engine, scheduler, SES, SQS, etc.
- **`backend/shared/`** – Repo-level shared config (e.g. `config.go`).
- **`frontend/`** – Create React App, React 18, TypeScript, React Router, React Bootstrap. Pages under `src/pages/` with route modules (e.g. `Routes.tsx`).
- **`bin/`** – Scripts: `deploy-backend.sh`, and `deploy-lambda/` (build + run deploy tool for each Lambda).
- **`ezlo/`** – Ezlo gateway/API docs and scripts (reference only).

Backend package layout follows [standard Go package layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1); imports use the `mlock` module (e.g. `mlock/lambdas/shared`, `mlock/lambdas/apis/devices/lockcodes`).

## Tech stack

- **Backend:** Go 1.21+, AWS Lambda, API Gateway, DynamoDB, SQS, SES. Mocks via `github.com/golang/mock/mockgen`.
- **Frontend:** React 18, TypeScript, React Router 6, React Bootstrap, Prettier (no semicolons per ESLint). Tests: Jest + React Testing Library.

## Conventions

- **Backend:** API handlers use `lambdas/shared` for responses (`NewAPIResponse`), context data, and auth. Auth is cookie-based (`AuthCookie`, `SetAuthCookie`). Each Lambda has `.env.example`; config is read via `shared` (e.g. `GetConfigUnsafe`).
- **Frontend:** API calls use `StandardFetch` from `pages/utils/FetchHelper.tsx` with `credentials: "include"`. Backend base URL: `REACT_APP_BACKEND_DOMAIN`. 401/403 redirect to `/sign-in`.
- **Testing:** Backend: `go test ./...` from repo root or `backend/`. Mocks live under `*_test.go` and `mocks/` in shared packages. Frontend: `npm test` in `frontend/`.

## Commands (from repo root)

| Task | Command |
|------|--------|
| Build frontend | `make build-frontend` or `./frontend/bin/build.sh` |
| Deploy backend | `make deploy-backend` or `./bin/deploy-backend.sh` |
| Deploy frontend | `make deploy-frontend` or `./frontend/bin/deploy.sh` |

**Backend deploy** (`bin/deploy-backend.sh`): runs `go generate ./...`, `go vet ./...`, `go test ./...` from `backend/`, then `bin/deploy-lambda/run.sh` for each Lambda (order in the script matters). Each Lambda is built and deployed from its own directory via the deploy tool in `bin/deploy-lambda/`.

## Mockgen

The project uses `mockgen` for Go mocks. If needed: `GO111MODULE=on go get github.com/golang/mock/mockgen@latest` and ensure `$HOME/go/bin` is on `PATH`.

## When editing

- **Backend:** Prefer shared types and helpers in `backend/lambdas/shared/`; keep Lambda handlers thin. Run `go generate ./...` if you change interfaces that have mocks.
- **Frontend:** Use existing patterns in `src/pages/` (e.g. List/Detail/Routes, `StandardFetch`). Run Prettier (`npm run code-style` in `frontend/`).
- **New Lambdas:** Add a directory under `backend/lambdas/apis/` or `backend/lambdas/jobs/` with `main.go` and `.env.example`, then add a line to `bin/deploy-backend.sh` for `deploy-lambda/run.sh` in the desired order.
