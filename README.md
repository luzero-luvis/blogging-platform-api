# Blogging Platform API

A small REST API for creating and reading blog posts, built with Go, `chi`, and PostgreSQL.

## What It Does

The service currently supports:

- Creating a post
- Listing all posts
- Fetching a post by ID
- Updating a post
- Checking API and database health

The project is structured with a simple layered architecture:

- `handler`: HTTP request/response handling
- `service`: validation and business logic
- `repository`: PostgreSQL queries
- `config` / `database` / `logger` / `middleware`: app setup

## Tech Stack

- Go `1.25.5`
- PostgreSQL
- [`chi`](https://github.com/go-chi/chi) router
- `slog` for logging
- Docker and Docker Compose

## Project Structure

```text
.
├── main.go
├── Dockerfile
├── docker-compose.yml
├── internal
│   ├── config
│   ├── database
│   ├── handler
│   ├── logger
│   ├── middleware
│   ├── model
│   ├── repository
│   └── service
└── README.md
```

## Environment Variables

The app reads the following variables from `.env`:

| Variable | Required | Description |
| --- | --- | --- |
| `DB_URL` | Yes | PostgreSQL connection string |
| `PORT` | No | HTTP server port. Defaults to `8080` |
| `ENV` | No | Logging mode. `dev` enables text logs; anything else uses JSON logs |

The PostgreSQL container uses `.db_env`:

- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`

Example app `.env`:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/blogdb?sslmode=disable
PORT=8080
ENV=dev
```

If you run the API inside Docker Compose, `DB_URL` should usually point to the `postgres` service, not `localhost`:

```env
DB_URL=postgres://postgres:postgres@postgres:5432/blogdb?sslmode=disable
PORT=8080
ENV=dev
```

## Database Setup

This repository does not include migrations yet, so you need to create the `posts` table yourself before using the API.

Example schema:

```sql
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT NOT NULL,
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Run Locally

1. Start PostgreSQL.
2. Create `.env` with a valid `DB_URL`.
3. Create the `posts` table.
4. Run the API:

```bash
go run .
```

The server starts on `http://localhost:8080` unless `PORT` is set.

## Run With Docker Compose

Build and start the stack:

```bash
docker compose up --build
```

Then create the schema inside Postgres if it does not already exist.

A typical flow is:

```bash
docker compose exec postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"
```

Then run the `CREATE TABLE` statement from the schema section above.

## API Endpoints

### Health Check

```http
GET /healthz
```

Checks whether the API can reach PostgreSQL.

### Create Post

```http
POST /posts
Content-Type: application/json
```

Request body:

```json
{
  "title": "My first post",
  "content": "Hello from the API",
  "category": "general",
  "tags": ["intro", "go"]
}
```

Notes:

- `title`, `content`, and `category` are required
- `tags` is optional and defaults to an empty array

### Get All Posts

```http
GET /posts
```

Returns all posts. If there are no posts, the API returns an empty array.

### Get Post By ID

```http
GET /posts/{id}
```

Example:

```bash
curl http://localhost:8080/posts/1
```

### Update Post

```http
PUT /posts/{id}
Content-Type: application/json
```

Uses the same request body shape as `POST /posts`.

## Example cURL Requests

Create a post:

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Building a Go API",
    "content": "A simple CRUD example",
    "category": "backend",
    "tags": ["go", "postgres", "rest"]
  }'
```

List posts:

```bash
curl http://localhost:8080/posts
```

Update a post:

```bash
curl -X PUT http://localhost:8080/posts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Building a Go API",
    "content": "Updated content",
    "category": "backend",
    "tags": ["go", "api"]
  }'
```

## Current Limitations

- No database migrations
- No delete endpoint
- No filtering, pagination, or search
- No automated tests in the repository yet
- Error responses are minimal and not standardized

## Next Improvements

Useful next steps for this project:

- Add SQL migrations
- Add `DELETE /posts/{id}`
- Add tests for handlers, services, and repositories
- Add request validation with clearer error responses
- Add pagination and filtering for `GET /posts`
