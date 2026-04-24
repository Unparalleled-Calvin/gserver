package server

import (
	"net/http"

	"github.com/Unparalleled-Calvin/gserver/internal/storage"
)

func Run(serverAddr string) error {

	// register handlers
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/register", registerHandler)

	_, _ = storage.GetDB(), storage.GetRedisClient()
	return http.ListenAndServe(serverAddr, nil)
}
