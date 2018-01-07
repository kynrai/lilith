package hello_world_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"

	"github.com/kynrai/lilith/services/hello_world"
)

func TestHello(t *testing.T) {
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
			handler := hello_world.Hello()
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
	t.Parallel()
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "happy path",
			input: "Bob",
			want:  "Hello Bob",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := mux.NewRouter()
			r.Handle("/hello/{name}", hello_world.HelloName())
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
