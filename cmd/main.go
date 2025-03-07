package main

import (
	"ShorteningURL/internal/database"
	"ShorteningURL/internal/handler"
)

func main() {

	// Get pool *handler.Pool for
	// sending requests to the DB
	pool := &handler.Pool{Pool: database.Connection()}

	go pool.Ticker()

	pool.StartServer()

}
