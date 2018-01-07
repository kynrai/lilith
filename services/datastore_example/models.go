package datastore_example

type Thing struct {
	id   string `datastore:"id" json:"id"`
	name string `datastore:"id" json:"name"`
}
