package server

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) { // response hello whatever the request is
	fmt.Fprintln(w, "Hello!")
}
