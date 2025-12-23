package http

import (
	"encoding/json"
	"net/http"

	authHandler "sense-backend/internal/delivery/http/handlers"
	"sense-backend/internal/delivery/http/middleware"
	"sense-backend/internal/infrastructure/jwt"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Router sets up all routes
type Router struct {
	router              *mux.Router
	validator           *validator.Validate
	logger              *logrus.Logger
	tokenSvc            *jwt.TokenService
	authHandler         *authHandler.AuthHandler
	publicationHandler  *authHandler.PublicationHandler
	commentHandler      *authHandler.CommentHandler
	profileHandler      *authHandler.ProfileHandler
	feedHandler         *authHandler.FeedHandler
	mediaHandler        *authHandler.MediaHandler
	aiHandler           *authHandler.AIHandler
	searchHandler       *authHandler.SearchHandler
	notificationHandler *authHandler.NotificationHandler
}

// NewRouter creates a new router
func NewRouter(
	validator *validator.Validate,
	logger *logrus.Logger,
	tokenSvc *jwt.TokenService,
	authHandler *authHandler.AuthHandler,
	publicationHandler *authHandler.PublicationHandler,
	commentHandler *authHandler.CommentHandler,
	profileHandler *authHandler.ProfileHandler,
	feedHandler *authHandler.FeedHandler,
	mediaHandler *authHandler.MediaHandler,
	aiHandler *authHandler.AIHandler,
	searchHandler *authHandler.SearchHandler,
	notificationHandler *authHandler.NotificationHandler,
) *Router {
	return &Router{
		router:              mux.NewRouter(),
		validator:           validator,
		logger:              logger,
		tokenSvc:            tokenSvc,
		authHandler:         authHandler,
		publicationHandler:  publicationHandler,
		commentHandler:      commentHandler,
		profileHandler:      profileHandler,
		feedHandler:         feedHandler,
		mediaHandler:        mediaHandler,
		aiHandler:           aiHandler,
		searchHandler:       searchHandler,
		notificationHandler: notificationHandler,
	}
}

// SetupRoutes configures all routes
func (r *Router) SetupRoutes() *mux.Router {
	// Health check
	r.router.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	// Auth routes (no auth required)
	r.authHandler.RegisterRoutes(r.router, r.tokenSvc)

	// Auth check (requires auth)
	r.router.Handle("/auth/check",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.authHandler.Check))).Methods("GET")

	// Publication routes (protected)
	publicationRouter := r.router.PathPrefix("/publication").Subrouter()
	publicationRouter.Use(middleware.AuthMiddleware(r.tokenSvc))
	r.publicationHandler.RegisterRoutes(publicationRouter, r.commentHandler)

	// Comment routes (protected)
	commentRouter := r.router.PathPrefix("/comment").Subrouter()
	commentRouter.Use(middleware.AuthMiddleware(r.tokenSvc))
	r.commentHandler.RegisterRoutes(commentRouter)

	// Profile routes (protected)
	profileRouter := r.router.PathPrefix("/profile").Subrouter()
	profileRouter.Use(middleware.AuthMiddleware(r.tokenSvc))
	r.profileHandler.RegisterRoutes(profileRouter)

	// Feed routes (some protected, some with optional auth)
	feedRouter := r.router.PathPrefix("/feed").Subrouter()
	// Public routes with optional auth (to get is_liked status for authenticated users)
	feedRouter.Handle("", middleware.OptionalAuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.feedHandler.GetFeed))).Methods("GET")
	feedRouter.Handle("/user/{id}", middleware.OptionalAuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.feedHandler.GetUser))).Methods("GET")
	// Protected routes
	feedRouter.Handle("/me", middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.feedHandler.GetMe))).Methods("GET")
	feedRouter.Handle("/me/saved", middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.feedHandler.GetSaved))).Methods("GET")

	// Media routes (protected)
	mediaRouter := r.router.PathPrefix("/media").Subrouter()
	mediaRouter.Use(middleware.AuthMiddleware(r.tokenSvc))
	r.mediaHandler.RegisterRoutes(mediaRouter)

	// AI routes (protected)
	r.router.Handle("/recommendations",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.aiHandler.GetRecommendations))).Methods("POST")
	r.router.Handle("/recommendations/feed",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.aiHandler.GetRecommendationsFeed))).Methods("GET")
	r.router.Handle("/recommendations/{id}/hide",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.aiHandler.HideRecommendation))).Methods("POST")
	r.router.Handle("/purify",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.aiHandler.PurifyText))).Methods("POST")

	// Search routes (mixed auth - some optional, some required)
	r.router.Handle("/search", middleware.OptionalAuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.searchHandler.SearchPublications))).Methods("GET")
	r.router.Handle("/search/users",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.searchHandler.SearchUsers))).Methods("GET")
	r.router.Handle("/search/warmup",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.searchHandler.WarmupIndex))).Methods("POST")
	r.router.Handle("/tags",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.searchHandler.GetTags))).Methods("GET")

	// Follow routes (protected)
	followRouter := r.router.PathPrefix("/follow").Subrouter()
	followRouter.Use(middleware.AuthMiddleware(r.tokenSvc))
	r.profileHandler.RegisterFollowRoutes(followRouter)

	// Notification routes (protected)
	r.router.Handle("/notifications",
		middleware.AuthMiddleware(r.tokenSvc)(http.HandlerFunc(r.notificationHandler.GetNotifications))).Methods("GET")

	return r.router
}
