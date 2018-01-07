package postgres_example

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/kynrai/lilith/server/http"
)

func GetThing(g Getter) h.ErrorHandler {
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

func PutThing(g Setter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var thing Thing
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
