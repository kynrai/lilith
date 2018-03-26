package models

type Thing struct {
	ID   string `datastore:"id" json:"id"`
	Name string `datastore:"name" json:"name"`
}
