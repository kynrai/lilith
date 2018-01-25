package datastore_example

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

var _ Repo = (*repo)(nil)

type Repo interface {
	Getter
	Putter
}

type (
	Getter interface {
		Get(ctx context.Context, id string) (*Thing, error)
	}
	GetterFunc func(ctx context.Context, id string) (*Thing, error)
)

func (f GetterFunc) Get(ctx context.Context, id string) (*Thing, error) {
	return f(ctx, id)
}

type (
	MultiGetter interface {
		GetMulti(ctx context.Context, ids ...string) ([]*Thing, error)
	}
	MultiGetterFunc func(ctx context.Context, ids ...string) ([]*Thing, error)
)

func (f MultiGetterFunc) GetMulti(ctx context.Context, ids ...string) ([]*Thing, error) {
	return f(ctx, ids...)
}

type (
	Putter interface {
		Put(ctx context.Context, t *Thing) error
	}
	PutterFunc func(ctx context.Context, t *Thing) error
)

func (f PutterFunc) Put(ctx context.Context, t *Thing) error {
	return f(ctx, t)
}

type (
	MultiPutter interface {
		PutMulti(ctx context.Context, ts ...*Thing) error
	}
	MultiPutterFunc func(ctx context.Context, ts ...*Thing) error
)

func (f MultiPutterFunc) PutMulti(ctx context.Context, ts ...*Thing) error {
	return f(ctx, ts...)
}

type repo struct {
	ds *datastore.Client
}

func New() *repo {
	ds, err := datastore.NewClient(context.Background(), projectID())
	if err != nil {
		// This repo wont work without a client so we can fatal here
		log.Fatal(err)
	}
	return &repo{ds}
}

func (r *repo) Get(ctx context.Context, id string) (*Thing, error) {
	t := &Thing{}
	return t, r.ds.Get(ctx, datastore.NameKey(kind, id, nil), t)
}

func (r *repo) Put(ctx context.Context, t *Thing) error {
	_, err := r.ds.Put(ctx, datastore.NameKey(kind, t.ID, nil), t)
	return err
}

// ProjectID will attempt to get the Google Cloud Project ID with the following rules:
// 1) Look for the DATASTORE_PROJECT_ID envar
// 2) Use the metadata API to get ID, this will only work in Google Cloud
// 3) Any failure or timeout (3s) will presume that the code is running outside the cloud
// in which case a default project ID is returned.
func projectID() string {
	if id := os.Getenv("DATASTORE_PROJECT_ID"); id != "" {
		return id
	}
	const defaultID = "project-id"
	req, err := http.NewRequest(http.MethodGet, "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	if err != nil {
		return defaultID
	}
	req.Header.Add("Metadata-Flavor", "Google")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return defaultID
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return defaultID
		}
		return string(b)
	}
	return defaultID
}
