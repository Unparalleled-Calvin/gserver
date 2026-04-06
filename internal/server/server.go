package server

import (
	"net/http"

	"github.com/Unparalleled-Calvin/gserver/internal/storage"
)

func Run(serverAddr string) error {
	http.HandleFunc("/", helloHandler)
	_, _ = storage.GetDB(), storage.GetRedisClient()
	return http.ListenAndServe(serverAddr, nil)
}
