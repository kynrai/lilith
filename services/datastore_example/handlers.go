package datastore_example

import (
	"net/http"

	h "github.com/kynrai/lilith/server/http"
)

func GetThing(g Getter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func PutThing(p Putter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
