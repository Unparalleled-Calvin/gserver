package server

import (
	"net/http"
)

func Run(serverAddr string) error {
	http.HandleFunc("/", helloHandler)
	return http.ListenAndServe(serverAddr, nil)
}
