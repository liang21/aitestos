// Package http provides HTTP transport layer implementation
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/liang21/aitestos/internal/transport/http/handler"
	httpMiddleware "github.com/liang21/aitestos/internal/transport/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
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
func NewRouter(handlers *Handlers, logger zerolog.Logger, metrics *httpMiddleware.Metrics) http.Handler {
	return NewRouterWithMiddleware(handlers, "", logger, metrics)
}

// NewRouterWithMiddleware creates a new HTTP router with authentication middleware
func NewRouterWithMiddleware(handlers *Handlers, jwtSecret string, logger zerolog.Logger, metrics *httpMiddleware.Metrics) http.Handler {
	r := chi.NewRouter()

	// Add base middleware - use our custom implementations
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpMiddleware.Recovery())
	r.Use(httpMiddleware.Logging(logger))
	if metrics != nil {
		r.Use(httpMiddleware.MetricsMiddleware(metrics))
	}
	r.Use(middleware.StripSlashes)

	// Health check endpoint (public)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Metrics endpoint (public) - Prometheus scraping
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// noOpHandler is used when handlers is nil (for testing route registration)
	noOpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			if handlers != nil && handlers.Identity != nil {
				r.Post("/register", http.HandlerFunc(handlers.Identity.Register))
				r.Post("/login", http.HandlerFunc(handlers.Identity.Login))
				r.Post("/refresh", http.HandlerFunc(handlers.Identity.RefreshToken))
			} else {
				r.Post("/register", noOpHandler)
				r.Post("/login", noOpHandler)
				r.Post("/refresh", noOpHandler)
			}
		})

		// Protected routes with authentication
		r.Group(func(r chi.Router) {
			// Add authentication middleware
			if jwtSecret != "" {
				r.Use(httpMiddleware.Auth(jwtSecret))
			}

			// Project routes
			r.Route("/projects", func(r chi.Router) {
				if handlers != nil && handlers.Project != nil {
					r.Get("/", http.HandlerFunc(handlers.Project.ListProjects))
					r.Post("/", http.HandlerFunc(handlers.Project.CreateProject))
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", http.HandlerFunc(handlers.Project.GetProject))
						r.Put("/", http.HandlerFunc(handlers.Project.UpdateProject))
						r.Delete("/", http.HandlerFunc(handlers.Project.DeleteProject))
						r.Get("/stats", http.HandlerFunc(handlers.Project.GetProjectStatistics))
						r.Get("/modules", http.HandlerFunc(handlers.Project.ListModules))
						r.Post("/modules", http.HandlerFunc(handlers.Project.CreateModule))
						r.Get("/configs", http.HandlerFunc(handlers.Project.ListConfigs))
						r.Post("/configs/import", http.HandlerFunc(handlers.Project.ImportConfigs))
						r.Get("/configs/export", http.HandlerFunc(handlers.Project.ExportConfigs))
						r.Put("/configs/{key}", http.HandlerFunc(handlers.Project.SetConfig))
					})
				} else {
					r.Get("/", noOpHandler)
					r.Post("/", noOpHandler)
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", noOpHandler)
						r.Put("/", noOpHandler)
						r.Delete("/", noOpHandler)
						r.Get("/modules", noOpHandler)
						r.Post("/modules", noOpHandler)
						r.Get("/configs", noOpHandler)
						r.Put("/configs/{key}", noOpHandler)
					})
				}
			})

			// Module routes
			r.Route("/modules", func(r chi.Router) {
				if handlers != nil && handlers.Project != nil {
					r.Put("/{id}", http.HandlerFunc(handlers.Project.UpdateModule))
					r.Delete("/{id}", http.HandlerFunc(handlers.Project.DeleteModule))
				} else {
					r.Put("/{id}", noOpHandler)
					r.Delete("/{id}", noOpHandler)
				}
			})

			// Test case routes
			r.Route("/testcases", func(r chi.Router) {
				if handlers != nil && handlers.TestCase != nil {
					r.Get("/", http.HandlerFunc(handlers.TestCase.ListCases))
					r.Post("/", http.HandlerFunc(handlers.TestCase.CreateCase))
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", http.HandlerFunc(handlers.TestCase.GetCase))
						r.Put("/", http.HandlerFunc(handlers.TestCase.UpdateCase))
						r.Delete("/", http.HandlerFunc(handlers.TestCase.DeleteCase))
					})
				} else {
					r.Get("/", noOpHandler)
					r.Post("/", noOpHandler)
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", noOpHandler)
						r.Put("/", noOpHandler)
						r.Delete("/", noOpHandler)
					})
				}
			})

			// Test plan routes
			r.Route("/plans", func(r chi.Router) {
				if handlers != nil && handlers.TestPlan != nil {
					r.Get("/", http.HandlerFunc(handlers.TestPlan.ListPlans))
					r.Post("/", http.HandlerFunc(handlers.TestPlan.CreatePlan))
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", http.HandlerFunc(handlers.TestPlan.GetPlan))
						r.Put("/", http.HandlerFunc(handlers.TestPlan.UpdatePlan))
						r.Patch("/status", http.HandlerFunc(handlers.TestPlan.UpdatePlanStatus))
						r.Delete("/", http.HandlerFunc(handlers.TestPlan.DeletePlan))
						r.Post("/cases", http.HandlerFunc(handlers.TestPlan.AddCases))
						r.Delete("/cases/{caseId}", http.HandlerFunc(handlers.TestPlan.RemoveCase))
						r.Post("/results", http.HandlerFunc(handlers.TestPlan.RecordResult))
						r.Get("/results", http.HandlerFunc(handlers.TestPlan.GetResults))
					})
				} else {
					r.Get("/", noOpHandler)
					r.Post("/", noOpHandler)
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", noOpHandler)
						r.Put("/", noOpHandler)
						r.Patch("/status", noOpHandler)
						r.Delete("/", noOpHandler)
						r.Post("/cases", noOpHandler)
						r.Delete("/cases/{caseId}", noOpHandler)
						r.Post("/results", noOpHandler)
						r.Get("/results", noOpHandler)
					})
				}
			})

			// Generation routes
			r.Route("/generation", func(r chi.Router) {
				if handlers != nil && handlers.Generation != nil {
					r.Route("/tasks", func(r chi.Router) {
						r.Get("/", http.HandlerFunc(handlers.Generation.ListTasks))
						r.Post("/", http.HandlerFunc(handlers.Generation.CreateTask))
						r.Get("/{id}", http.HandlerFunc(handlers.Generation.GetTask))
						r.Get("/{id}/drafts", http.HandlerFunc(handlers.Generation.GetDrafts))
					})
					r.Route("/drafts", func(r chi.Router) {
						r.Get("/", http.HandlerFunc(handlers.Generation.ListAllDrafts))
						r.Post("/{id}/confirm", http.HandlerFunc(handlers.Generation.ConfirmDraft))
						r.Post("/{id}/reject", http.HandlerFunc(handlers.Generation.RejectDraft))
						r.Post("/batch-confirm", http.HandlerFunc(handlers.Generation.BatchConfirm))
					})
				} else {
					r.Route("/tasks", func(r chi.Router) {
						r.Get("/", noOpHandler)
						r.Post("/", noOpHandler)
						r.Get("/{id}", noOpHandler)
						r.Get("/{id}/drafts", noOpHandler)
					})
					r.Route("/drafts", func(r chi.Router) {
						r.Get("/", noOpHandler)
						r.Post("/{id}/confirm", noOpHandler)
						r.Post("/{id}/reject", noOpHandler)
						r.Post("/batch-confirm", noOpHandler)
					})
				}
			})

			// Knowledge routes
			r.Route("/knowledge", func(r chi.Router) {
				if handlers != nil && handlers.Knowledge != nil {
					r.Route("/documents", func(r chi.Router) {
						r.Get("/", http.HandlerFunc(handlers.Knowledge.ListDocuments))
						r.Post("/", http.HandlerFunc(handlers.Knowledge.UploadDocument))
						r.Route("/{id}", func(r chi.Router) {
							r.Get("/", http.HandlerFunc(handlers.Knowledge.GetDocument))
							r.Delete("/", http.HandlerFunc(handlers.Knowledge.DeleteDocument))
							r.Get("/chunks", http.HandlerFunc(handlers.Knowledge.GetChunks))
						})
					})
				} else {
					r.Route("/documents", func(r chi.Router) {
						r.Get("/", noOpHandler)
						r.Post("/", noOpHandler)
						r.Route("/{id}", func(r chi.Router) {
							r.Get("/", noOpHandler)
							r.Delete("/", noOpHandler)
							r.Get("/chunks", noOpHandler)
						})
					})
				}
			})
		})
	})

	return r
}
