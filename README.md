# 📚 Library Management API

[![Go CI/CD Pipeline](https://github.com/jdiitm/library-api/actions/workflows/ci.yml/badge.svg)](https://github.com/jdiitm/library-api/actions/workflows/ci.yml)
[![Coverage Status](https://codecov.io/gh/jdiitm/library-api/branch/main/graph/badge.svg)](https://codecov.io/gh/jdiitm/library-api)
![Go Version](https://img.shields.io/github/go-mod/go-version/jdiitm/library-api)
![Docker Pulls](https://img.shields.io/docker/pulls/jdiitm/library-api)
![License: MIT](https://img.shields.io/github/license/jdiitm/library-api)

A **production-ready REST API** for managing library resources — built with **Go (Gin)**, **PostgreSQL**, and **Docker**.

---

## 🚀 Features

- Complete CRUD operations for books
- Normalized relational schema (Books ↔ Authors ↔ Publishers)  
- Input validation & structured error responses  
- Pagination support  
- Docker + Docker Compose setup for local & CI environments  
- Unit & integration tests (with isolated Postgres test DB)  
- GitHub Actions CI/CD pipeline with linting, security scans, and Docker image publishing  

---

## 🧱 Project Structure

```plaintext
.
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── api
│   │   ├── book_handler.go
│   │   ├── book_handler_test.go
│   │   └── dto.go
│   ├── models
│   │   ├── author.go
│   │   ├── book.go
│   │   └── publisher.go
│   ├── repository
│   │   ├── book_repository.go
│   │   └── book_repository_test.go
│   ├── service
│   │   └── book_service.go
│   └── tests
│       └── integration
│           └── book_integration_test.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── .env.example
```
---
## ⚙️ Prerequisites

- **Go** 1.24 or later  
- **Docker** and **Docker Compose**  
- **PostgreSQL** (only if running manually, not via Docker)

---

## 🏃 Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/jdiitm/library-api.git
   cd library-api
2. **Copy and configure environment variables**
   ```bash
   cp .env.example .env
3. **Run with Docker Compose**
   ```bash
   docker compose up --build
The API will be available at 👉 http://localhost:8080

---

📚 API Endpoints
---
| Method | Endpoint                   | Description                |
| ------ | -------------------------- | -------------------------- |
| GET    | `/api/v1/books`            | List all books (paginated) |
| GET    | `/api/v1/books/:id`        | Get book by ID             |
| POST   | `/api/v1/books`            | Create a new book          |
| PUT    | `/api/v1/books/:id`        | Update a book              |
| DELETE | `/api/v1/books/:id`        | Delete a book              |
| POST   | `/api/v1/books/:id/issue`  | Issue a book               |
| POST   | `/api/v1/books/:id/return` | Return a book              |

---
## 🧪 Running Tests
---
1. **Unit Tests**
   ```bash
   go test ./... -v
2. **Integration Tests**
   ```bash
   DATABASE_URL="postgres://postgres:postgres@localhost:5432/library_test?sslmode=disable" 
   go test ./... -tags=integration
3. **Automated Test Script**
   ```bash
   ./run-tests.sh
---
## 🗄 Database Schema
---
1. **Books**

| Column       | Type      | Notes              |
| ------------ | --------- | ------------------ |
| id           | UUID      | Primary key        |
| title        | String    | Indexed            |
| isbn         | String    | Unique             |
| author_id    | UUID      | FK → authors.id    |
| publisher_id | UUID      | FK → publishers.id |
| year         | Integer   |                    |
| genre        | String    | Indexed            |
| quantity     | Integer   |                    |
| created_at   | Timestamp |                    |
| updated_at   | Timestamp |                    |

2. **Authors**

| Column     | Type      | Notes       |
| ---------- | --------- | ----------- |
| id         | UUID      | Primary key |
| name       | String    | Indexed     |
| biography  | Text      |             |
| created_at | Timestamp |             |
| updated_at | Timestamp |             |

3. **Publishers**

| Column     | Type      | Notes       |
| ---------- | --------- | ----------- |
| id         | UUID      | Primary key |
| name       | String    | Unique      |
| location   | String    |             |
| created_at | Timestamp |             |
| updated_at | Timestamp |             |
---
## 🧰 CI/CD Pipeline
---
This project uses GitHub Actions for automated testing, linting, and Docker image deployment.

Stages:
1. Lint & Static Analysis — Runs golangci-lint and gosec
2. Run Tests — Spins up a PostgreSQL container and executes
3. unit/integration tests
4. Build & Push Docker Image — Builds multi-arch images and pushes to Docker Hub
---
## 💻 Local Testing (Seeding, cURL Examples, Sample End-To-End Manual Testing)
---
1. **Run Docker**
      ```bash
      docker compose up --build
2. **Seeding Author(s) & Publisher(s)**
      ```bash
      # exec into db
      docker exec -it library-api-library-db-1 psql -U postgres -d library
      INSERT INTO authors (id, name) VALUES (gen_random_uuid(), 'Manual Author') RETURNING id;
      #example: 3f6e961b-3b01-4918-9863-a1e7efef71c4
      
      INSERT INTO publishers (id, name) VALUES (gen_random_uuid(), 'Manual Publisher') RETURNING id;
      #example: c8e9111b-f931-4407-9eca-871e203810a6
3. **Create a Book**
      ```bash
      curl --location 'http://localhost:8080/api/v1/books' \
      --header 'Content-Type: application/json' \
      --data '{
         "title": "Manual Test Book",
         "isbn": "9781234567890",
         "author_id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
         "publisher_id": "c8e9111b-f931-4407-9eca-871e203810a6",
         "year": 2025,
         "genre": "Fiction",
         "quantity": 3
      }'
      #201 Created
      #{
      # "id": "13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc",
      # "title": "Manual Test Book",
      # "isbn": "9781234567890",
      # "author_id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
      # "author": {
      #     "id": "00000000-0000-0000-0000-000000000000",
      #     "name": "",
      #     "biography": "",
      #     "created_at": "0001-01-01T00:00:00Z",
      #     "updated_at": "0001-01-01T00:00:00Z"
      # },
      # "publisher_id": "c8e9111b-f931-4407-9eca-871e203810a6",
      # "publisher": {
      #     "id": "00000000-0000-0000-0000-000000000000",
      #     "name": "",
      #     "location": "",
      #     "created_at": "0001-01-01T00:00:00Z",
      #     "updated_at": "0001-01-01T00:00:00Z"
      # },
      # "year": 2025,
      # "genre": "Fiction",
      # "quantity": 3,
      # "quantity_issued": 0,
      # "created_at": "2025-10-29T17:25:38.596431452Z",
      # "updated_at": "2025-10-29T17:25:38.596431452Z"
      #}
4. **List Books**
      ```bash
      curl --location 'http://localhost:8080/api/v1/books'
      #200 OK
      #{
      #    "data": [
      #        {
      #            "id": "13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc",
      #            "title": "Manual Test Book",
      #            "isbn": "9781234567890",
      #            "author_id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
      #            "author": {
      #                "id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
      #                "name": "Manual Author",
      #                "biography": "",
      #                "created_at": "0001-01-01T00:00:00Z",
      #                "updated_at": "0001-01-01T00:00:00Z"
      #            },
      #            "publisher_id": "c8e9111b-f931-4407-9eca-871e203810a6",
      #            "publisher": {
      #                "id": "c8e9111b-f931-4407-9eca-871e203810a6",
      #                "name": "Manual Publisher",
      #                "location": "",
      #                "created_at": "0001-01-01T00:00:00Z",
      #                "updated_at": "0001-01-01T00:00:00Z"
      #            },
      #            "year": 2025,
      #            "genre": "Fiction",
      #            "quantity": 3,
      #            "quantity_issued": 0,
      #            "created_at": "2025-10-29T17:25:38.596431Z",
      #            "updated_at": "2025-10-29T17:25:38.596431Z"
      #        }
      #    ],
      #    "limit": 10,
      #    "page": 1,
      #    "total": 1
      #}
4. **Get Book by ID**
      ```bash
      curl --location 'http://localhost:8080/api/v1/books/13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc'
      #200 OK
      #{
      #    "id": "13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc",
      #    "title": "Manual Test Book",
      #    "isbn": "9781234567890",
      #    "author_id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
      #    "author": {
      #        "id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
      #        "name": "Manual Author",
      #        "biography": "",
      #        "created_at": "0001-01-01T00:00:00Z",
      #        "updated_at": "0001-01-01T00:00:00Z"
      #    },
      #    "publisher_id": "c8e9111b-f931-4407-9eca-871e203810a6",
      #    "publisher": {
      #        "id": "c8e9111b-f931-4407-9eca-871e203810a6",
      #        "name": "Manual Publisher",
      #        "location": "",
      #        "created_at": "0001-01-01T00:00:00Z",
      #        "updated_at": "0001-01-01T00:00:00Z"
      #    },
      #    "year": 2025,
      #    "genre": "Fiction",
      #    "quantity": 3,
      #    "quantity_issued": 0,
      #    "created_at": "2025-10-29T17:25:38.596431Z",
      #    "updated_at": "2025-10-29T17:25:38.596431Z"
      #}
5. **Update Book**
      ```bash
      curl --location --request PUT 'http://localhost:8080/api/v1/books/13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc' \
      --header 'Content-Type: application/json' \
      --data '{
         "title": "Manual Test Book Updated",
         "isbn": "9781234567890",
         "author_id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
         "publisher_id": "c8e9111b-f931-4407-9eca-871e203810a6",
         "year": 2025,
         "genre": "Fiction",
         "quantity": 5
      }'
      #200 OK
      #{
      #    "id": "13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc",
      #    "title": "Manual Test Book Updated",
      #    "isbn": "9781234567890",
      #    "author_id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
      #    "author": {
      #        "id": "00000000-0000-0000-0000-000000000000",
      #        "name": "",
      #        "biography": "",
      #        "created_at": "0001-01-01T00:00:00Z",
      #        "updated_at": "0001-01-01T00:00:00Z"
      #    },
      #    "publisher_id": "c8e9111b-f931-4407-9eca-871e203810a6",
      #    "publisher": {
      #        "id": "00000000-0000-0000-0000-000000000000",
      #        "name": "",
      #        "location": "",
      #        "created_at": "0001-01-01T00:00:00Z",
      #        "updated_at": "0001-01-01T00:00:00Z"
      #    },
      #    "year": 2025,
      #    "genre": "Fiction",
      #    "quantity": 5,
      #    "quantity_issued": 0,
      #    "created_at": "2025-10-29T17:25:38.596431Z",
      #    "updated_at": "2025-10-29T17:30:35.208994972Z"
      #}
7. **Issue Book (5 times)**
      ```bash
      curl --location --request POST 'http://localhost:8080/api/v1/books/13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc/issue'
      #200 OK
      #{
      #    "id": "13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc",
      #    "title": "Manual Test Book Updated",
      #    "isbn": "9781234567890",
      #    "author_id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
      #    "author": {
      #        "id": "00000000-0000-0000-0000-000000000000",
      #        "name": "",
      #        "biography": "",
      #        "created_at": "0001-01-01T00:00:00Z",
      #        "updated_at": "0001-01-01T00:00:00Z"
      #    },
      #    "publisher_id": "c8e9111b-f931-4407-9eca-871e203810a6",
      #    "publisher": {
      #        "id": "00000000-0000-0000-0000-000000000000",
      #        "name": "",
      #        "location": "",
      #        "created_at": "0001-01-01T00:00:00Z",
      #        "updated_at": "0001-01-01T00:00:00Z"
      #    },
      #    "year": 2025,
      #    "genre": "Fiction",
      #    "quantity": 5,
      #    "quantity_issued": 5,
      #    "created_at": "2025-10-29T17:25:38.596431Z",
      #    "updated_at": "2025-10-29T17:31:29.891243214Z"
      #}
8. **Issue Book (6th time, expect failure)**
      ```bash
      curl --location --request POST 'http://localhost:8080/api/v1/books/13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc/issue'
      #409 Conflict
      #{
      #    "error": "failed to issue book: no available copies to issue"
      #}
9. **Return Book (5 times)**
      ```bash
      curl --location --request POST 'http://localhost:8080/api/v1/books/13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc/return'
      #200 OK
      #{
      #    "id": "13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc",
      #    "title": "Manual Test Book Updated",
      #    "isbn": "9781234567890",
      #    "author_id": "3f6e961b-3b01-4918-9863-a1e7efef71c4",
      #    "author": {
      #        "id": "00000000-0000-0000-0000-000000000000",
      #        "name": "",
      #        "biography": "",
      #        "created_at": "0001-01-01T00:00:00Z",
      #        "updated_at": "0001-01-01T00:00:00Z"
      #    },
      #    "publisher_id": "c8e9111b-f931-4407-9eca-871e203810a6",
      #    "publisher": {
      #        "id": "00000000-0000-0000-0000-000000000000",
      #        "name": "",
      #        "location": "",
      #        "created_at": "0001-01-01T00:00:00Z",
      #        "updated_at": "0001-01-01T00:00:00Z"
      #    },
      #    "year": 2025,
      #    "genre": "Fiction",
      #    "quantity": 5,
      #    "quantity_issued": 0,
      #    "created_at": "2025-10-29T17:25:38.596431Z",
      #    "updated_at": "2025-10-29T17:39:14.960720947Z"
      #}
10. **Return Book (6th time, expect failure)**
      ```bash
      curl --location --request POST 'http://localhost:8080/api/v1/books/13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc/return'
      #409 Conflict
      #{
      #    "error": "failed to return book: no issued copies to return"
      #}
11. **Delete Book**
      ```bash
      curl --location --request DELETE 'http://localhost:8080/api/v1/books/13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc'
      #204 No Content
12. **Issue Book After Delete**
      ```bash
      curl --location --request POST 'http://localhost:8080/api/v1/books/13a5a9ce-9f80-4b5e-8d8a-97ce7af1e7fc/issue'
      #409 Conflict
      #{
      #    "error": "failed to issue book: record not found"
      #}
13. **Get Book by ID After Delete**
      ```bash
      curl --location 'http://localhost:8080/api/v1/books/13a5
      #404 Not Found
      #{
      #    "error": "Book not found"
      #}

### 📝 Note

During local testing, both **Author** and **Publisher** objects are not properly preloaded.

---
## 💻 Docker Image Testing (Seeding, cURL Examples, Sample End-To-End Manual Testing)
---
1. **Pull & Run Docker**
      ```bash
      #Pull image from Docker
      docker pull jdiitm/library-api:latest
      #Run Postgres container
      docker run -d \
      --name library-db \
      -e POSTGRES_USER=postgres \
      -e POSTGRES_PASSWORD=postgres \
      -e POSTGRES_DB=library_dev \
      -p 5432:5432 \
      postgres:15-alpine
      #Run Service Container
      docker run -d \
      --name library-api \
      -p 8080:8080 \
      -e DATABASE_URL="postgres://postgres:postgres@library-db:5432/library_dev?sslmode=disable" \
      --link library-db \
      jdiitm/library-api:latest
2. **Seeding Author(s) & Publisher(s)**
      ```bash
      # exec into db
      docker exec -it library-db /bin/sh
      psql -U postgres
      \c library_dev
      INSERT INTO authors (id, name) VALUES (gen_random_uuid(), 'Manual Author') RETURNING id;
      #example: 3f6e961b-3b01-4918-9863-a1e7efef71c4

      INSERT INTO publishers (id, name) VALUES (gen_random_uuid(), 'Manual Publisher') RETURNING id;
      #example: c8e9111b-f931-4407-9eca-871e203810a6
3. **Follow steps 3-13 from section above**
### 📝 Note

During image testing, both **Author** and **Publisher** objects are properly preloaded, ensuring that related data is included in API responses.


---
## 🤝 Contributing
1. Fork the repository
2. Create your feature branch (git checkout -b feature/awesome-feature)
3. Commit changes (git commit -m 'Add awesome feature')
4. Push to your branch
5. Create a Pull Request 🚀
---

## 🧾 License

This project is licensed under the [MIT License](LICENSE).