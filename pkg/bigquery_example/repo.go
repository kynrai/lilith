package bigquery_example

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	bqi "github.com/kynrai/lilith/internal/storage/bigquery"
	"github.com/kynrai/lilith/pkg/bigquery_example/models"
)

var (
	bqDataset = "demo"
	bqTable   = "example"
)

var _ Repo = (*repo)(nil)

type Repo interface {
	Setter
}

type (
	Setter interface {
		Set(ctx context.Context, p []*models.Thing) error
	}
	SetterFunc func(ctx context.Context, p []*models.Thing) error
)

func (f SetterFunc) Set(ctx context.Context, p []*models.Thing) error {
	return f(ctx, p)
}

type repo struct {
	bq *bigquery.Client
	u  *bigquery.Uploader
}

func new(ds *datastore.Client, bq *bigquery.Client) *repo {
	// create the table
	if err := bqi.Create(context.Background(), bqDataset, bqTable, models.ThingBQSchema{}); err != nil {
		log.Fatalf("failed to create bigquery table for example: %v", err)
	}
	return &repo{
		bq: bq,
		u:  bq.Dataset(bqDataset).Table(bqTable).Uploader(),
	}
}

func (r *repo) Set(ctx context.Context, ds []*models.Thing) error {
	err := r.u.Put(ctx, ds)
	if pmErr, ok := err.(bigquery.PutMultiError); ok {
		for _, rowInsertionError := range pmErr {
			log.Println(rowInsertionError.Errors)
		}
	}
	return err

}

func BulkUploader(hugeData []*models.Thing) error {
	bqBulk := bqi.NewBulk("dataset", "table")
	// Deletes the table first then recreates
	bqBulk.Delete = true
	// Max batch size as the data cannot be more than 10mb per call
	bqBulk.BatchSize = 8000
	// Concurrently insert all batches
	if err := bqBulk.Stream(context.Background(), hugeData); err != nil {
		return err
	}
	return nil
}
