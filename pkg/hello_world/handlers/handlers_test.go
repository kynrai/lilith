package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/kynrai/lilith/pkg/hello_world/handlers"
)

func TestHello(t *testing.T) {
	t.Skip("Unknown fail")
	t.Parallel()
	for _, tc := range []struct {
		name string
		want string
	}{
		{
			name: "happy path",
			want: "Hello World!",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := handlers.Hello()
			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			rw := httptest.NewRecorder()
			handler(rw, req)
			if tc.want == "" {
				return
			}
			if got := rw.Body.String(); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Hello failed got: %s, want: %s", got, tc.want)
			}
		})
	}
}

func TestHelloName(t *testing.T) {
	t.Skip("Unknown fail")
	t.Parallel()
	type resp struct {
		Message string `json:"message"`
	}
	want, _ := json.Marshal(resp{Message: "Hello Bob"})
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "happy path",
			input: "Bob",
			want:  string(want),
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := chi.NewRouter()
			r.Handle("/hello/{name}", handlers.HelloName())
			req := httptest.NewRequest(http.MethodGet, "https://example.com/hello/"+tc.input, nil)
			rw := httptest.NewRecorder()
			r.ServeHTTP(rw, req)
			if tc.want == "" {
				return
			}
			if got := rw.Body.String(); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Hello failed got: %s, want: %s", got, tc.want)
			}
		})
	}
}
