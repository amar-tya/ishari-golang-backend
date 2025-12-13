# Clean Architecture (Hexagonal) - Ishari Backend

Proyek ini menggunakan struktur **Clean Architecture** dengan pendekatan **Ports & Adapters (Hexagonal)**. Struktur ini memisahkan secara tegas antara **Core Business Logic** (Domain) dengan **Technical Implementation** (Adapter).

Konsep ini sangat mirip dengan pemisahan **Data - Domain - Bloc** yang umum di framework seperti Flutter.

## Struktur Folder

```
ishari-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point & Server start
├── internal/
│   ├── bootstrap/                  # Wiring / Dependency Injection setup
│   │
│   ├── core/                       # DOMAIN LAYER (Pure Business Logic)
│   │   ├── entity/                 # Data Structures (Structs Only)
│   │   ├── factory/                # Test Data Factories
│   │   ├── port/                   # Interface / Contracts (What we need)
│   │   │   ├── repository/         # Repository Interfaces
│   │   │   └── usecase/            # UseCase Interfaces (DIP)
│   │   └── usecase/                # Business Logic Implementation
│   │       ├── user/               # User feature (per-feature package)
│   │       └── book/               # Book feature
│   │
│   └── adapter/                    # ADAPTER LAYER (Implementation Details)
│       ├── repository/             # DATA LAYER (Database Impl)
│       │   └── postgres/           # Postgres Implementation of Repositories
│       └── handler/                # PRESENTATION LAYER (Routes & Controllers)
│           └── http/               # HTTP Handlers (Fiber)
│               ├── controller/
│               ├── middleware/
│               └── dto/            # Data Transfer Objects
└── pkg/                            # Shared Utilities (Config, Logger, Errors)
    ├── hasher/                     # Password hashing adapters
```

## Layer Architecture

### 1. Core Layer (`internal/core`)
Ini adalah jantung aplikasi. Layer ini **TIDAK BOLEH** bergantung pada library eksternal teknis (seperti Gorm, Fiber, Redis).

*   **Entity (`core/entity`)**: Definisi objek bisnis (misal: `Book`, `User`). Hanya berisi struct.
*   **Factory (`core/factory`)**: Builder pattern untuk membuat test data. Mirip Laravel Factory tapi untuk unit testing.
*   **Port (`core/port`)**: Interface/Kontrak mengikuti **Dependency Inversion Principle**.
    *   `port/repository/`: Interface untuk data access (diimplementasi oleh adapter/repository)
    *   `port/usecase/`: Interface untuk business logic (diimplementasi oleh usecase/)
    *   Contoh: `UserRepository` dan `UserUseCase` interface ada di sini.
*   **UseCase (`core/usecase`)**: Implementasi logika bisnis. Diorganisir per-feature dalam package terpisah.
    *   Contoh: `usecase/user/user_usecase.go` implements `port/usecase/user.go`

### 2. Adapter Layer (`internal/adapter`)
Layer ini berisi detail teknis. Layer ini tugasnya memenuhi kontrak (interface) yang dibuat oleh Core, atau menghubungkan Core dengan dunia luar.

*   **Repository (`adapter/repository`)**: Implementasi dari interface repository.
    *   Contoh: `postgres.bookRepository` yang menggunakan Gorm ada di sini.
*   **Handler (`adapter/handler`)**: Menghubungkan HTTP Request dari user ke UseCase.
    *   Mirip dengan **Bloc/View** di Flutter.
    *   Menggunakan DTO untuk validasi input/output.

### 3. Bootstrap (`internal/bootstrap`)
Tempat dimana semua komponen "dicolokkan" (wired together).
*   Membuat koneksi DB.
*   Repository Adapter dimasukkan ke UseCase.
*   UseCase dimasukkan ke Handler.
*   Handler didaftarkan ke Server.

## Dependency Rule (Dependency Inversion Principle)

Arah ketergantungan (import) harus selalu **mengarah ke dalam**:

```
┌─────────────────────────────────────────────────┐
│  Adapter Layer (Infrastructure)                 │
│  - postgres/user_repository.go                  │
│  - http/controller/user_controller.go           │
│  - pkg/hasher/bcrypt.go                         │
└────────────────┬────────────────────────────────┘
                 │ implements
                 ▼
┌─────────────────────────────────────────────────┐
│  Port Layer (Interfaces)                        │
│  - port/repository/user_repository.go           │
│  - port/usecase/user.go                         │
└────────────────┬────────────────────────────────┘
                 │ depends on
                 ▼
┌─────────────────────────────────────────────────┐
│  Core Layer (Business Logic)                    │
│  - usecase/user/user_usecase.go                 │
│  - entity/user.go                               │
└─────────────────────────────────────────────────┘
```

**Prinsip:**
*   **Core** tidak boleh import **Adapter**.
*   **Core** mendefinisikan interface di **Port**.
*   **Adapter** mengimplementasikan interface dari **Port**.
*   Dependencies mengalir dari luar ke dalam (Adapter → Port → Core).

## Menambah Fitur Baru

Contoh: Menambah fitur **Product** dengan Clean Architecture.

### 1. Domain Layer (Core)
1.  **Entity**: `core/entity/product.go` - Struct `Product`
2.  **Port Interfaces**:
    *   `core/port/repository/product_repository.go` - Interface `ProductRepository`
    *   `core/port/usecase/product.go` - Interface `ProductUseCase` + Input DTOs
3.  **UseCase Implementation**: `core/usecase/product/product_usecase.go` - Business logic
4.  **Domain Errors**: `core/usecase/product/errors.go` - Domain-specific errors
5.  **Factory**: `core/factory/product_factory.go` - Test data builder

### 2. Adapter Layer
6.  **Repository**: `adapter/repository/postgres/product_repository.go` - DB implementation
7.  **HTTP DTOs**: `adapter/handler/http/dto/product.go` - Request/Response
8.  **Controller**: `adapter/handler/http/controller/product_controller.go` - HTTP handlers
9.  **Routes**: `adapter/handler/http/product_routes.go` - Route registration

### 3. Wiring
10. **Bootstrap**: Wire dependencies di `internal/bootstrap/bootstrap.go`
11. **Routes**: Register di `adapter/handler/http/routes.go`

### 4. Testing
12. **Unit Tests**: `core/usecase/product/product_usecase_test.go` - dengan mock repository

## Testing Strategy

### Unit Testing
Dengan struktur ini, Unit Test sangat mudah dibuat:
*   Test **UseCase** dengan **Mock Repository** (tidak butuh DB nyata)
*   Gunakan **Factory** untuk membuat test data
*   Mock semua dependencies via interface

Contoh:
```go
// Mock repository
mockRepo := &MockUserRepository{
    GetByIDFunc: func(ctx, id) (*entity.User, error) {
        return factory.NewUserFactory().WithID(id).Build(), nil
    },
}

// Test usecase
uc := user.NewUserUseCase(mockRepo, mockHasher)
result, err := uc.GetByID(ctx, 1)
```

### Integration Testing
Untuk test dengan DB real, bisa buat factory yang save ke DB:
```go
type UserDBFactory struct {
    db *gorm.DB
}

func (f *UserDBFactory) Create() *entity.User {
    user := factory.NewUserFactory().Build()
    f.db.Create(user)
    return user
}
```

## Command Penting

```bash
# Menjalankan aplikasi
make run

# Menjalankan test
make test

# Menjalankan test dengan verbose
make test-v
```
