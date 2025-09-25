package main

import (
	"log"

	"github.com/V4T54L/goship/pkg/goship"
)

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
	log.Println("Hi", db, err)

	server := goship.NewChiServer()
	err = server.Run("8000")
	if err != nil {
		log.Println("Error when running the server: ", err)
	}
}
