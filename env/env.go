package env

import (
	"fmt"
	"os"
)

type EnvVarMissing string

func (err EnvVarMissing) Error() string {
	return fmt.Sprintf("Missing env variable %q", string(err))
}

func Get(key string) (string, error) {
	val, found := os.LookupEnv(key)
	if !found {
		return "", EnvVarMissing(key)
	} else {
		return val, nil
	}
}
