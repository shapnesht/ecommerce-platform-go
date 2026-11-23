# Ecommerce Platform (Go)

A REST API for an online shop, written in Go with Gin and GORM. It covers authentication, a product catalogue with image uploads, shopping carts and orders. Order events are published to SQS, and a notifier service is planned to send email.

## Features

- JWT authentication with refresh tokens, and admin/user roles
- Categories and products (admin-managed, publicly readable)
- Product image uploads to S3 or the local filesystem
- Shopping cart and order management
- Domain events published through [Watermill](https://github.com/ThreeDotsLabs/watermill) to AWS SQS
- SMTP email notifications (`internal/notifications`)
- Swagger/OpenAPI docs and structured logging (zerolog)

## Tech Stack

| Concern | Choice |
| --- | --- |
| Language | Go 1.26 |
| HTTP | Gin |
| ORM / DB | GORM, PostgreSQL |
| Migrations | golang-migrate (SQL files in `db/migrations`) |
| Auth | golang-jwt, bcrypt (`x/crypto`) |
| Storage | AWS S3 (LocalStack locally) or local disk |
| Messaging | Watermill + AWS SQS |
| Local mail | Mailpit |
| API docs | swaggo |

## Project Layout

```
cmd/api          API entry point
cmd/notifier     Event consumer / email sender (work in progress)
internal/
  config         Environment-based configuration
  database       DB connection
  dto, models    Request/response types and GORM models
  server         Routes, handlers, middleware
  services       Business logic
  events         Event publisher (Watermill)
  notifications  SMTP email
  providers      S3 / local upload providers
  utils          JWT, password hashing, response helpers
db/migrations    SQL migrations
docker           docker-compose, nginx (CDN), LocalStack init
docs             Generated Swagger docs
```

## Getting Started

### Prerequisites

- Go 1.26+
- Docker and Docker Compose
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI
- [`swag`](https://github.com/swaggo/swag) (only to regenerate docs)
- `golangci-lint` (optional)

### Setup

```sh
# 1. Start Postgres, LocalStack (S3 + SQS), Mailpit and the nginx CDN
make docker-up

# 2. Create a .env file (see Configuration below)

# 3. Apply migrations
make migrate-up

# 4. Run the API
make run
```

The API listens on `http://localhost:8080` by default.

### Local services

| Service | URL / Port |
| --- | --- |
| API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/docs/index.html |
| Postgres | localhost:5432 |
| LocalStack (S3, SQS) | http://localhost:4566 |
| Mailpit UI / SMTP | http://localhost:8025 / localhost:1025 |
| CDN (nginx, serves `uploads/`) | http://localhost:8081 |

## Configuration

Configuration is read from environment variables (a `.env` file is loaded if present). Defaults are shown below.

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8080` | HTTP port |
| `GIN_MODE` | `debug` | Gin mode |
| `DB_HOST` / `DB_PORT` | `localhost` / `5432` | Database address |
| `DB_USERNAME` / `DB_PASSWORD` | `postgres` / `password` | Database credentials |
| `DB_NAME` | `ecommerce_shop` | Database name |
| `DB_SSLmode` | `disable` | Postgres SSL mode |
| `JWT_SECRET_KEY` | `secretkey` | **Change in production** |
| `JWT_EXPIRY` | `24h` | Access token lifetime |
| `REFRESH_TOKEN_EXPIRY` | `72h` | Refresh token lifetime |
| `UPLOAD_PROVIDER` | `s3` | `s3` or `local` |
| `UPLOAD_PATH` | `./uploads` | Local upload directory |
| `MAX_UPLOAD_SIZE` | `10485760` | Max upload size in bytes |
| `AWS_REGION` | `us-east-1` | AWS region |
| `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` | `test` / `secret` | AWS credentials |
| `AWS_S3_BUCKET_NAME` | `ecommerce-uploads` | Upload bucket |
| `AWS_S3_ENDPOINT` | `localhost:9000` | S3 endpoint (LocalStack: `http://localhost:4566`) |
| `AWS_EventQueueName` | `ecommerce-events` | SQS queue for events |
| `SMTP_HOST` / `SMTP_PORT` | `localhost` / `1025` | SMTP server |
| `SMTP_USERNAME` / `SMTP_PASSWORD` | empty | SMTP credentials |
| `SMTP_FROM` | `noreply@shop.com` | Sender address |

The migration targets in the Makefile use database `ecommerce-shop`, while the config default is `ecommerce_shop`. Set `DB_NAME` to match.

## API Overview

Base path: `/api/v1`. Full details are in Swagger. Protected routes need an `Authorization: Bearer <token>` header.

| Area | Endpoints | Access |
| --- | --- | --- |
| Health | `GET /health` | Public |
| Auth | `POST /auth/register`, `/auth/login`, `/auth/refresh`, `/auth/logout` | Public |
| Catalogue | `GET /products`, `GET /products/:id`, `GET /categories` | Public |
| Users | `GET`/`PUT /users/profile` | Authenticated |
| Categories | `POST /categories`, `PUT`/`DELETE /categories/:id` | Admin |
| Products | `POST`/`PUT`/`DELETE /products`, `POST /products/:id/images` | Admin |
| Cart | `GET /cart`, `POST /cart/items`, `PUT`/`DELETE /cart/items/:id` | Authenticated |
| Orders | `POST /orders`, `GET /orders`, `GET /orders/:id` | Authenticated |

## Make Targets

| Command | Purpose |
| --- | --- |
| `make build` | Build the API binary to `bin/app` |
| `make run` | Run the API |
| `make lint` | Run golangci-lint |
| `make migrate-up` / `migrate-down` | Apply / revert migrations |
| `make docker-up` | Start local infrastructure |
| `make docs-generate` | Regenerate Swagger docs |

## Contributing

Work on a feature branch and keep commits focused. Run `make lint` before opening a PR.
