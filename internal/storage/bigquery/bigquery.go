package bigquery

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// set at build time or set from init function based on environment
var bqProject string

func New(ctx context.Context) (*bigquery.Client, error) {
	return bigquery.NewClient(ctx, bqProject)
}

func MustNew(ctx context.Context) *bigquery.Client {
	bq, err := New(ctx)
	if err != nil {
		log.Fatalf("could not create bigquery client: %v", err)
	}
	return bq
}

// Create will create a dataset and table if it does not exist
func Create(ctx context.Context, dataset, table string, data interface{}) error {
	bq := MustNew(ctx)
	defer bq.Close()
	dsit := bq.Datasets(ctx)
	// check if dataset exists
	dsExists := false
	for {
		ds, err := dsit.Next()
		if err == iterator.Done {
			break
		}
		if ds.DatasetID == dataset {
			dsExists = true
			break
		}
	}
	// create dataset if it does not exist
	if !dsExists {
		meta := &bigquery.DatasetMetadata{Location: "EU"}
		if err := bq.Dataset(dataset).Create(ctx, meta); err != nil {
			return fmt.Errorf("failed to create dataset for bigquery dataset %s: %v", dataset, err)
		}
	}
	tit := bq.Dataset(dataset).Tables(ctx)
	tExists := false
	for {
		t, err := tit.Next()
		if err == iterator.Done {
			break
		}
		if t.TableID == table {
			tExists = true
			break
		}
	}
	if !tExists {
		schema, err := bigquery.InferSchema(data)
		if err != nil {
			return fmt.Errorf("failed to infer schema for bigquery table %s: %v", table, err)
		}
		err = bq.Dataset(dataset).Table(table).Create(context.Background(), &bigquery.TableMetadata{Schema: schema})
		if err != nil {
			return fmt.Errorf("failed to create table for bigquery destiation: %v", err)
		}
	}
	return nil
}

// Delete deletes a table from a dataset
func Delete(ctx context.Context, dataset, table string) error {
	bq := MustNew(ctx)
	defer bq.Close()
	return bq.Dataset(dataset).Table(table).Delete(ctx)
}

type Bulk struct {
	client    *bigquery.Client
	dataset   string
	table     string
	BatchSize int
	// If true will delete the table before streaming
	Delete bool
}

func NewBulk(dataset, table string) *Bulk {
	b := new(Bulk)
	b.client = MustNew(context.Background())
	b.dataset = dataset
	b.table = table
	b.BatchSize = 500
	return b
}

func (b *Bulk) Stream(ctx context.Context, data interface{}) error {
	v := reflect.ValueOf(data)
	if k := v.Kind(); k != reflect.Slice {
		return fmt.Errorf("must use a slice with Stream, not a %s", k.String())
	}

	if b.Delete {
		if err := b.client.Dataset(b.dataset).Table(b.table).Delete(ctx); err != nil {
			return err
		}
	}
	if err := Create(context.Background(), b.dataset, b.table, reflect.New(reflect.TypeOf(data).Elem()).Interface()); err != nil {
		return err
	}
	up := b.client.Dataset(b.dataset).Table(b.table).Uploader()

	callCount := (v.Len()-1)/b.BatchSize + 1
	var wg sync.WaitGroup
	wg.Add(callCount)
	for i := 0; i < callCount; i++ {
		lo := i * b.BatchSize
		hi := (i + 1) * b.BatchSize
		if hi > v.Len() {
			hi = v.Len()
		}
		log.Printf("loading from %d to %d\n", lo, hi)
		go func(vals reflect.Value) {
			defer wg.Done()
			if err := up.Put(ctx, vals.Interface()); err != nil {
				if pmErr, ok := err.(bigquery.PutMultiError); ok {
					for _, rowInsertionError := range pmErr {
						log.Fatal(rowInsertionError.Errors)
					}
				}
			}
		}(v.Slice(lo, hi))
	}
	wg.Wait()

	return nil
}
