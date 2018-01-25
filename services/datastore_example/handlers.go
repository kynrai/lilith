package datastore_example

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

func PutThing(p Putter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		t := Thing{}
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			return h.HTTPError{Code: http.StatusBadRequest, Err: err}
		}
		if err := p.Put(r.Context(), &t); err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		if err := json.NewEncoder(w).Encode(&t); err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		return nil
	}
}
