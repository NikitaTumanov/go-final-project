package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/api"
)

const (
	defaultPort = 7540
	webDir      = "./web"
)

func Run() {
	http.Handle("/", http.FileServer(http.Dir(webDir)))

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

	api.Init()

	addr := fmt.Sprintf(":%d", port)

	log.Println("Server started on port", addr)
	http.ListenAndServe(addr, nil)
}
