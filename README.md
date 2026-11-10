# drp-identity-api

Identity HTTP API (Go). **Hexagonal:** `internal/domain` has no web/DB imports. DDL lives in [`drp-identity-db`](https://github.com/code-corhuila/drp-identity-db). Engine: [`drp-infra-postgres`](https://github.com/code-corhuila/drp-infra-postgres). This repo must not define a database container.

Contract: `drp-docs` `07-api/contracts/openapi/identity-service.yaml`. RS256 is validated **here**, not only at the gateway.

## Corte 2

**Corte 2 is the Angular shell** (`drp-front`, `environment.mode = 'failover'`). That UI works with the API off. This process is optional.

This increment: `GET /health`, `POST /api/v1/auth/login` (RS256), `GET /api/v1/auth/jwks`, `GET /api/v1/users/me`. Empty `DATABASE_URL` keeps the **in-memory Corte 2 seed** (same emails as the shell). After Flyway in `drp-identity-db`, set `DATABASE_URL` as `identity_app` — this process **reads** `identity.users` / `identity.user_roles` (no DDL here). Placeholder hashes `PENDING%` are replaced with bcrypt on boot (V002). Do not point `drp-front` at live mode for the Corte 2 checkpoint unless you are demoing this API on purpose.

## Layout

| Path | Layer |
|------|--------|
| `internal/domain` | entities (no adapters) |
| `internal/app` | use cases |
| `internal/adapters/http` | JSON, correlation id |
| `cmd/api` | process |

## Run

```bash
go test ./...

# Persist (host): Flyway first, then this process against the published port.
# $env:DATABASE_URL="postgres://identity_app:change-me-identity@localhost:5436/drp?sslmode=disable"
go run ./cmd/api
curl http://localhost:8081/health
```

Compose (infra network must exist). Do **not** pass `.env.example` as `--env-file`: that DSN is for a **host** process (`localhost`). The compose file already injects `DATABASE_URL` to `postgres:5432` on the infra network.

```bash
docker compose -f deploy/compose.yml up --build
```

## Branching

Child of `develop` named `hu-c3-*-dev`. Promote with `cherry-pick -x`.
