package sse

import (
	"bufio"
	"fmt"
	"net/http"
	"time"
)

type Writer interface {
	Event(event, data string) error
	Comment(comment string) error
	Retry(duration time.Duration) error
	Flush() error
}

type sseWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
	bw      *bufio.Writer
}

func newSSEWriter(w http.ResponseWriter) (*sseWriter, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("response writer does not support flushing")
	}

	bw := bufio.NewWriter(w)

	return &sseWriter{
		w:       w,
		flusher: flusher,
		bw:      bw,
	}, nil
}

func (s *sseWriter) Event(event, data string) error {
	if len(event) > 0 {
		if _, err := fmt.Fprintf(s.bw, "event: %s\n", event); err != nil {
			return err
		}
	}

	if len(data) > 0 {
		if _, err := fmt.Fprintf(s.bw, "data: %s\n\n", data); err != nil {
			return err
		}
	}
	return s.Flush()
}

func (s *sseWriter) Comment(comment string) error {
	if _, err := fmt.Fprintf(s.bw, ": %s\n\n", comment); err != nil {
		return err
	}
	return s.Flush()
}

func (s *sseWriter) Retry(duration time.Duration) error {
	if _, err := fmt.Fprintf(s.bw, "retry: %d\n\n", duration.Milliseconds()); err != nil {
		return err
	}
	return s.Flush()
}

func (s *sseWriter) Flush() error {
	if err := s.bw.Flush(); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}
