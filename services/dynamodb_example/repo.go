package dynamodb_example

import (
	"context"

	"github.com/kynrai/lilith/storage/dynamodb"
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
	db dynamodb.Repo
}

func New(db dynamodb.Repo) *repo {
	return &repo{db}
}

func (r *repo) Get(ctx context.Context, id string) (*Thing, error) {
	var key struct {
		ID string
	}
	key.ID = id
	t := &Thing{}
	return t, r.db.Get(ctx, tableName, &key, t)
}

func (r *repo) Put(ctx context.Context, t *Thing) error {
	return r.db.Put(ctx, tableName, t)
}
