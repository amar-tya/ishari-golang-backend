# Agents and Tools - Ishari Backend

This document describes the primary **Agents** (logic handlers) and **Tools** (infrastructure/utilities) in the Ishari backend project. It serves as a guide for developers and AI assistants to understand the system's capabilities and architecture.

## 🛠️ Infrastructure & API Engine

These components handle the "Outer Layer" (Ports & Adapters) of the system.

### [Data Engine]
- **Location**: `internal/adapter/repository/postgres`
- **Description**: Handles persistent data storage using **PostgreSQL** and **GORM**. It implements the repository interfaces defined in the core layer.
- **Capabilities**: CRUD operations, complex queries for books, verses, and users.

### [Web Engine (Fiber)]
- **Location**: `internal/adapter/handler/http`
- **Description**: The HTTP delivery layer built with the **Go Fiber** framework.
- **Capabilities**: Request routing, middleware execution, DTO validation, and JSON response orchestration.

---

## 🧠 Core Domain Agents (Business Logic)

Agents located in `internal/core/usecase` that encapsulate pure business rules.

| Agent | Responsibility | Key Features |
| :--- | :--- | :--- |
| **Auth Agent** | Security & Identity | Registration, Login, and Session management. |
| **Book Agent** | Content Catalog | Managing the library of Islamic books (Kitab). |
| **Chapter Agent** | Structure Parser | Managing chapters (Bab) within specific books. |
| **Verse Agent** | Verse Management | Handling individual verses (Bait) and lyrics. |
| **Translation Agent** | Multi-language Logic | Managing translations for verses and chapters. |
| **User Agent** | Identity Management | User profiles, permissions, and account updates. |
| **Health Agent** | System Monitor | Providing status checks for the API and Database. |

---

## 🔧 Utility Agents

Cross-cutting concerns located in the `pkg/` directory.

- **Hasher Agent** (`pkg/hasher`): Handles secure password hashing (BCRYPT).
- **JWT Agent** (`pkg/jwt`): Manages JSON Web Token signing and parsing for authentication.
- **Validator Agent** (`pkg/validation`): Performs structural validation on DTOs using `go-playground/validator`.
- **Logger Agent** (`pkg/logger`): Provides structured logging for the application.
- **Config Agent** (`pkg/config`): Loads and manages environment-specific variables.
- **Error Agent** (`pkg/errors`): Standardizes error types and HTTP status mapping.

---

## 🚀 Workflow & Development Tools

- **Makefile**: The central command runner.
  - `make run`: Starts the development server.
  - `make test`: Executes the test suite.
  - `make migrate-*`: Manages database schema versioning.
- **GolangCI-Lint**: Ensures code quality and consistency through static analysis.
- **Go Test**: The built-in testing tool used for Unit and Integration tests.
