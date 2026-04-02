package server

import (
	"log"
	"net/http"
)

func Run() {
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
