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
	"github.com/kynrai/lilith/pkg/datastore_example"
	"github.com/kynrai/lilith/pkg/datastore_example/handlers"
	"github.com/kynrai/lilith/pkg/datastore_example/models"
)

func TestGet(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		id     string
		getter datastore_example.Getter
		want   *models.Thing
		err    *h.HTTPError
	}{
		{
			name: "happy path",
			id:   "1",
			getter: datastore_example.GetterFunc(func(ctx context.Context, id string) (*models.Thing, error) {
				return &models.Thing{ID: "1", Name: "test"}, nil
			}),
			want: &models.Thing{ID: "1", Name: "test"},
		},
		{
			name: "error path",
			id:   "1",
			getter: datastore_example.GetterFunc(func(ctx context.Context, id string) (*models.Thing, error) {
				return nil, errors.New("boom")
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
			if tc.want != nil {
				got := &models.Thing{}
				if err := json.NewDecoder(rw.Body).Decode(got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatal(pretty.Compare(got, tc.want))
				}
			}
			if tc.err != nil {
				if got := rw.Code; !reflect.DeepEqual(got, tc.err.Code) {
					t.Fatal(pretty.Compare(got, tc.err.Code))
				}
			}
		})
	}
}

func TestPut(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		putter datastore_example.Putter
		body   *models.Thing
		want   *models.Thing
		err    *h.HTTPError
	}{
		{
			name: "happy path",
			body: &models.Thing{ID: "1", Name: "test"},
			putter: datastore_example.PutterFunc(func(ctx context.Context, t *models.Thing) error {
				return nil
			}),
			want: &models.Thing{ID: "1", Name: "test"},
		},
		{
			name: "error path",
			body: &models.Thing{ID: "1", Name: "test"},
			putter: datastore_example.PutterFunc(func(ctx context.Context, t *models.Thing) error {
				return errors.New("boom!")
			}),
			err: &h.HTTPError{Code: 500},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := chi.NewRouter()
			r.Handle("/thing", handlers.PutThing(tc.putter))
			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.body); err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "https://example.com/thing", b)
			rw := httptest.NewRecorder()
			r.ServeHTTP(rw, req)
			if tc.want != nil {
				got := &models.Thing{}
				if err := json.NewDecoder(rw.Body).Decode(got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatal(pretty.Compare(got, tc.want))
				}
			}
			if tc.err != nil {
				if got := rw.Code; !reflect.DeepEqual(got, tc.err.Code) {
					t.Fatal(pretty.Compare(got, tc.err.Code))
				}
			}
		})
	}
}
