# Ishari Backend

Backend API untuk aplikasi Ishari menggunakan Go, Fiber, dan PostgreSQL dengan Clean Architecture.

## 🏗️ Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan struktur:

```
ishari-backend/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── domain/          # Business entities & interfaces
│   ├── usecase/         # Business logic
│   ├── repository/      # Data access implementation
│   └── delivery/        # HTTP handlers & routes
└── pkg/                 # Shared packages
    ├── config/          # Configuration
    ├── database/        # Database connection
    └── errors/          # Error handling
```

Lihat [ARCHITECTURE.md](./ARCHITECTURE.md) untuk detail lengkap.

## 🚀 Quick Start

### Prerequisites

- Go 1.25.1 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi

### Installation

1. Clone repository
```bash
git clone <repository-url>
cd ishari-backend
```

2. Install dependencies
```bash
go mod download
```

3. Setup environment variables
```bash
cp .env.example .env
# Edit .env sesuai konfigurasi Anda
```

4. Run application
```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:3000`

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## 📦 Build

```bash
# Build binary
go build -o bin/api cmd/api/main.go

# Run binary
./bin/api
```

## 🔗 API Endpoints

### Health Check

- `GET /` - Welcome message
- `GET /health` - Basic health check
- `GET /health/db` - Database health check

## 🛠️ Development

### Project Structure

- **cmd/api**: Entry point aplikasi
- **internal/domain**: Business entities dan repository interfaces
- **internal/usecase**: Business logic layer
- **internal/repository**: Data access implementation
- **internal/delivery/http**: HTTP handlers dan routing
- **pkg**: Shared utilities (config, database, errors)

### Adding New Features

1. Define domain entities dan interfaces di `internal/domain/`
2. Implement repository di `internal/repository/`
3. Create use case di `internal/usecase/`
4. Add HTTP handlers di `internal/delivery/http/`
5. Wire dependencies di `cmd/api/main.go`

Lihat [ARCHITECTURE.md](./ARCHITECTURE.md) untuk contoh lengkap.

## 📚 Tech Stack

- **Framework**: [Fiber](https://gofiber.io/) - Express-inspired web framework
- **ORM**: [GORM](https://gorm.io/) - The fantastic ORM library for Golang
- **Database**: PostgreSQL
- **Config**: [Viper](https://github.com/spf13/viper) - Configuration management
- **Architecture**: Clean Architecture

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.
