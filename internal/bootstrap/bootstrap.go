package bootstrap

import (
	"fmt"

	"ishari-backend/internal/adapter/handler/http"
	"ishari-backend/internal/adapter/handler/http/controller"
	"ishari-backend/internal/adapter/handler/http/middleware"
	"ishari-backend/internal/adapter/repository/postgres"
	"ishari-backend/internal/core/usecase"
	authusecase "ishari-backend/internal/core/usecase/auth"
	bookusecase "ishari-backend/internal/core/usecase/book"
	bookmarkusecase "ishari-backend/internal/core/usecase/bookmark"
	chapterusecase "ishari-backend/internal/core/usecase/chapter"
	dashboardusecase "ishari-backend/internal/core/usecase/dashboard"
	hadiusecase "ishari-backend/internal/core/usecase/hadi"
	translationusecase "ishari-backend/internal/core/usecase/translation"
	userusecase "ishari-backend/internal/core/usecase/user"
	verseusecase "ishari-backend/internal/core/usecase/verse"
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
	chapterRepo := postgres.NewChapterRepository(db)
	healthRepo := postgres.NewHealthRepository(db)
	userRepo := postgres.NewUserRepository(db)
	verseRepo := postgres.NewVerseRepository(db)
	translationRepo := postgres.NewTranslationRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	bookmarkRepo := postgres.NewBookmarkRepository(db)
	hadiRepo := postgres.NewHadiRepository(db)
	dashboardRepo := postgres.NewDashboardRepository(db)

	// Token blacklist (database-backed)
	tokenBlacklist := jwt.NewDatabaseBlacklist(refreshTokenRepo)

	// Use cases
	healthUC := usecase.NewHealthUseCase(healthRepo)
	bookUC := bookusecase.NewBookUseCase(bookRepo)
	chapterUC := chapterusecase.NewChapterUsecase(chapterRepo, bookRepo, l)
	userUC := userusecase.NewUserUseCase(userRepo, passwordHasher)
	verseUC := verseusecase.NewVerseUsecase(verseRepo, chapterRepo, l)
	translationUC := translationusecase.NewTranslationUsecase(translationRepo, verseRepo, l)
	authUC := authusecase.NewAuthUseCase(userRepo, jwtService, tokenBlacklist, passwordHasher)
	bookmarkUC := bookmarkusecase.NewBookmarkUsecase(bookmarkRepo, verseRepo, l)
	hadiUC := hadiusecase.NewHadiUseCase(hadiRepo)
	dashboardUC := dashboardusecase.NewDashboardUseCase(dashboardRepo)

	// HTTP server
	server := http.NewServer(cfg.Server, l)
	middleware.Setup(server.App)

	// Controllers and routes
	v := validation.New()
	healthCtrl := controller.NewHealthController(healthUC)
	bookCtrl := controller.NewBookController(bookUC, v, l)
	chapterCtrl := controller.NewChapterController(chapterUC, v, l)
	userCtrl := controller.NewUserController(userUC, v, l)
	authCtrl := controller.NewAuthController(authUC, v, l)
	verseCtrl := controller.NewVerseController(verseUC, v, l)
	translationCtrl := controller.NewTranslationController(translationUC, v, l)
	bookmarkCtrl := controller.NewBookmarkController(bookmarkUC, v, l)
	hadiCtrl := controller.NewHadiController(hadiUC, v, l)
	dashboardCtrl := controller.NewDashboardController(dashboardUC, l)

	http.RegisterRoutes(server.App, http.Controllers{
		Health:      healthCtrl,
		Book:        bookCtrl,
		Chapter:     chapterCtrl,
		User:        userCtrl,
		Auth:        authCtrl,
		Verse:       verseCtrl,
		Translation: translationCtrl,
		Bookmark:    bookmarkCtrl,
		Hadi:        hadiCtrl,
		Dashboard:   dashboardCtrl,
	}, &http.AuthDeps{
		AuthUC: authUC,
	})

	return &App{Server: server, Cleanup: cleanup}, nil
}
