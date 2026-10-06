# drp-identity-api

Identity HTTP API (Go). **Hexagonal:** `internal/domain` has no web/DB imports. DDL lives in [`drp-identity-db`](https://github.com/code-corhuila/drp-identity-db). Engine: [`drp-infra-postgres`](https://github.com/code-corhuila/drp-infra-postgres). This repo must not define a database container.

Contract: `drp-docs` `07-api/contracts/openapi/identity-service.yaml`. RS256 is validated **here**, not only at the gateway.

## Corte 2

This increment is a walking skeleton: `GET /health` and `{error, message, traceId}` on unknown routes. Login / JWKS / bcrypt persist in a later `feat/` on this repo. `drp-front` does not need this process (synthetic contract data).

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
