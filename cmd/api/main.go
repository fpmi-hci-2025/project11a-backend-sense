package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpDelivery "sense-backend/internal/delivery/http"
	authHandler "sense-backend/internal/delivery/http/handlers"
	"sense-backend/internal/infrastructure/database"
	"sense-backend/internal/infrastructure/jwt"
	"sense-backend/internal/infrastructure/repository"
	authUsecase "sense-backend/internal/usecase/auth"
	commentUsecase "sense-backend/internal/usecase/comment"
	feedUsecase "sense-backend/internal/usecase/feed"
	mediaUsecase "sense-backend/internal/usecase/media"
	profileUsecase "sense-backend/internal/usecase/profile"
	publicationUsecase "sense-backend/internal/usecase/publication"
	"sense-backend/pkg/config"
	"sense-backend/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func main() {
	// Load configuration - try multiple paths
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "/app/config.yaml" // Docker path
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger := logger.New()

	// Connect to database
	dbPool, err := database.NewPool(&cfg.Database)
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close(dbPool)

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)
	publicationRepo := repository.NewPublicationRepository(dbPool)
	commentRepo := repository.NewCommentRepository(dbPool)
	mediaRepo := repository.NewMediaRepository(dbPool)

	// Initialize JWT service
	tokenSvc := jwt.NewTokenService(&cfg.JWT)

	// Initialize use cases
	authUC := authUsecase.NewUseCase(userRepo, tokenSvc)
	publicationUC := publicationUsecase.NewUseCase(publicationRepo, userRepo, mediaRepo)
	commentUC := commentUsecase.NewUseCase(commentRepo)
	profileUC := profileUsecase.NewUseCase(userRepo)
	feedUC := feedUsecase.NewUseCase(publicationRepo)
	mediaUC := mediaUsecase.NewUseCase(mediaRepo)

	// Initialize validator
	validator := validator.New()

	// Initialize handlers
	authH := authHandler.NewAuthHandler(authUC, validator)
	publicationH := authHandler.NewPublicationHandler(publicationUC, validator)
	commentH := authHandler.NewCommentHandler(commentUC, validator)
	profileH := authHandler.NewProfileHandler(profileUC, validator)
	feedH := authHandler.NewFeedHandler(feedUC, validator)
	mediaH := authHandler.NewMediaHandler(mediaUC, validator, cfg.Media.MaxFileSize)

	// Initialize router
	router := httpDelivery.NewRouter(validator, appLogger, tokenSvc, authH, publicationH, commentH, profileH, feedH, mediaH)
	muxRouter := router.SetupRoutes()

	// Setup server
	srv := &http.Server{
		Handler:      muxRouter,
		Addr:         ":" + fmt.Sprintf("%d", cfg.Server.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		appLogger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			appLogger.WithError(err).Error("Server shutdown error")
		}
	}()

	appLogger.Infof("Server starting on :%d", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.WithError(err).Fatal("Server failed to start")
	}
}
