package datastore_example

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Putter interface {
		Put(ctx context.Context, t *Thing) error
	}
	PutterFunc func(ctx context.Context, t *Thing) error
)

func (f PutterFunc) Put(ctx context.Context, t *Thing) error {
	return f(ctx, t)
}

type repo struct {
}

func New() *repo {
	return &repo{}
}

func (r *repo) Get(ctx context.Context, id string) (*Thing, error) {
	return nil, nil
}

func (r *repo) Put(ctx context.Context, t *Thing) error {
	return nil
}

func projectID() (string, error) {
	// Datastore needs the project ID when running outside of AppEngine. We can get it from the metadata.
	req, err := http.NewRequest(http.MethodGet, "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", fmt.Errorf("got status %d", resp.StatusCode)
}
