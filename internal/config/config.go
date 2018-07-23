package config

type Vars struct {
	Secret string
}

const (
	SearchBackendAWS = "aws"
)

func Load() (*Vars, error) {
	vars := &Vars{}
	err := SetRequiredVars([]Var{
		{Variable: "OPTIONAL_SECRET", Value: &vars.Secret, Optional: true},
		{Variable: "SECRET", Value: &vars.Secret},
	})
	if err != nil {
		return nil, err
	}
	return vars, nil
}
