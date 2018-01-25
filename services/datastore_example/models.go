package datastore_example

const kind = "datastore_example_kind"

type Thing struct {
	ID   string `datastore:"id" json:"id"`
	Name string `datastore:"name" json:"name"`
}
