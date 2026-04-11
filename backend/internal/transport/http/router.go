// Package http provides HTTP transport layer implementation
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpMiddleware "github.com/liang21/aitestos/internal/transport/http/middleware"
	"github.com/liang21/aitestos/internal/transport/http/handler"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Identity   *handler.IdentityHandler
	Project    *handler.ProjectHandler
	TestCase   *handler.TestCaseHandler
	TestPlan   *handler.TestPlanHandler
	Generation *handler.GenerationHandler
	Knowledge  *handler.KnowledgeHandler
}

// NewRouter creates a new HTTP router with all routes registered (without auth)
func NewRouter(handlers *Handlers) http.Handler {
	return NewRouterWithMiddleware(handlers, "")
}

// NewRouterWithMiddleware creates a new HTTP router with authentication middleware
func NewRouterWithMiddleware(handlers *Handlers, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	// Add base middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	// Import middleware package
	// Note: We need to use the middleware package from the same directory
	// For now, we'll add authentication logic directly

	// Health check endpoint (public)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", http.HandlerFunc(handlers.Identity.Register))
			r.Post("/login", http.HandlerFunc(handlers.Identity.Login))
			r.Post("/refresh", http.HandlerFunc(handlers.Identity.RefreshToken))
		})

		// Protected routes with authentication
		r.Group(func(r chi.Router) {
			// Add authentication middleware
			if jwtSecret != "" {
				r.Use(httpMiddleware.Auth(jwtSecret))
			}

			// Project routes
			r.Route("/projects", func(r chi.Router) {
				r.Get("/", http.HandlerFunc(handlers.Project.ListProjects))
				r.Post("/", http.HandlerFunc(handlers.Project.CreateProject))
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", http.HandlerFunc(handlers.Project.GetProject))
					r.Put("/", http.HandlerFunc(handlers.Project.UpdateProject))
					r.Delete("/", http.HandlerFunc(handlers.Project.DeleteProject))
					r.Get("/modules", http.HandlerFunc(handlers.Project.ListModules))
					r.Post("/modules", http.HandlerFunc(handlers.Project.CreateModule))
					r.Get("/configs", http.HandlerFunc(handlers.Project.ListConfigs))
					r.Put("/configs/{key}", http.HandlerFunc(handlers.Project.SetConfig))
				})
			})

			// Module routes
			r.Route("/modules", func(r chi.Router) {
				r.Delete("/{id}", http.HandlerFunc(handlers.Project.DeleteModule))
			})

			// Test case routes
			r.Route("/testcases", func(r chi.Router) {
				r.Get("/", http.HandlerFunc(handlers.TestCase.ListCases))
				r.Post("/", http.HandlerFunc(handlers.TestCase.CreateCase))
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", http.HandlerFunc(handlers.TestCase.GetCase))
					r.Put("/", http.HandlerFunc(handlers.TestCase.UpdateCase))
					r.Delete("/", http.HandlerFunc(handlers.TestCase.DeleteCase))
				})
			})

			// Test plan routes
			r.Route("/plans", func(r chi.Router) {
				r.Get("/", http.HandlerFunc(handlers.TestPlan.ListPlans))
				r.Post("/", http.HandlerFunc(handlers.TestPlan.CreatePlan))
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", http.HandlerFunc(handlers.TestPlan.GetPlan))
					r.Put("/", http.HandlerFunc(handlers.TestPlan.UpdatePlan))
					r.Delete("/", http.HandlerFunc(handlers.TestPlan.DeletePlan))
					r.Post("/cases", http.HandlerFunc(handlers.TestPlan.AddCases))
					r.Delete("/cases/{caseId}", http.HandlerFunc(handlers.TestPlan.RemoveCase))
					r.Post("/results", http.HandlerFunc(handlers.TestPlan.RecordResult))
					r.Get("/results", http.HandlerFunc(handlers.TestPlan.GetResults))
				})
			})

			// Generation routes
			r.Route("/generation", func(r chi.Router) {
				r.Route("/tasks", func(r chi.Router) {
					r.Post("/", http.HandlerFunc(handlers.Generation.CreateTask))
					r.Get("/{id}", http.HandlerFunc(handlers.Generation.GetTask))
					r.Get("/{id}/drafts", http.HandlerFunc(handlers.Generation.GetDrafts))
				})
				r.Route("/drafts", func(r chi.Router) {
					r.Post("/{id}/confirm", http.HandlerFunc(handlers.Generation.ConfirmDraft))
					r.Post("/{id}/reject", http.HandlerFunc(handlers.Generation.RejectDraft))
					r.Post("/batch-confirm", http.HandlerFunc(handlers.Generation.BatchConfirm))
				})
			})

			// Knowledge routes
			r.Route("/knowledge", func(r chi.Router) {
				r.Route("/documents", func(r chi.Router) {
					r.Get("/", http.HandlerFunc(handlers.Knowledge.ListDocuments))
					r.Post("/", http.HandlerFunc(handlers.Knowledge.UploadDocument))
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", http.HandlerFunc(handlers.Knowledge.GetDocument))
						r.Delete("/", http.HandlerFunc(handlers.Knowledge.DeleteDocument))
						r.Get("/chunks", http.HandlerFunc(handlers.Knowledge.GetChunks))
					})
				})
			})
		})
	})

	return r
}
