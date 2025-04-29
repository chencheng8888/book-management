# Book Management System

A comprehensive book management system for library operations including book borrowing, donation tracking, and activity management.

## Technology Stack

### Core
- **Go 1.23+** - Primary programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library for database operations

### Database
- **MySQL 5.7+** - Primary relational database
- **Redis 6+** - Caching layer for performance optimization

### Infrastructure
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **Wire** - Compile-time dependency injection

### Documentation
- **Swagger** - API documentation generation
- **Makefile** - Build automation

### Testing
- **Go Test** - Unit testing framework
- **GitHub Actions** - CI/CD pipeline

## Architecture Overview

The system follows a clean architecture pattern with clear separation of concerns:

```
├── cmd/            # Application entry points
├── configs/        # Configuration files
├── docs/           # API documentation
├── internal/       # Core application logic
│   ├── controller/ # HTTP handlers
│   ├── service/    # Business logic
│   ├── repository/ # Data access layer
│   │   ├── cache/  # Redis caching
│   │   ├── dao/    # Database operations
│   │   └── repo/   # Repository interfaces
│   ├── ioc/        # Dependency injection
│   └── middleware/ # HTTP middleware
```

Key architectural features:
- Dependency injection using Wire
- Layered architecture (Controller-Service-Repository)
- Caching layer for performance
- Swagger API documentation
- Containerized deployment


## Prerequisites

- Go 1.23+
- MySQL 5.7+
- Redis 6+
- Docker (optional)

## Configuration

Create `configs/config.yaml` with following structure:

```yaml
db:
  addr: "mysql_user:mysql_password@tcp(localhost:3306)/bm?charset=utf8mb4&parseTime=True&loc=Local"
cache:
  addr: "localhost:6379"
users:
  - username: "admin"
    password: "admin123"
```

## Local Development

1. Install dependencies:
```bash
go mod download
```

2. Generate Swagger docs:
```bash
make swagger
```

3. Run application:
```bash
go run cmd/main.go -conf configs/config.yaml
```

## Docker Deployment

1. Build Docker image:
```bash
docker build -t bm .
```

2. Start services:
```bash
docker-compose up -d
```

## Status Codes

### Book Status
- `waiting_return`: Book is borrowed and waiting for return
- `returned`: Book has been returned
- `overdue`: Book is overdue

### Stock Status
- `adequate`: Stock is sufficient
- `early_warning`: Stock is running low
- `shortage`: Stock is insufficient

### Book Categories
- `children_story`: Children story books
- `science_knowledge`: Science knowledge books
- `art_enlightenment`: Art enlightenment books

### Activity Status
- `pending`: Activity is pending
- `ongoing`: Activity is ongoing
- `ended`: Activity has ended

### User Status
- `normal`: User is normal
- `overdue`: User has overdue books
- `freeze`: User is frozen

## API Documentation

### Authentication
- `POST /api/auth/login` - User login
  - Request:
    ```json
    {
      "username": "admin",
      "password": "admin123"
    }
    ```
  - Response:
    ```json
    {
      "code": 200,
      "message": "success",
      "data": {
        "token": "JWT_TOKEN"
      }
    }
    ```

### Book Management
- `GET /api/books` - Get all books
  - Response:
    ```json
    {
      "code": 200,
      "message": "success",
      "data": [
        {
          "id": 1,
          "title": "Book Title",
          "author": "Author Name",
          "status": "waiting_return"
        }
      ]
    }
    ```

- `POST /api/books/borrow` - Borrow a book
  - Request:
    ```json
    {
      "book_id": 1,
      "user_id": 1
    }
    ```
  - Response:
    ```json
    {
      "code": 200,
      "message": "success",
      "data": {
        "borrow_id": 1,
        "due_date": "2025-05-15"
      }
    }
    ```

### Activity Management
- `GET /api/activities` - Get all activities
  - Response:
    ```json
    {
      "code": 200,
      "message": "success",
      "data": [
        {
          "id": 1,
          "name": "Reading Activity",
          "status": "ongoing"
        }
      ]
    }
    ```

### Donation Management
- `POST /api/donations` - Donate a book
  - Request:
    ```json
    {
      "book_id": 1,
      "donor_id": 1,
      "quantity": 2
    }
    ```
  - Response:
    ```json
    {
      "code": 200,
      "message": "success",
      "data": {
        "donation_id": 1,
        "status": "pending"
      }
    }
    ```

For more detailed API documentation with all endpoints and parameters, you can still access the Swagger UI at:
```
http://localhost:8989/swagger/index.html
