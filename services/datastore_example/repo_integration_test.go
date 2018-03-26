package datastore_example

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/kynrai/lilith/services/datastore_example/models"
)

const emulator = "localhost:8081"

func sethost() {
	host := os.Getenv("DATASTORE_EMULATOR_HOST")
	if host == "" {
		os.Setenv("DATASTORE_EMULATOR_HOST", emulator)
	}
}

func TestPutGet_Integration(t *testing.T) {
	t.Parallel()
	sethost()
	repo := New()
	for _, tc := range []struct {
		name string
		id   string
		body *models.Thing
		want *models.Thing
	}{
		{
			name: "happy path",
			id:   "1",
			body: &models.Thing{ID: "1", Name: "test"},
			want: &models.Thing{ID: "1", Name: "test"},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := repo.Put(context.Background(), tc.body); err != nil {
				t.Fatal(err)
			}
			thing, err := repo.Get(context.Background(), tc.id)
			if err != nil {
				t.Fatal(err)
			}
			if tc.want != nil {
				if !reflect.DeepEqual(thing, tc.want) {
					t.Fatal(pretty.Compare(thing, tc.want))
				}
			}
		})
	}
}
