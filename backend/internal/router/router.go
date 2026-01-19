// Package router sets up the HTTP router and middleware.
package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urinvitedto-my/backend/internal/config"
	"github.com/urinvitedto-my/backend/internal/handlers"
)

type Router struct {
	cfg    *config.Config
	router *chi.Mux
	db     *pgxpool.Pool
}

func NewRouter(cfg *config.Config, db *pgxpool.Pool) *Router {
	return &Router{
		cfg:    cfg,
		router: chi.NewRouter(),
		db:     db,
	}
}

// SetupRouter configures and returns the chi router.
func (rm *Router) SetupRouter() *chi.Mux {
	r := rm.router

	// middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByRealIP(60, time.Minute))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Compress(5))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{rm.cfg.FrontEndURL},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Invite-Code",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// init handlers
	h := handlers.New(rm.db)

	// API routes
	r.Route("/api/v1", func(api chi.Router) {
		// health check
		api.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
		})

		// event routes
		api.Route("/events/{type}/{slug}", func(er chi.Router) {
			er.Get("/summary", h.GetEventSummary)
			er.Get("/details", h.GetEventDetails)
			er.Get("/confirmed-guests", h.GetConfirmedGuests)
			er.Post("/rsvp", h.PostRSVP)
		})
	})

	return r
}
