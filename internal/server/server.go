package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nathanael/organizr/internal/config"
	"github.com/nathanael/organizr/internal/downloads"
	"github.com/nathanael/organizr/internal/search"
)

type Config struct {
	Port            string
	AllowedOrigins  []string
	DownloadService *downloads.Service
	SearchService   *search.Service
	ConfigService   *config.Service
}

type Server struct {
	router          chi.Router
	httpServer      *http.Server
	downloadService *downloads.Service
	searchService   *search.Service
	configService   *config.Service
}

func New(cfg Config) *Server {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s := &Server{
		router:          router,
		downloadService: cfg.DownloadService,
		searchService:   cfg.SearchService,
		configService:   cfg.ConfigService,
	}

	s.registerRoutes()

	s.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
