package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	defaultPort = 7540
	webDir      = "./web"
)

func StartServer() {
	router := chi.NewRouter()

	router.Handle("/*", http.FileServer(http.Dir(webDir)))

	port := defaultPort

	envPort := os.Getenv("TODO_PORT")
	if envPort != "" {
		envPortInt, err := strconv.Atoi(envPort)
		if err == nil {
			port = envPortInt
		} else {
			log.Printf("invalid TODO_PORT: %s", envPort)
		}
	}

	addr := fmt.Sprintf(":%d", port)

	log.Println("Server started on port", addr)
	http.ListenAndServe(addr, router)
}
