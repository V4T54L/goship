package goship

type Server interface {
	// Adds an open CORS that allows all traffic
	AddCORS()
	// Returns the router as interface with some middlewares & handlers attached,
	// 
	// Like IP, logger, rate limiter, recoverer, and health & prometheus(metrics) handler 
	GetRouter() interface{}
	// Runs server with graceful shutdown on the provided port
	Run(port string) error
}
