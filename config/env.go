package config

import (
	"fmt"
	"os"
	"strings"
)

type ErrorSlice struct {
	errs []string
}

func (e ErrorSlice) Error() string {
	return strings.Join(e.errs, ", ")
}

type Var struct {
	Variable string
	Value    *string
	Optional bool
}

// SetRequiredVars given a slice of Var, will load the environment variable and set the value
// returning errors for any env variable not found
func SetRequiredVars(required []Var) error {
	var errs []string
	for _, v := range required {
		if *v.Value = os.Getenv(v.Variable); *v.Value == "" && !v.Optional {
			errs = append(errs, fmt.Sprintf("%s not set", v.Variable))
		}
	}
	if len(errs) > 0 {
		return ErrorSlice{errs}
	}
	return nil
}
