# Culinary Blog API

A high-performance REST API for a culinary blog, rewritten in Go (Fiber + SQLX) from a legacy PHP codebase. Built using Clean Architecture for a final examination project.

## About

This repository contains the backend services for a Culinary Blog platform. Originally developed in PHP, this project is a complete rework using **Go (Golang)** and the **Fiber** framework to demonstrate improved performance, type safety, and scalability.

**Context:** Developed as a final examination test, focusing on backend logic with a frontend implementation to follow.

## Tech Stack

- **Language:** [Go](https://go.dev/) (Golang)
- **Framework:** [Fiber v2](https://gofiber.io/)
- **Database:** MySQL
- **Database Access:** [SQLX](https://github.com/jmoiron/sqlx)
- **Logging:** [Zerolog](https://github.com/rs/zerolog)
- **Architecture:** Modular Clean Architecture

## Features

- **User Management:** Registration, Login, and Role-based access control (RBAC).
- **Database Separation:** Configured for Reader/Writer database splitting for scalability.
- **High Performance:** Built on Fiber, one of the fastest Go web frameworks.
- **Structured Logging:** JSON logging using Zerolog for observability.

## Project Structure

```text
├── internal
|   ├── config       # Configuration
│   ├── domain       # Domain entities (Structs)
│   ├── modules      # Feature modules (e.g., Users) containing Service & Repository layers
│   ├── platform     # Infrastructure (Database connections, SQL preparation)
│   └── server       # Server configuration, wiring, and dependency injection
└── cmd/app          # Application entry point
```

## ⚙️ Configuration

The application requires the following database configuration (typically via environment variables or a config file):

- `DB_HOST_WRITER`: Hostname for the writer database instance.
- `DB_HOST_READER`: Hostname for the reader database instance.
- `DB_PORT`: Database port.
- `DB_USER`: Database username.
- `DB_PASS`: Database password.
- `DB_NAME`: Database name.
- `DB_DIALECT`: Database driver name (e.g., `mysql`).

## 🏃 Getting Started

1.  **Clone the repository**

    ```bash
    git clone https://github.com/rafli2460/culinary-blog-api.git
    cd culinary-blog-api
    ```

2.  **Install Dependencies**

    ```bash
    go mod tidy
    ```

3.  **Run the Application**
    ```bash
    go run main.go
    ```
