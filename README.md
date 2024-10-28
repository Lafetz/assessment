# Introduction

Basic CRUD Api using http/net using Hexagonal Architecture.
[**LIVE SWAGGER**](https://assessment-8uar.onrender.com/swagger/index.html)

- Unit and Integration Tests
- Continuous Integration (CI) using github actions
- API Documentation using Swagger
- net/http Package
- Docker

## How to Run

### Using Docker

1. Clone the Repository
2. Build the Docker Image:

```sh
 docker build -t person-api .
 Run the Docker Container:
 docker run -p 8080:8080 person-api
```

Access the API:
The API will be available at /api/v1/persons.

View Swagger Documentation:
Access the Swagger UI at /swagger/index.html.

Run Tests:

 ```sh
    go test ./... -v
 ```

### Without Docker

Clone the Repository:

Install Dependencies: Make sure you have Go installed and set up. Then, install any required dependencies:

```sh
go mod tidy

```

Run the Application:

```sh

go run main.go

```

```sh
go test ./... -v
```
