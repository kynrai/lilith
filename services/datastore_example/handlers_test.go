package datastore_example

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kylelemons/godebug/pretty"
	h "github.com/kynrai/lilith/server/http"
)

func TestGet(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		id     string
		getter Getter
		want   *Thing
		err    *h.HTTPError
	}{
		{
			name: "happy path",
			id:   "1",
			getter: GetterFunc(func(ctx context.Context, id string) (*Thing, error) {
				return &Thing{Id: "1", Name: "test"}, nil
			}),
			want: &Thing{Id: "1", Name: "test"},
		},
		{
			name: "error path",
			id:   "1",
			getter: GetterFunc(func(ctx context.Context, id string) (*Thing, error) {
				return nil, errors.New("boom")
			}),
			err: &h.HTTPError{Code: 500},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := mux.NewRouter()
			r.Handle("/thing/{id}", GetThing(tc.getter))
			req := httptest.NewRequest(http.MethodGet, "https://example.com/thing/"+tc.id, nil)
			rw := httptest.NewRecorder()
			r.ServeHTTP(rw, req)
			if tc.want != nil {
				got := &Thing{}
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
		putter Putter
		body   *Thing
		want   *Thing
		err    *h.HTTPError
	}{
		{
			name: "happy path",
			body: &Thing{Id: "1", Name: "test"},
			putter: PutterFunc(func(ctx context.Context, t *Thing) error {
				return nil
			}),
			want: &Thing{Id: "1", Name: "test"},
		},
		{
			name: "error path",
			body: &Thing{Id: "1", Name: "test"},
			putter: PutterFunc(func(ctx context.Context, t *Thing) error {
				return errors.New("boom!")
			}),
			err: &h.HTTPError{Code: 500},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := mux.NewRouter()
			r.Handle("/thing", PutThing(tc.putter))
			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.body); err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "https://example.com/thing", b)
			rw := httptest.NewRecorder()
			r.ServeHTTP(rw, req)
			if tc.want != nil {
				got := &Thing{}
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
