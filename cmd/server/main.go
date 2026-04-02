package main

import (
	"log"
	"os"

	"github.com/Unparalleled-Calvin/gserver/internal/server"
)

func main() {
	serverAddr := os.Getenv("GSERVER_ADDR")

	if serverAddr == "" {
		serverAddr = "localhost:8000"
	}

	log.Fatal(server.Run(serverAddr))
}
