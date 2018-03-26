package handlers_test

import "testing"

func TestGetThing(t *testing.T) {
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
		})
	}
}

func TestPutThing(t *testing.T) {
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
		})
	}
}
