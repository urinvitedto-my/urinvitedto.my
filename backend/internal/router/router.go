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
	mw "github.com/urinvitedto-my/backend/internal/middleware"
)

// Router manages HTTP routing and middleware.
type Router struct {
	cfg      *config.Config
	router   *chi.Mux
	db       *pgxpool.Pool
	handlers *handlers.Handlers
	mw       *mw.Middleware
}

func NewRouter(cfg *config.Config, db *pgxpool.Pool) *Router {
	h := handlers.New(db)
	return &Router{
		cfg:      cfg,
		router:   chi.NewRouter(),
		db:       db,
		handlers: h,
		mw:       mw.New(cfg, h),
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

	h := rm.handlers

	// API routes
	r.Route("/api/v1", func(api chi.Router) {
		// health check
		api.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
		})

		// event routes (public)
		api.Route("/events/{type}/{slug}", func(er chi.Router) {
			er.Get("/summary", h.GetEventSummary)
			er.Get("/details", h.GetEventDetails)
			er.Get("/confirmed-guests", h.GetConfirmedGuests)
			er.Post("/rsvp", h.PostRSVP)
		})

		// admin routes (protected)
		api.Route("/admin", func(ar chi.Router) {
			ar.Use(rm.mw.Auth)
			ar.Use(rm.mw.RequireAdmin)
			ar.Get("/events", h.ListEvents)
			ar.Post("/events", h.CreateEvent)
			ar.Put("/events/{id}", h.UpdateEvent)
			ar.Delete("/events/{id}", h.DeleteEvent)
			ar.Post("/events/{id}/hosts", h.AddHost)
			ar.Delete("/events/{id}/hosts/{hostId}", h.DeleteHost)
			ar.Get("/events/{id}/invites", h.ListInvites)
			ar.Post("/events/{id}/invites", h.CreateInvite)
			ar.Delete("/events/{id}/invites/{inviteId}", h.DeleteInvite)
			ar.Post("/events/{id}/invites/{inviteId}/guests", h.AddGuest)
			ar.Put("/events/{id}/guests/{guestId}", h.UpdateGuest)
			ar.Delete("/events/{id}/guests/{guestId}", h.DeleteGuest)
			ar.Get("/events/{id}/schedule", h.ListSchedule)
			ar.Post("/events/{id}/schedule", h.CreateScheduleItem)
			ar.Put("/events/{id}/schedule/{itemId}", h.UpdateScheduleItem)
			ar.Delete("/events/{id}/schedule/{itemId}", h.DeleteScheduleItem)
			ar.Get("/events/{id}/faqs", h.ListFAQs)
			ar.Post("/events/{id}/faqs", h.CreateFAQ)
			ar.Put("/events/{id}/faqs/{itemId}", h.UpdateFAQ)
			ar.Delete("/events/{id}/faqs/{itemId}", h.DeleteFAQ)
			ar.Get("/events/{id}/gifts", h.ListGifts)
			ar.Post("/events/{id}/gifts", h.CreateGift)
			ar.Put("/events/{id}/gifts/{itemId}", h.UpdateGift)
			ar.Delete("/events/{id}/gifts/{itemId}", h.DeleteGift)
		})

		// host routes (protected - any authenticated user)
		api.Route("/host", func(hr chi.Router) {
			hr.Use(rm.mw.Auth)
			hr.Get("/events", h.GetHostEvents)
		})
	})

	return r
}
