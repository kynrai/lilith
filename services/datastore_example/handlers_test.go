package datastore_example

import "testing"

func TestGet(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
	}{
		{
			name: "happy path",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			id, err := projectID()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(id)
		})
	}
}
