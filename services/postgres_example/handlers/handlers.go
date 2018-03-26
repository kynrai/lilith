package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/kynrai/lilith/server/http"
	"github.com/kynrai/lilith/services/postgres_example"
	"github.com/kynrai/lilith/services/postgres_example/models"
)

func GetThing(g postgres_example.Getter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		thing, err := g.Get(r.Context(), mux.Vars(r)["id"])
		if err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		if err := json.NewEncoder(w).Encode(thing); err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		return nil
	}
}

func PutThing(g postgres_example.Setter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var thing models.Thing
		if err := json.NewDecoder(r.Body).Decode(&thing); err != nil {
			return h.HTTPError{Code: http.StatusBadRequest, Err: err}
		}
		if err := g.Set(r.Context(), &thing); err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		if err := json.NewEncoder(w).Encode(thing); err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		return nil
	}
}
