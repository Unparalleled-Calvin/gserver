package main

import (
	"log"

	"github.com/Unparalleled-Calvin/gserver/internal/server"
	"github.com/Unparalleled-Calvin/gserver/internal/settings"
)

func main() {
	settings.Load()
	log.Fatal(server.Run(settings.ServerAddr))
}
