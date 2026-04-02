package main

import (
	"log"

	"calvin.com/gserver/internal/client"
)

func main() {
	if err := client.Run("http://localhost:8000/"); err != nil {
		log.Fatal(err)
	}
}
