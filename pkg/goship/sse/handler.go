package sse

import (
	"log"
	"net/http"
)

// Handler wraps a function that uses the SSE Writer.
func Handler(fn func(Writer, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set necessary headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no") // For nginx buffering

		sseW, err := newSSEWriter(w)
		if err != nil {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		// Allow handler to use the writer
		err = fn(sseW, r)

		if err != nil {
			log.Printf("SSE handler error: %v", err)
			return
		}
	}
}
