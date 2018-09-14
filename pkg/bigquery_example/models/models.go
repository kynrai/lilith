package models

import (
	"cloud.google.com/go/bigquery"
)

type Thing struct {
	Name          string  `json:"name"`
	OptionalField *string `json:"optional_field,omitempty"`
}

type ThingBQSchema struct {
	Name          string              `bigquery:"name"`
	OptionalField bigquery.NullString `bigquery:"optional_field,omitempty"`
}

// Save implements the ValueSaver interface.
// Name used as dedupe key
func (t *ThingBQSchema) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"name":           t.Name,
		"optional_field": t.OptionalField,
	}, t.Name, nil
}
