# goland-eda

A small event-driven architecture demo: a Go API on Kafka, run with Docker Compose.

**Stack:** Go · GraphQL · REST (Swagger) · Kafka + ZooKeeper · k6 (TypeScript) · Docker Compose

## Quick start

Requires [Docker Desktop](https://www.docker.com/products/docker-desktop/).

```bash
cp .env.example .env
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build --wait
docker compose --profile test run --rm --build k6 run smoke.js
```

## URLs

| What        | URL                                    |
|-------------|----------------------------------------|
| Swagger UI  | http://localhost:8080/swagger          |
| GraphQL playground | http://localhost:8080/playground |
| GraphQL     | `POST` http://localhost:8080/graphql   |
| Health      | http://localhost:8080/health           |

Stop with `docker compose down`.

## Learn more

📖 **[docs/development.md](docs/development.md)** covers architecture, configuration, adding endpoints and tests, tooling, and deployment.

## Contributing

Work on a `feature/<name>` branch and open a pull request into `main`.
