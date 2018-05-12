package postgres_example

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx"
	"github.com/kynrai/lilith/services/postgres_example/models"
	"github.com/kynrai/lilith/storage/postgres"
)

var _ Repo = (*repo)(nil)

type Repo interface {
	Getter
	Setter
}

type (
	Getter interface {
		Get(ctx context.Context, id string) (*models.Thing, error)
	}
	GetterFunc func(ctx context.Context, id string) (*models.Thing, error)
)

func (f GetterFunc) Get(ctx context.Context, id string) (*models.Thing, error) {
	return f(ctx, id)
}

type (
	Setter interface {
		Set(ctx context.Context, t *models.Thing) error
	}
	SetterFunc func(ctx context.Context, t *models.Thing) error
)

func (f SetterFunc) Set(ctx context.Context, t *models.Thing) error {
	return f(ctx, t)
}

type repo struct {
	db postgres.Repo
}

func (r *repo) Get(ctx context.Context, id string) (t *models.Thing, err error) {
	return t, r.db.Conn(ctx, func(db *pgx.ConnPool) error {
		return db.QueryRow(getThing, id).Scan(t)
	})
}

func (r *repo) Set(ctx context.Context, t *models.Thing) error {
	tBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return r.db.Exec(ctx, setThing, t.ID, string(tBytes))
}
