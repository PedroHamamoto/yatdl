# YATDL - Yet  Another To-Do List

A lightweight _To-Do List_ API, built for learning Go

---

## Requirements
* **[Go](https://go.dev/)** (version 1.25+)
* **[Docker & Docker Compose](https://www.docker.com/)**
* **[Golangci-lint](https://golangci-lint.run/)**

---
## Getting Started

### 1. Spin up the Database
```bash
make db-up
```
### 2. Execute the Migrations
```bash
make migrate-up
```

### 3. Run the application
```bash
make api-start
```
The API will be running at http://localhost:8080

---

## API
Execute `docker-compose up` and `go run cmd/api/main.go` to run the API locally.
Access the swagger docs at http://localhost:8081

---

## Development
### Creating migrations
To create a new set of migration files (up and down) with the current timestamp:
```bash
make migrate-create NAME=create_users_table
```

### Lint
```bash
make lint
```