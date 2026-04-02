package main

import (
	"log"

	"github.com/Unparalleled-Calvin/gserver/internal/client"
)

func main() {
	if err := client.Run("http://localhost:8000/"); err != nil {
		log.Fatal(err)
	}
}
