package server

// Server defines the interface for an HTTP server with middleware, CORS,
// routing, and graceful shutdown support.
type Server interface {
	// AddDefaultMiddleware attaches common middleware like logging,
	// recovery, request ID, real IP, and rate limiting.
	AddDefaultMiddleware()

	// AddPermissiveCORS attaches a permissive CORS policy allowing all origins,
	// methods, and headers.
	// 
	// WARNING: This is insecure for production use.
	AddPermissiveCORS()

	// AddDefaultRoutes registers default endpoints such as /health and /metrics.
	AddDefaultRoutes()

	// GetRouter returns the underlying router for mounting custom routes.
	GetRouter() interface{}

	// Run starts the HTTP server on the specified port and gracefully shuts down
	// on interrupt signals.
	Run(port string) error
}
