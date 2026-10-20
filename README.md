# todo

REST API for task management in Go: users sign up, sign in with JWT and manage their lists and items.

> Work in progress: tests, Docker and CI are coming next.

## Stack

- Go, Gin
- PostgreSQL, sqlx, SQL migrations
- JWT (golang-jwt/jwt v5), bcrypt for passwords
- Viper for config, secrets from environment variables
- log/slog for structured JSON logs
- Swagger for API docs

## What's inside

- Layered structure: handler → service → repository
- Transactions when a list or an item is created
- Fail-fast startup: the server refuses to start without required secrets or with a weak JWT key
- Graceful shutdown on SIGINT / SIGTERM

## Run locally

1. Start PostgreSQL (the config expects port 5436):

   ```bash
   docker run --name todo-db -e POSTGRES_PASSWORD=yourpassword -p 5436:5432 -d postgres:16
   ```

2. Apply migrations with [golang-migrate](https://github.com/golang-migrate/migrate):

   ```bash
   migrate -path ./schema -database "postgres://postgres:yourpassword@localhost:5436/postgres?sslmode=disable" up
   ```

3. Create `.env` from the example and fill in `DB_PASSWORD` and `JWT_SIGNING_KEY` (how to generate the key is described in `.env.example`):

   ```bash
   cp .env.example .env
   ```

4. Run the server:

   ```bash
   go run ./cmd/todo
   ```

API docs: http://localhost:8000/swagger/index.html

## API

| Method | Path | Description |
|---|---|---|
| POST | /auth/sign-up | Register a user |
| POST | /auth/sign-in | Get a JWT |
| POST, GET | /api/lists/ | Create a list, get all lists |
| GET, PUT, DELETE | /api/lists/:id | Get, update, delete a list |
| POST, GET | /api/lists/:id/items/ | Create an item, get all items of a list |
| GET, PUT, DELETE | /api/items/:id | Get, update, delete an item |

All `/api` routes require the `Authorization: Bearer <token>` header.

## License

GPL-3.0