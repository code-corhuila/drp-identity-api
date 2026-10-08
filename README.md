# drp-identity-api

Identity HTTP API (Go). **Hexagonal:** `internal/domain` has no web/DB imports. DDL lives in [`drp-identity-db`](https://github.com/code-corhuila/drp-identity-db). Engine: [`drp-infra-postgres`](https://github.com/code-corhuila/drp-infra-postgres). This repo must not define a database container.

Contract: `drp-docs` `07-api/contracts/openapi/identity-service.yaml`. RS256 is validated **here**, not only at the gateway.

## Corte 2

**Corte 2 is the Angular shell** (`drp-front`, `environment.mode = 'failover'`). That UI works with the API off. This process is optional.

This increment: `GET /health`, `POST /api/v1/auth/login` (RS256), `GET /api/v1/auth/jwks`, `GET /api/v1/users/me`. Users are an **in-memory Corte 2 seed** (same emails as the shell). Flyway/`identity_app` persistence is a later `feat/`. Do not point `drp-front` at live mode for the checkpoint unless you are demoing this API on purpose.

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
go run ./cmd/api
curl http://localhost:8081/health
```

Compose (infra network must exist):

```bash
docker compose --env-file .env.example -f deploy/compose.yml up --build
```

## Branching

Child of `develop` named `feat/…`. Promote with `cherry-pick -x`.
