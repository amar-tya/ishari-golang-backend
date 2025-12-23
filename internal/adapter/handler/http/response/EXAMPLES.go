package response

/*
CONTOH REFACTORING USER CONTROLLER

=== BEFORE ===
func (h *UserController) Register(c *fiber.Ctx) error {
    var req dto.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        h.log.Error("Register body parse error", "error", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid request body",
            "error":   err.Error(),
        })
    }

    if err := h.validate.Struct(req); err != nil {
        h.log.Error("Register validation failed", "error", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "validation failed",
            "error":   err.Error(),
        })
    }

    input := portuc.RegisterUserInput{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    user, err := h.userUseCase.Register(c.UserContext(), input)
    if err != nil {
        h.log.Error("Register failed", "error", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status":  "success",
        "message": "user registered successfully",
        "data":    h.toUserResponse(user),
    })
}

=== AFTER ===
func (h *UserController) Register(c *fiber.Ctx) error {
    var req dto.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return response.SendParseError(c, err, h.log, "Register body parse error")
    }

    if err := h.validate.Struct(req); err != nil {
        return response.SendValidationError(c, err, h.log, "Register validation failed")
    }

    input := portuc.RegisterUserInput{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    user, err := h.userUseCase.Register(c.UserContext(), input)
    if err != nil {
        return response.SendBadRequest(c, err.Error(), err, h.log, "Register failed")
    }

    return response.SendCreated(c, "user registered successfully", h.toUserResponse(user))
}

PERBANDINGAN:
- Before: 38 lines
- After: 20 lines
- Pengurangan: 47% lebih sedikit kode!
- Lebih mudah dibaca dan di-maintain

=== CONTOH LAIN: GetUserByID ===

BEFORE:
func (h *UserController) GetUserByID(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid user id",
        })
    }

    user, err := h.userUseCase.GetByID(c.UserContext(), uint(id))
    if err != nil {
        h.log.Error("GetUserByID failed", "error", err)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   h.toUserResponse(user),
    })
}

AFTER:
func (h *UserController) GetUserByID(c *fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return response.SendBadRequest(c, "invalid user id", err, nil, "")
    }

    user, err := h.userUseCase.GetByID(c.UserContext(), uint(id))
    if err != nil {
        return response.SendNotFound(c, err.Error(), err, h.log, "GetUserByID failed")
    }

    return response.SendOK(c, h.toUserResponse(user))
}

=== CONTOH: ListUsers dengan Pagination ===

BEFORE:
func (h *UserController) ListUsers(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "20"))
    search := c.Query("search", "")

    params := portuc.ListUserParams{
        Page:   page,
        Limit:  limit,
        Search: search,
    }

    result, err := h.userUseCase.List(c.UserContext(), params)
    if err != nil {
        h.log.Error("ListUsers failed", "error", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "failed to list users",
            "error":   err.Error(),
        })
    }

    var totalPages int
    if limit > 0 {
        totalPages = int(math.Ceil(float64(result.Total) / float64(limit)))
    }

    out := make([]dto.UserResponse, 0, len(result.Data))
    for _, u := range result.Data {
        out = append(out, h.toUserResponse(&u))
    }

    return c.JSON(fiber.Map{
        "data": out,
        "meta": fiber.Map{
            "page":        page,
            "limit":       limit,
            "total":       result.Total,
            "total_pages": totalPages,
            "count":       len(result.Data),
        },
    })
}

AFTER:
func (h *UserController) ListUsers(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "20"))
    search := c.Query("search", "")

    params := portuc.ListUserParams{
        Page:   page,
        Limit:  limit,
        Search: search,
    }

    result, err := h.userUseCase.List(c.UserContext(), params)
    if err != nil {
        return response.SendInternalError(c, "failed to list users", err, h.log, "ListUsers failed")
    }

    var totalPages int
    if limit > 0 {
        totalPages = int(math.Ceil(float64(result.Total) / float64(limit)))
    }

    out := make([]dto.UserResponse, 0, len(result.Data))
    for _, u := range result.Data {
        out = append(out, h.toUserResponse(&u))
    }

    return response.SendPaginated(c, out, page, limit, result.Total, totalPages, len(result.Data))
}

TIPS REFACTORING:
1. Import package: import "ishari-backend/internal/adapter/handler/http/response"
2. Replace semua error response dengan helper function yang sesuai
3. Replace success response dengan SendOK, SendCreated, atau SendPaginated
4. Jika tidak perlu logging, pass nil untuk parameter log
5. Jika log nil, logContext bisa string kosong ""
*/
