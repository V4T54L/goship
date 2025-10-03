package server

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

var _ Server = &chiServer{}

// NewChiServer initializes and returns a new instance of a Chi-based server.
// Middleware and routes must be added separately after initialization.
func NewChiServer() *chiServer {
	return &chiServer{
		router: chi.NewRouter(),
	}
}

// AddDefaultMiddleware attaches standard middleware to the router.
// This includes:
// - RequestID: Generates a unique request ID for each request
// - RealIP: Extracts client IP from X-Forwarded-For
// - Logger: Logs request and response metadata
// - Recoverer: Recovers from panics and returns 500 error
// - RateLimiter: Basic rate limiting to protect server
func (cs *chiServer) AddDefaultMiddleware() {
	cs.router.Use(middleware.RequestID)
	cs.router.Use(middleware.RealIP)
	cs.router.Use(middleware.Logger)
	cs.router.Use(middleware.Recoverer)
	cs.router.Use(chiRateLimiter(10, 20)) // Limit: 10 requests/sec, burst: 20
}

// AddPermissiveCORS attaches a permissive CORS policy to the router.
//
// ⚠ WARNING: This configuration allows all origins, methods, and headers
// and should NOT be used in production environments without restrictions.
func (cs *chiServer) AddPermissiveCORS() {
	cs.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
}

// AddDefaultRoutes registers common endpoints on the router.
//
// - GET /health: Returns "OK" (used for health checks)
// - GET /metrics: Prometheus metrics handler
func (cs *chiServer) AddDefaultRoutes() {
	cs.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	cs.router.Handle("/metrics", promhttp.Handler())
}

// GetRouter returns the internal Chi router instance.
// This allows the caller to mount additional custom routes externally.
func (cs *chiServer) GetRouter() interface{} {
	return cs.router
}

// Run starts the HTTP server on the specified port and handles
// graceful shutdown on system interrupt (e.g., Ctrl+C).
//
// It listens for incoming connections and shuts down cleanly within
// 10 seconds after receiving a termination signal.
func (cs *chiServer) Run(port string) error {
	cs.server = &http.Server{
		Addr:         ":" + port,
		Handler:      cs.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to signal when shutdown is complete
	idleConnsClosed := make(chan struct{})

	// Listen for interrupt signal and initiate graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down server gracefully...")

		// Create context with timeout for shutdown operations
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := cs.server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Starting server on port %s", port)
	if err := cs.server.ListenAndServe(); err != http.ErrServerClosed {
		// Unexpected error occurred
		return err
	}

	// Wait for graceful shutdown to complete
	<-idleConnsClosed
	log.Println("Server stopped.")
	return nil
}

// chiRateLimiter returns a middleware that rate-limits incoming requests
// using a token bucket algorithm.
//
// rps: Requests per second
// burst: Maximum burst size
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
