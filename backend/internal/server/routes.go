package server

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) registerRoutes() {
	// Swagger documentation
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	s.router.Route("/api", func(r chi.Router) {
		r.Get("/health", s.handleHealth)

		r.Route("/downloads", func(r chi.Router) {
			r.Post("/", s.handleCreateDownload)
			r.Post("/batch", s.handleBatchCreateDownload)
			r.Get("/", s.handleListDownloads)
			r.Get("/{id}", s.handleGetDownload)
			r.Delete("/{id}", s.handleCancelDownload)
			r.Post("/{id}/organize", s.handleOrganize)
		})

		r.Route("/config", func(r chi.Router) {
			r.Get("/", s.handleGetAllConfig)
			r.Get("/{key}", s.handleGetConfig)
			r.Put("/{key}", s.handleUpdateConfig)
			r.Post("/preview-path", s.handlePreviewPath)
		})

		r.Route("/search", func(r chi.Router) {
			r.Get("/", s.handleSearch)
			r.Post("/test", s.handleTestConnection)
		})

		r.Route("/qbittorrent", func(r chi.Router) {
			r.Get("/test", s.handleTestQBittorrentConnection)
		})
	})
}
