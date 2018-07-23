package server

import (
	"encoding/json"
	"net/http"

	h "github.com/kynrai/lilith/internal/http"
)

func Health() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		type resp struct {
			Health string `json:"health"`
		}
		return json.NewEncoder(w).Encode(
			resp{
				Health: "OK",
			})
	}
}
