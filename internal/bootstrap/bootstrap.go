package bootstrap

import (
	"fmt"

	"ishari-backend/internal/adapter/handler/http"
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	"ishari-backend/internal/adapter/repository/postgres"
	"ishari-backend/internal/core/usecase"
	authusecase "ishari-backend/internal/core/usecase/auth"
	userusecase "ishari-backend/internal/core/usecase/user"
	"ishari-backend/pkg/config"
	"ishari-backend/pkg/database"
	"ishari-backend/pkg/hasher"
	"ishari-backend/pkg/jwt"
	"ishari-backend/pkg/logger"
	"ishari-backend/pkg/validation"
)

type App struct {
	Server  *http.Server
	Cleanup func() error
}

// Build composes all dependencies and registers routes, returning a ready Server and cleanup.
func Build(cfg config.Config) (*App, error) {
	// DB
	db, err := database.Connect(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}
	cleanup := func() error { return database.Close(db) }

	// Logger (file with fallback)
	l, lerr := logger.NewFile("logs/app.log")
	if lerr != nil {
		l = logger.New()
	}

	// Infrastructure
	passwordHasher := hasher.NewBcryptHasher(12)
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)

	// Repositories
	bookRepo := postgres.NewBookRepository(db)
	healthRepo := postgres.NewHealthRepository(db)
	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)

	// Token blacklist (database-backed)
	tokenBlacklist := jwt.NewDatabaseBlacklist(refreshTokenRepo)

	// Use cases
	healthUC := usecase.NewHealthUseCase(healthRepo)
	bookUC := usecase.NewBookUseCase(bookRepo)
	userUC := userusecase.NewUserUseCase(userRepo, passwordHasher)
	authUC := authusecase.NewAuthUseCase(userRepo, jwtService, tokenBlacklist, passwordHasher)

	// HTTP server
	server := http.NewServer(cfg.Server, l)
	middleware.Setup(server.App)

	// Controllers and routes
	v := validation.New()
	healthCtrl := controller.NewHealthController(healthUC)
	bookCtrl := controller.NewBookController(bookUC, v, l)
	userCtrl := controller.NewUserController(userUC, v, l)
	authCtrl := controller.NewAuthController(authUC, v, l)

	http.RegisterRoutes(server.App, http.Controllers{
		Health: healthCtrl,
		Book:   bookCtrl,
		User:   userCtrl,
		Auth:   authCtrl,
	}, &http.AuthDeps{
		AuthUC: authUC,
	})

	return &App{Server: server, Cleanup: cleanup}, nil
}
