package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/V4T54L/goship/pkg/goship"
	"github.com/V4T54L/goship/pkg/goship/sse"
	"github.com/V4T54L/goship/pkg/goship/ws"
	"github.com/go-chi/chi/v5"
)

type PingMessage struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

func main() {
	db, err := goship.ConnectToSqliteDb("./deleteThisDB.sqlite")
	defer func() {
		log.Println("Closing the db...")
		err := goship.CloseSqlDBConn(db)
		if err != nil {
			log.Println("Error closing the db")
		} else {
			log.Println("DB closed.")
		}
	}()
	log.Printf("Database: %v \n\nError: %v", db, err)

	// ---

	server := goship.NewChiServer()
	server.AddDefaultMiddleware()
	server.AddPermissiveCORS()
	server.AddDefaultRoutes()

	r, ok := server.GetRouter().(*chi.Mux)
	if !ok {
		log.Fatal("Error obtaining the router")
	}

	// ---

	r.Get("/events", sse.Handler(func(w sse.Writer, r *http.Request) error {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// Optional: send retry interval
		w.Retry(2 * time.Second)

		for {
			select {
			case <-r.Context().Done():
				// Client disconnected
				return nil
			case t := <-ticker.C:
				msg := fmt.Sprintf(`{"time": "%s"}`, t.Format(time.RFC3339))
				if err := w.Event("tick", msg); err != nil {
					return err
				}
			}
		}
	}))

	// ---

	r.Get("/ws", ws.Handler(func(conn ws.Conn) {
		log.Println("Client connected")

		for {
			var msg PingMessage
			if err := conn.ReadJSON(&msg); err != nil {
				log.Println("Read error:", err)
				break
			}

			log.Printf("Received: %+v\n", msg)

			// Echo back a response
			resp := PingMessage{
				Type: "pong",
				Time: time.Now().Format(time.RFC3339),
			}
			if err := conn.WriteJSON(resp); err != nil {
				log.Println("Write error:", err)
				break
			}
		}

		log.Println("Client disconnected")
	}))

	// ---

	err = server.Run("8000")
	if err != nil {
		log.Println("Error when running the server: ", err)
	}
}
