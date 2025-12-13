# Migration Guide - Old Structure to Clean Architecture

## Overview

Panduan ini menjelaskan cara migrasi dari struktur lama ke struktur Clean Architecture yang baru.

## Struktur Lama vs Baru

### Struktur Lama
```
ishari-backend/
├── config/
├── database/
├── errors/
├── handlers/
├── routes/
├── server/
└── main.go
```

### Struktur Baru
```
ishari-backend/
├── cmd/api/
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   └── delivery/
└── pkg/
    ├── config/
    ├── database/
    └── errors/
```

## Mapping File Lama ke Baru

| File Lama | File Baru | Keterangan |
|-----------|-----------|------------|
| `main.go` | `cmd/api/main.go` | Entry point dengan dependency injection |
| `config/config.go` | `pkg/config/config.go` | Tanpa perubahan signifikan |
| `database/database.go` | `pkg/database/postgres.go` | Rename dan update import |
| `errors/errors.go` | `pkg/errors/errors.go` | Tanpa perubahan |
| `handlers/health.go` | `internal/delivery/http/health_handler.go` | Refactor dengan use case |
| `routes/routes.go` | Integrated ke handlers | Routes didefinisikan di masing-masing handler |
| `server/server.go` | `internal/delivery/http/server.go` | Refactor struktur |
| - | `internal/domain/health.go` | **NEW**: Domain entities & interfaces |
| - | `internal/usecase/health_usecase.go` | **NEW**: Business logic |
| - | `internal/repository/health_repository.go` | **NEW**: Data access |

## Step-by-Step Migration

### Step 1: Backup Kode Lama
```bash
# Buat branch baru untuk backup
git checkout -b backup/old-structure
git add .
git commit -m "Backup old structure before migration"

# Kembali ke main branch
git checkout main
```

### Step 2: Verifikasi Struktur Baru
Pastikan semua file baru sudah dibuat dengan benar:
```bash
# Check struktur folder
tree -L 3

# Atau di Windows PowerShell
Get-ChildItem -Recurse -Depth 3
```

### Step 3: Update Import Paths
Jika ada file test atau file lain yang masih menggunakan import lama, update:

**Lama:**
```go
import (
    "ishari-backend/config"
    "ishari-backend/database"
    "ishari-backend/handlers"
)
```

**Baru:**
```go
import (
    "ishari-backend/pkg/config"
    "ishari-backend/pkg/database"
    "ishari-backend/internal/delivery/http"
)
```

### Step 4: Test Aplikasi Baru
```bash
# Run dari entry point baru
go run cmd/api/main.go

# Test endpoints
curl http://localhost:3000/health
curl http://localhost:3000/health/db
```

### Step 5: Run Tests
```bash
# Run semua tests
go test ./...

# Jika ada test yang gagal, update import paths
```

### Step 6: Hapus File Lama (Optional)
Setelah verifikasi berhasil, Anda bisa hapus file-file lama:

```bash
# Hapus folder lama
rm -rf config/
rm -rf database/
rm -rf errors/
rm -rf handlers/
rm -rf routes/
rm -rf server/
rm main.go

# Atau pindahkan ke folder backup
mkdir -p _old_structure
mv config/ database/ errors/ handlers/ routes/ server/ main.go _old_structure/
```

## Perubahan Penting

### 1. Dependency Injection
**Lama** (Global variable):
```go
// database/database.go
var DB *gorm.DB

// handlers/health.go
db := database.GetDB()
```

**Baru** (Injected via constructor):
```go
// internal/usecase/health_usecase.go
func NewHealthUseCase(db *gorm.DB) HealthUseCase {
    return &healthUseCase{
        healthRepo: repository.NewHealthRepository(db),
    }
}
```

### 2. Handler Structure
**Lama** (Function-based):
```go
func HealthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"status": "ok"})
}
```

**Baru** (Struct-based dengan use case):
```go
type HealthHandler struct {
    healthUseCase usecase.HealthUseCase
}

func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
    status := h.healthUseCase.CheckHealth()
    return c.JSON(status)
}
```

### 3. Routes Registration
**Lama** (Centralized routes):
```go
// routes/routes.go
func Register(app *fiber.App) {
    app.Get("/health", handlers.HealthCheck)
}
```

**Baru** (Handler-based routes):
```go
// internal/delivery/http/health_handler.go
func NewHealthHandler(app *fiber.App, healthUseCase usecase.HealthUseCase) {
    handler := &HealthHandler{healthUseCase: healthUseCase}
    app.Get("/health", handler.HealthCheck)
}
```

## Testing Strategy

### Unit Tests
Test setiap layer secara terpisah:

```go
// internal/usecase/health_usecase_test.go
func TestHealthUseCase_CheckHealth(t *testing.T) {
    // Mock repository
    mockRepo := &mockHealthRepository{}
    uc := &healthUseCase{healthRepo: mockRepo}
    
    // Test
    status := uc.CheckHealth()
    assert.Equal(t, "ok", status.Status)
}
```

### Integration Tests
Test interaksi antar layer dengan database test.

## Troubleshooting

### Import Error
**Problem**: `package xxx is not in GOROOT`

**Solution**: 
```bash
go mod tidy
go mod download
```

### Database Connection Error
**Problem**: Database tidak terkoneksi

**Solution**: 
- Pastikan `.env` sudah dikonfigurasi dengan benar
- Cek PostgreSQL sudah running
- Verifikasi credentials

### Port Already in Use
**Problem**: `bind: address already in use`

**Solution**:
```bash
# Windows
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:3000 | xargs kill -9
```

## Rollback Plan

Jika terjadi masalah, Anda bisa rollback ke struktur lama:

```bash
# Checkout branch backup
git checkout backup/old-structure

# Atau restore dari _old_structure
cp -r _old_structure/* .
```

## Next Steps

Setelah migrasi berhasil:

1. ✅ Update CI/CD pipeline untuk menggunakan `cmd/api/main.go`
2. ✅ Update Docker/Kubernetes configuration
3. ✅ Tambahkan fitur baru mengikuti clean architecture
4. ✅ Buat documentation untuk team
5. ✅ Setup linting dan code quality tools

## Questions?

Jika ada pertanyaan atau masalah, silakan buka issue di repository atau hubungi team lead.
