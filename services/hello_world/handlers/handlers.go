package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/kynrai/lilith/server/http"
)

func Hello() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "Hello World!")
		return nil
	}
}

func HelloName() h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "Hello %s", mux.Vars(r)["name"])
		return nil
	}
}
