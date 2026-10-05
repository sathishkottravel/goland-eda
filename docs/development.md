# Development guide

- [Architecture](#architecture)
- [Folder structure](#folder-structure)
- [Tooling](#tooling)
- [Configuration](#configuration)
- [Running the stack](#running-the-stack)
- [api service](#api-service)
- [Kafka](#kafka)
- [k6 API tests](#k6-api-tests)
- [Deployment notes](#deployment-notes)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## Architecture

```
            ┌──────────────────────── eda-net (Docker network) ────────────────────────┐
 host       │                                                                           │
 :8080 ───► │  api (Go)  ──(KAFKA_BROKERS=kafka:29092, wired up in a later feature)──►  │
            │                                       kafka :29092 (INTERNAL)  ◄──► zookeeper :2181
 :9092 ───► │                                       kafka :9092  (EXTERNAL)             │
            │  k6 (profile "test", runs on demand) ──► http://api:8080                  │
            └───────────────────────────────────────────────────────────────────────────┘
```

- **zookeeper / kafka** — Confluent Platform 7.x images (the last line that supports ZooKeeper). Single broker; data kept in named volumes.
- **api** — the only backend service. Serves REST, GraphQL and Swagger UI. Kafka producers/consumers will be added here.
- **k6** — load/API tests, only started when explicitly run.

## Folder structure

```
goland-eda/
├── backend/
│   └── api/                    # service "api" — its own Go module
│       ├── main.go             # wiring, HTTP server, graceful shutdown, `api healthcheck`
│       ├── internal/config/    # env-based configuration
│       ├── internal/service/   # business logic shared by REST and GraphQL
│       ├── internal/rest/      # REST handlers (+ tests)
│       ├── internal/graph/     # GraphQL schema + resolvers (+ tests)
│       ├── internal/docs/      # openapi.yaml + Swagger UI page (embedded)
│       ├── Dockerfile          # multi-stage → distroless, non-root
│       ├── .air.toml           # live reload config
│       └── .golangci.yml       # lint config
├── tests/k6/                   # isolated k6 project (TypeScript → JavaScript)
├── docs/                       # this guide
├── docker-compose.yml          # base stack, deploy-safe
├── docker-compose.dev.yml      # local overrides (published ports, topic auto-create)
└── .env.example                # all configuration variables
```

New backend services go in `backend/<name>/` as their own module. The frontend will live in `frontend/`.

## Tooling

Only Docker is required to run the stack. For working on Go code on the host (Windows, via winget):

```powershell
winget install --id GoLang.Go -e                  # Go 1.27
winget install --id GolangCI.golangci-lint -e     # golangci-lint v2
go install golang.org/x/tools/gopls@latest        # language server
go install github.com/go-delve/delve/cmd/dlv@latest   # debugger
go install github.com/air-verse/air@latest        # live reload
```

k6 and Node are **not** needed on the host: the k6 project is built and run inside Docker.

## Configuration

All settings are environment variables. Copy `.env.example` to `.env` (git-ignored); Compose reads it automatically.

| Variable                   | Default          | Used by     | Purpose |
|----------------------------|------------------|-------------|---------|
| `CP_VERSION`               | `7.7.1`          | zookeeper, kafka | Confluent Platform image tag |
| `KAFKA_ADVERTISED_HOST`    | `localhost`      | kafka       | Hostname external clients use to reach Kafka |
| `KAFKA_EXTERNAL_PORT`      | `9092`           | kafka       | External listener port |
| `KAFKA_REPLICATION_FACTOR` | `1`              | kafka       | Internal topic replication |
| `KAFKA_AUTO_CREATE_TOPICS` | `false`          | kafka       | Auto-create topics (forced `true` in dev override) |
| `PORT`                     | `8080`           | api         | HTTP listen port |
| `KAFKA_BROKERS`            | `kafka:29092`    | api         | Broker list (`localhost:9092` when running api on host) |
| `API_PORT`                 | `8080`           | dev override | Host port the api is published on |
| `K6_BASE_URL`              | `http://api:8080`| k6          | Target for tests |

**Base vs dev compose:** `docker-compose.yml` publishes no ports and is what you'd deploy. `docker-compose.dev.yml` adds host ports for Kafka (9092) and api (8080) and enables topic auto-creation.

| Port  | Service   | Published in |
|-------|-----------|--------------|
| 8080  | api       | dev override |
| 9092  | kafka (EXTERNAL) | dev override |
| 29092 | kafka (INTERNAL) | never (containers only) |
| 2181  | zookeeper | never |

## Running the stack

```bash
# start (dev) and wait until everything is healthy
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build --wait
docker compose ps

# logs
docker compose logs -f api

# stop (keeps Kafka data) / stop and wipe volumes
docker compose down
docker compose down -v
```

## api service

### Design

Kept deliberately slim — one third-party dependency:

- **REST:** standard library `net/http` with method + path patterns (`GET /health`).
- **GraphQL:** [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go) — schema-first, resolvers are plain Go methods, no code generation.
- **Swagger:** a hand-written `internal/docs/openapi.yaml`, embedded into the binary with `go:embed`. `/swagger` serves a small HTML page that loads Swagger UI from the jsDelivr CDN.
- **Shared service layer:** REST handlers and GraphQL resolvers both call `internal/service`, so logic is written once.

### Endpoints

| Method | Path            | Description |
|--------|-----------------|-------------|
| GET    | `/health`       | Liveness: `{"status":"ok"}` |
| GET    | `/api/v1/hello` | Greeting, optional `?name=` |
| POST   | `/graphql`      | GraphQL (`health`, `hello(name)`) |
| GET    | `/openapi.yaml` | OpenAPI 3 spec |
| GET    | `/swagger`      | Swagger UI |

```bash
curl localhost:8080/health
curl "localhost:8080/api/v1/hello?name=demo"
curl -X POST localhost:8080/graphql -H "Content-Type: application/json" \
  -d '{"query":"{ hello(name: \"demo\") }"}'
```

### Adding an endpoint

1. Add the logic to `internal/service` (new file or method).
2. **REST:** add a handler in `internal/rest` and register it in `Register`; document it in `internal/docs/openapi.yaml`.
3. **GraphQL:** add the field to `internal/graph/schema.graphql` and a matching method on `Resolver` (method name = field name, capitalised; arguments as a struct).
4. Add tests next to the code, and a k6 check (see below).

### Working on the host

```bash
cd backend/api
go test ./...
golangci-lint run
air                      # live reload on :8080; set PORT to change
```

If the dev stack is also running, either stop its `api` container or run on another port (`PORT=8090 air`). To reach Kafka from the host, set `KAFKA_BROKERS=localhost:9092`.

The container image is built from `backend/api/Dockerfile`: Go builds a static binary, which runs in `distroless/static` as a non-root user. Because distroless has no shell or curl, the container healthcheck runs `/api healthcheck`.

## Kafka

Kafka has two listeners:

- **INTERNAL** `kafka:29092` — used by other containers (api, CLI tools inside the network).
- **EXTERNAL** `${KAFKA_ADVERTISED_HOST}:9092` — used by clients outside Docker. The advertised host is what Kafka tells clients to reconnect to, so it must be resolvable by those clients.

Quick check from inside the broker container:

```bash
docker compose exec kafka kafka-topics --bootstrap-server kafka:29092 --create --if-not-exists --topic demo --partitions 1 --replication-factor 1
echo "hello" | docker compose exec -T kafka kafka-console-producer --bootstrap-server kafka:29092 --topic demo
docker compose exec kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic demo --from-beginning --max-messages 1
docker compose exec kafka kafka-topics --bootstrap-server kafka:29092 --list
```

## k6 API tests

Tests are written in TypeScript and compiled to JavaScript before k6 runs them. The project is fully isolated: its own `package.json` and lockfile, built and run in Docker, and the `k6` compose service sits behind the `test` profile so `docker compose up` never starts it.

```
tests/k6/
├── src/lib/config.ts   # BASE_URL (from env) and shared thresholds
├── src/lib/http.ts     # get(path), graphql(query, vars)
├── src/tests/smoke.ts  # one file per test → dist/<name>.js
├── build.mjs           # esbuild bundles every src/tests/*.ts
├── tsconfig.json       # strict typecheck only (no emit)
└── Dockerfile          # node: npm ci → typecheck → build; then grafana/k6 with /scripts
```

Run (the stack must be up):

```bash
docker compose --profile test run --rm --build k6 run smoke.js
# against another environment
docker compose --profile test run --rm -e BASE_URL=https://api.example.com k6 run smoke.js
```

**Add a test:** create `src/tests/<name>.ts` exporting `options` and a default function, reuse helpers from `src/lib`, then run `... k6 run <name>.js`. The build picks up new files automatically.

Optional host workflow (needs Node and k6 installed locally): `cd tests/k6 && npm ci && npm run build && k6 run dist/smoke.js`. Running `npm ci` locally also makes your editor find the `k6` type definitions.

## Deployment notes

- Use only `docker-compose.yml` (no dev override) and publish ports deliberately.
- Set `KAFKA_ADVERTISED_HOST` to the server's DNS name if external clients need Kafka; otherwise don't publish 9092 at all.
- ZooKeeper and the INTERNAL listener are never published.
- Kafka/ZooKeeper data live in named volumes (`kafka-data`, `zk-data`, `zk-log`); back them up or use a managed Kafka.
- Keep `KAFKA_AUTO_CREATE_TOPICS=false` and create topics explicitly.
- Secrets go in the platform's secret store, never in `.env` files in git.
- This is a single-broker setup: fine for a demo, not highly available.

## Troubleshooting

- **Git Bash rewrites paths** like `/scripts/x.js` into `C:/Program Files/Git/...`. The k6 service sets `working_dir: /scripts`, so use `k6 run smoke.js`; or prefix commands with `MSYS_NO_PATHCONV=1`.
- **Line endings:** `.gitattributes` forces LF so shell scripts and YAML work in Linux containers. If a file shows as modified after checkout, run `git add --renormalize .`.
- **Port already in use:** change `API_PORT` / `KAFKA_EXTERNAL_PORT` in `.env`.
- **Kafka stays unhealthy:** check `docker compose logs kafka`; after config changes a stale volume can be the cause — `docker compose down -v` resets it (deletes data).
- **Editor says "Cannot find type definition file for 'k6'":** expected until you run `npm ci` in `tests/k6`; Docker builds are unaffected.

## Contributing

1. Branch from `main`: `git checkout -b feature/<name>`.
2. Make changes; run `go test ./...`, `golangci-lint run` and the k6 smoke test.
3. Push and open a pull request into `main`.
4. Never commit `.env`, secrets, or local tool config (`.claude/`, IDE folders) — `.gitignore` covers them, but check `git status` before committing.
