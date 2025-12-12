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
│   │   ├── port/                   # Interface / Contracts (What we need)
│   │   │   └── repository/         # Repository Interfaces
│   │   └── usecase/                # Business Logic Implementation
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
```

## Layer Architecture

### 1. Core Layer (`internal/core`)
Ini adalah jantung aplikasi. Layer ini **TIDAK BOLEH** bergantung pada library eksternal teknis (seperti Gorm, Fiber, Redis).

*   **Entity (`core/entity`)**: Definisi objek bisnis (misal: `Book`, `User`). Hanya berisi struct.
*   **Port (`core/port`)**: Interface/Kontrak. Mendefinisikan *apa* yang dibutuhkan oleh aplikasi, tapi tidak tahu *bagaimana* cara kerjanya.
    *   Contoh: `BookRepository` interface ada di sini.
*   **UseCase (`core/usecase`)**: Logika bisnis utama. Menggabungkan entity dan menggunakan port untuk menyelesaikan tugas.

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

## Dependency Rule

Arah ketergantungan (import) harus selalu **mengarah ke dalam**:

```
DB Impl (Adapter) ────┐
                      ▼
HTTP Impl (Adapter) ──▶  Core (Domain/Ports/UseCase)
```

*   **Core** tidak boleh import **Adapter**.
*   **Adapter** bergantung pada **Core** (karena adapter mengimplementasikan port yang ada di core).

## Menambah Fitur Baru

Contoh: Menambah fitur **Product**.

1.  **Core/Entity**: Buat `core/entity/product.go` (Struct `Product`).
2.  **Core/Port**: Buat `core/port/repository/product_repo.go` (Interface `ProductRepository`).
3.  **Adapter/Repository**: Buat `adapter/repository/postgres/product_repo.go` (Implementasi DB).
4.  **Core/UseCase**: Buat `core/usecase/product_usecase.go` (Logic).
5.  **Adapter/Handler**: Buat `adapter/handler/http/controller/product_controller.go` (HTTP Endpoints).
6.  **Bootstrap**: Wire semuanya di `internal/bootstrap/bootstrap.go`.

## Testing

Dengan struktur ini, Unit Test sangat mudah dibuat:
*   Test **UseCase** dengan melakukan **Mocking** pada **Repository Port**.
*   Kita tidak butuh database nyata untuk mengetes logika bisnis.

## Command Penting

```bash
# Menjalankan aplikasi
make run

# Menjalankan test
make test

# Menjalankan test dengan verbose
make test-v
```
