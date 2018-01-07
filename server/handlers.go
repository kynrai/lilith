package server

import (
	"net/http"

	h "github.com/kynrai/lilith/server/http"
)

func Health() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("OK"))
		return nil
	}
}
