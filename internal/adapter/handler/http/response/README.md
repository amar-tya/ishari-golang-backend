# Global Response Handler

Package ini menyediakan helper functions untuk standardisasi response API di seluruh controller.

## Struktur

```
internal/adapter/handler/http/response/
├── error.go    # Error response helpers
└── success.go  # Success response helpers
```

## Error Response Helpers

### Fungsi Utama

#### `SendError(ctx, statusCode, message, err, log, logContext)`
Fungsi dasar untuk mengirim error response dengan custom status code.

#### Helper Functions

1. **`SendBadRequest(ctx, message, err, log, logContext)`**
   - Status Code: `400 Bad Request`
   - Digunakan untuk: Invalid input, validation errors, parsing errors

2. **`SendUnauthorized(ctx, message, err, log, logContext)`**
   - Status Code: `401 Unauthorized`
   - Digunakan untuk: Authentication failures

3. **`SendNotFound(ctx, message, err, log, logContext)`**
   - Status Code: `404 Not Found`
   - Digunakan untuk: Resource tidak ditemukan

4. **`SendInternalError(ctx, message, err, log, logContext)`**
   - Status Code: `500 Internal Server Error`
   - Digunakan untuk: Database errors, unexpected errors

5. **`SendValidationError(ctx, err, log, logContext)`**
   - Shortcut untuk validation errors
   - Message: "validation failed"

6. **`SendParseError(ctx, err, log, logContext)`**
   - Shortcut untuk body parsing errors
   - Message: "invalid request body"

### Contoh Penggunaan Error

```go
import "ishari-backend/internal/adapter/handler/http/response"

// Bad Request
if err := strconv.Atoi(id); err != nil {
    return response.SendBadRequest(ctx, "invalid ID format", err, nil, "")
}

// Internal Server Error dengan logging
result, err := c.usecase.List(ctx.UserContext(), params)
if err != nil {
    return response.SendInternalError(ctx, "failed to list items", err, c.log, "List items failed")
}

// Validation Error
if err := c.validate.Struct(req); err != nil {
    return response.SendValidationError(ctx, err, c.log, "Validation failed")
}

// Parse Error
var req dto.CreateRequest
if err := ctx.BodyParser(&req); err != nil {
    return response.SendParseError(ctx, err, c.log, "Body parse error")
}
```

## Success Response Helpers

### Fungsi Utama

#### `SendSuccess(ctx, statusCode, message, data)`
Fungsi dasar untuk mengirim success response dengan custom status code.

#### Helper Functions

1. **`SendOK(ctx, data)`**
   - Status Code: `200 OK`
   - Digunakan untuk: GET requests yang berhasil

2. **`SendCreated(ctx, message, data)`**
   - Status Code: `201 Created`
   - Digunakan untuk: POST requests yang berhasil membuat resource

3. **`SendPaginated(ctx, data, page, limit, total, totalPages, count)`**
   - Status Code: `200 OK`
   - Digunakan untuk: List endpoints dengan pagination
   - Otomatis menambahkan metadata pagination

### Contoh Penggunaan Success

```go
import "ishari-backend/internal/adapter/handler/http/response"

// Simple OK response
user, err := c.usecase.GetByID(ctx.UserContext(), id)
if err != nil {
    return response.SendInternalError(ctx, "failed to get user", err, c.log, "GetByID failed")
}
return response.SendOK(ctx, c.toUserResponse(user))

// Created response
user, err := c.usecase.Create(ctx.UserContext(), input)
if err != nil {
    return response.SendBadRequest(ctx, err.Error(), err, c.log, "Create failed")
}
return response.SendCreated(ctx, "user created successfully", c.toUserResponse(user))

// Paginated response
result, err := c.usecase.List(ctx.UserContext(), params)
if err != nil {
    return response.SendInternalError(ctx, "failed to list users", err, c.log, "List failed")
}

var totalPages int
if limit > 0 {
    totalPages = int(math.Ceil(float64(result.Total) / float64(limit)))
}

out := make([]dto.UserResponse, 0, len(result.Data))
for _, u := range result.Data {
    out = append(out, c.toUserResponse(&u))
}

return response.SendPaginated(ctx, out, page, limit, result.Total, totalPages, len(result.Data))
```

## Response Format

### Error Response
```json
{
  "status": "error",
  "message": "failed to list chapters",
  "error": "database connection failed"
}
```

### Success Response
```json
{
  "status": "success",
  "message": "user created successfully",
  "data": {
    "id": 1,
    "username": "john_doe"
  }
}
```

### Paginated Response
```json
{
  "data": [
    {
      "id": 1,
      "title": "Chapter 1"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5,
    "count": 20
  }
}
```

## Parameter Explanation

### Error Helpers
- `ctx`: Fiber context
- `message`: User-friendly error message
- `err`: Error object (bisa nil)
- `log`: Logger instance (bisa nil jika tidak perlu logging)
- `logContext`: Context untuk log message (kosongkan jika log nil)

### Success Helpers
- `ctx`: Fiber context
- `data`: Response data (any type)
- `message`: Success message
- `page`: Current page number
- `limit`: Items per page
- `total`: Total items (int64)
- `totalPages`: Total pages
- `count`: Items in current page

## Migration Guide

### Before (Old Code)
```go
if err != nil {
    c.log.Error("List chapters failed", "error", err)
    return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "status":  "error",
        "message": "failed to list chapters",
        "error":   err.Error(),
    })
}
```

### After (New Code)
```go
if err != nil {
    return response.SendInternalError(ctx, "failed to list chapters", err, c.log, "List chapters failed")
}
```

## Benefits

✅ **Konsistensi**: Semua response memiliki format yang sama  
✅ **DRY**: Tidak ada duplikasi kode error handling  
✅ **Maintainability**: Mudah mengubah format response di satu tempat  
✅ **Logging**: Otomatis log error dengan context  
✅ **Type Safety**: Compile-time checking untuk parameter
