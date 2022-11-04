package env

import (
	"fmt"
	"os"
)

type EnvVarMissing string

func (err EnvVarMissing) Error() string {
	return fmt.Sprintf("Missing env variable %q", string(err))
}

type EnvGetter struct{ Prefix string }

func (e *EnvGetter) Get(key string, fallback *string) (string, error) {
	key = e.Prefix + key
	val, found := os.LookupEnv(key)
	if !found {
		if fallback != nil {
			return *fallback, nil
		}
		return "", EnvVarMissing(key)
	} else {
		return val, nil
	}
}
