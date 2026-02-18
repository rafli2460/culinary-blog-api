# Culinary Blog API

A high-performance REST API for a culinary blog, rewritten in Go (Fiber + SQLX) from a legacy PHP codebase. Built using Clean Architecture for a final examination project.

## About

This repository contains the backend services for a Culinary Blog platform. Originally developed in PHP, this project is a complete rework using **Go (Golang)** and the **Fiber** framework to demonstrate improved performance, type safety, and scalability.

**Context:** Developed as a final examination test, focusing on backend logic with a frontend implementation to follow.

## Tech Stack

- **Language:** [Go](https://go.dev/) (Golang)
- **Framework:** [Fiber v3](https://gofiber.io/)
- **Database:** MySQL
- **Database Access:** [SQLX](https://github.com/jmoiron/sqlx)
- **Authentication:** JWT (JSON Web Token)
- **Logging:** [Zerolog](https://github.com/rs/zerolog)
- **Environment Management:** [Godotenv](https://github.com/joho/godotenv)

## Features

- **User Authentication:** Registration, Login (JWT-based), and Logout using HttpOnly cookies.
- **Admin Management:**
  - View user statistics.
  - Manage user list with search functionality.
  - Update user roles (RBAC).
  - Delete user accounts.
- **Database Separation:** Configured for Reader/Writer database splitting for optimized scalability.
- **Structured Logging:** Console-friendly JSON logging using Zerolog.
- **Health Check:** System health monitoring endpoint.

## Project Structure

```text
├── cmd
│   └── api
│       └── app.go          # Application entry point
├── internal
│   ├── config             # Database and environment configuration
│   ├── handlers           # Request handlers (Controllers)
│   ├── middleware         # Authentication and authorization middleware
│   ├── models             # Data models and structures
│   ├── repository         # Database access layer (SQL queries)
│   ├── routes             # API route definitions
│   └── service            # Business logic layer
├── .env.example           # Example environment configuration
├── go.mod                 # Go modules file
└── go.sum                 # Go checksums file
```

## ⚙️ Configuration

Create a `.env` file based on `.env.example` and configure the following variables:

- `APP_PORT`: Port for the application to run on (default: 3000).
- `DB_WRITE_HOST`: Hostname for the writer database instance.
- `DB_READ_HOST`: Hostname for the reader database instance.
- `DB_PORT`: Database port.
- `DB_USER`: Database username.
- `DB_PASS`: Database password.
- `DB_NAME`: Database name.
- `DB_DIALECT`: Database driver name (e.g., `mysql`).
- `JWT_SECRET`: Secret key for JWT signing.

## Getting Started

1.  **Clone the repository**

    ```bash
    git clone https://github.com/rafli2460/culinary-blog-api.git
    cd culinary-blog-api
    ```

2.  **Install Dependencies**

    ```bash
    go mod tidy
    ```

3.  **Setup Environment**
    ```bash
    cp .env.example .env
    # Edit .env with your configuration
    ```

4.  **Run the Application**
    ```bash
    go run cmd/api/app.go
    ```