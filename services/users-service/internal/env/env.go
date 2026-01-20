package env

import (
	"fmt"
	"os"
)

type ConfigLoader struct {
	err error
}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

func (cl *ConfigLoader) GetEnv(key string) string {
	if cl.err != nil {
		return ""
	}

	val, err := GetEnv(key)
	if err != nil {
		cl.err = fmt.Errorf("failed to load %s: %w", key, err)
		return ""
	}
	return val
}

func (cl *ConfigLoader) Error() error {
	return cl.err
}

func GetEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}

	return value, nil
}
