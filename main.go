package main

import (
	"log"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/db"
	"github.com/NikitaTumanov/go-final-project/internal/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	server.StartServer()
}
