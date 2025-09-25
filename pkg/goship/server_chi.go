package goship

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

// chiServer is a concrete implementation of the Server interface using the Chi router.
type chiServer struct {
	router *chi.Mux
	server *http.Server
}

// NewChiServer creates and initializes a new chiServer with required middleware and routes.
func NewChiServer() Server {
	r := chi.NewRouter()

	// Add common middleware
	r.Use(middleware.RequestID)   // Assigns a unique ID to each request
	r.Use(middleware.RealIP)      // Gets the real IP from X-Forwarded-For
	r.Use(middleware.Logger)      // Logs the start and end of each request with metadata
	r.Use(middleware.Recoverer)   // Recovers from panics and writes a 500 error
	r.Use(chiRateLimiter(10, 20)) // Basic rate limiter: 10 rps with burst of 20

	// Health and metrics endpoints
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.Handle("/metrics", promhttp.Handler())

	return &chiServer{
		router: r,
	}
}

// AddCORS adds a permissive CORS policy allowing all origins, methods, and headers.
// WARNING: This should be locked down in production.
func (cs *chiServer) AddCORS() {
	cs.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
}

// GetRouter exposes the underlying router interface so additional routes can be mounted externally.
func (cs *chiServer) GetRouter() interface{} {
	return cs.router
}

// Run starts the HTTP server and sets up graceful shutdown using signals.
func (cs *chiServer) Run(port string) error {
	cs.server = &http.Server{
		Addr:         ":" + port,
		Handler:      cs.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt or terminate signal
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Shutdown signal received, attempt graceful shutdown
		log.Println("Shutting down server gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := cs.server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Starting server on port %s", port)
	if err := cs.server.ListenAndServe(); err != http.ErrServerClosed {
		// Unexpected error. Port may be in use, etc.
		return err
	}

	<-idleConnsClosed
	log.Println("Server stopped.")
	return nil
}

// chiRateLimiter creates a middleware that limits incoming requests per client IP.
// rps: requests per second, burst: maximum burst capacity
func chiRateLimiter(rps float64, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
