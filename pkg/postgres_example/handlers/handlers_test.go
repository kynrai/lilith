package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/kylelemons/godebug/pretty"
	h "github.com/kynrai/lilith/internal/http"
	"github.com/kynrai/lilith/pkg/postgres_example"
	"github.com/kynrai/lilith/pkg/postgres_example/handlers"
	"github.com/kynrai/lilith/pkg/postgres_example/models"
)

func TestGetThing(t *testing.T) {
	thing1 := &models.Thing{ID: "1"}
	t.Parallel()
	for _, tc := range []struct {
		name   string
		id     string
		getter postgres_example.Getter
		want   *models.Thing
		err    *h.HTTPError
	}{
		{
			name: "happy path",
			id:   "1",
			getter: postgres_example.GetterFunc(func(ctx context.Context, id string) (*models.Thing, error) {
				if id != "1" {
					return nil, nil
				}
				return thing1, nil
			}),
			want: thing1,
		},
		{
			name: "error path",
			id:   "1",
			getter: postgres_example.GetterFunc(func(ctx context.Context, id string) (*models.Thing, error) {
				return nil, h.HTTPError{Code: 500}
			}),
			err: &h.HTTPError{Code: 500},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := chi.NewRouter()
			r.Handle("/thing/{id}", handlers.GetThing(tc.getter))
			req := httptest.NewRequest(http.MethodGet, "https://example.com/thing/"+tc.id, nil)
			rw := httptest.NewRecorder()
			r.ServeHTTP(rw, req)
			if tc.err != nil && tc.err.Code != rw.Code {
				t.Fatalf("status codes dont match, Got: %d, Want: %d", rw.Code, tc.err.Code)
			}
			if tc.want == nil {
				return
			}
			if rw.Body.String() == "" {
				t.Fatal("GetThing failed, body should not be empty")
			}
			var got models.Thing
			if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
				t.Fatal("GetThing failed, failed to decode response")
			}
			if !reflect.DeepEqual(&got, tc.want) {
				t.Fatal(pretty.Compare(got, tc.want))
			}
		})
	}
}

func TestSetThing(t *testing.T) {
	thing1 := &models.Thing{ID: "1"}
	t.Parallel()
	for _, tc := range []struct {
		name   string
		id     string
		setter postgres_example.Setter
		thing  string
		want   *models.Thing
		err    *h.HTTPError
	}{
		{
			name:  "happy path",
			thing: `{"id":"1"}`,
			want:  thing1,
			setter: postgres_example.SetterFunc(func(ctx context.Context, t *models.Thing) error {
				return nil
			}),
		},
		{
			name:  "decode fail",
			thing: "boom",
			err:   &h.HTTPError{Code: 400},
			setter: postgres_example.SetterFunc(func(ctx context.Context, t *models.Thing) error {
				return nil
			}),
		},
		{
			name:  "error path",
			thing: `{"id":"1"}`,
			setter: postgres_example.SetterFunc(func(ctx context.Context, t *models.Thing) error {
				return errors.New("boom")
			}),
			err: &h.HTTPError{Code: http.StatusInternalServerError},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := chi.NewRouter()
			r.Handle("/thing", handlers.PutThing(tc.setter))
			req := httptest.NewRequest(http.MethodPost, "/thing", bytes.NewBufferString(tc.thing))
			rw := httptest.NewRecorder()
			r.ServeHTTP(rw, req)
			if tc.err != nil && tc.err.Code != rw.Code {
				t.Fatalf("status codes dont match, Got: %d, Want: %d", rw.Code, tc.err.Code)
			}
			if tc.want == nil {
				return
			}
			if rw.Body.String() == "" {
				t.Fatal("GetThing failed, body should not be empty")
			}
			var got models.Thing
			if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
				t.Fatal("GetThing failed, failed to decode response")
			}
			if !reflect.DeepEqual(&got, tc.want) {
				t.Fatal(pretty.Compare(got, tc.want))
			}
		})
	}
}
