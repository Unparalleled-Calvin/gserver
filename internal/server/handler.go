package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Unparalleled-Calvin/gserver/internal/schema"
	"github.com/Unparalleled-Calvin/gserver/internal/service"
)

func helloHandler(w http.ResponseWriter, r *http.Request) { // response hello whatever the request is
	fmt.Fprintln(w, "Hello!")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var payload schema.UserRegister
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if payload.ID == 0 || payload.UserName == "" || payload.Password == "" {
		http.Error(w, "id, username, and password are required", http.StatusBadRequest)
		return
	}

	if err := service.RegisterUser(payload); err != nil {
		http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "ok: id=%d\n", payload.ID)
}
