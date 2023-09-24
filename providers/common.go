package providers

import (
	"fmt"
	"os"
)

func ProvideEnvironmentVariable[T ~string](envName string) func() (T, error) {
	return func() (T, error) {
		if token, ok := os.LookupEnv(envName); ok {
			return T(token), nil
		} else {
			return "", fmt.Errorf("environment variable %q was not found in environment variables", envName)
		}
	}
}
