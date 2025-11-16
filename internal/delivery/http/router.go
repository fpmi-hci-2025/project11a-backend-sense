package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	authHandler "sense-backend/internal/delivery/http/handlers"
	"sense-backend/internal/delivery/http/middleware"
	"sense-backend/internal/infrastructure/jwt"
)

// Router sets up all routes
type Router struct {
	router            *mux.Router
	validator         *validator.Validate
	logger            *logrus.Logger
	tokenSvc          *jwt.TokenService
	authHandler       *authHandler.AuthHandler
	publicationHandler *authHandler.PublicationHandler
	commentHandler    *authHandler.CommentHandler
	profileHandler    *authHandler.ProfileHandler
	feedHandler       *authHandler.FeedHandler
	mediaHandler      *authHandler.MediaHandler
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
) *Router {
	return &Router{
		router:            mux.NewRouter(),
		validator:         validator,
		logger:            logger,
		tokenSvc:          tokenSvc,
		authHandler:       authHandler,
		publicationHandler: publicationHandler,
		commentHandler:    commentHandler,
		profileHandler:    profileHandler,
		feedHandler:       feedHandler,
		mediaHandler:      mediaHandler,
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

	// Feed routes (some protected, some not)
	feedRouter := r.router.PathPrefix("/feed").Subrouter()
	r.feedHandler.RegisterRoutes(feedRouter)

	// Media routes (protected)
	mediaRouter := r.router.PathPrefix("/media").Subrouter()
	mediaRouter.Use(middleware.AuthMiddleware(r.tokenSvc))
	r.mediaHandler.RegisterRoutes(mediaRouter)

	return r.router
}

