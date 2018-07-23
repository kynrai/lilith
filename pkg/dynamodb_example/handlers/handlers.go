package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	h "github.com/kynrai/lilith/internal/http"
	"github.com/kynrai/lilith/pkg/dynamodb_example"
	"github.com/kynrai/lilith/pkg/dynamodb_example/models"
)

func GetThing(g dynamodb_example.Getter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		t, err := g.Get(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		if err := json.NewEncoder(w).Encode(t); err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		return nil
	}
}

func GetThings() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func PutThing(p dynamodb_example.Putter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		t := models.Thing{}
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			return h.HTTPError{Code: http.StatusBadRequest, Err: err}
		}
		if err := p.Put(r.Context(), &t); err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		return nil
	}
}

func PutThings() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func DelThing() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func DelThings() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
