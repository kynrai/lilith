package datastore_example

const kind = "datastore_example_kind"

type Thing struct {
	Id   string `datastore:"id" json:"id"`
	Name string `datastore:"id" json:"name"`
}
