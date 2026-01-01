package server

import "github.com/go-chi/chi/v5"

func (s *Server) registerRoutes() {
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/health", s.handleHealth)

		r.Route("/downloads", func(r chi.Router) {
			r.Post("/", s.handleCreateDownload)
			r.Get("/", s.handleListDownloads)
			r.Get("/{id}", s.handleGetDownload)
			r.Delete("/{id}", s.handleCancelDownload)
			r.Post("/{id}/organize", s.handleOrganize)
		})

		r.Route("/config", func(r chi.Router) {
			r.Get("/", s.handleGetAllConfig)
			r.Get("/{key}", s.handleGetConfig)
			r.Put("/{key}", s.handleUpdateConfig)
		})

		r.Route("/search", func(r chi.Router) {
			r.Get("/", s.handleSearch)
			r.Get("/providers", s.handleListProviders)
		})
	})
}
