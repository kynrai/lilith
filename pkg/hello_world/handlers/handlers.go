package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	h "github.com/kynrai/lilith/internal/http"
)

func Hello() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		type resp struct {
			Message string `json:"message"`
		}
		return json.NewEncoder(w).Encode(resp{Message: "Hello World!"})
	}
}

func HelloName() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		name := chi.URLParam(r, "name")
		type resp struct {
			Message string `json:"message"`
		}
		return json.NewEncoder(w).Encode(resp{Message: fmt.Sprintf("Hello %s", name)})
	}
}
