package env

import (
	"fmt"
	"os"
)

func GetEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}

	return value
}
