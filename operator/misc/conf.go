package misc

import (
	"os"
	"strings"
)

const SA_TOKEN_PATH = "/var/run/secrets/kubernetes.io/serviceaccount"

func GetEnvValue(key string, def string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return def
	}

	return val
}

func GetEnvBool(key string, def bool) bool {
	val := os.Getenv(key)

	if len(val) == 0 {
		return def
	}

	if strings.Compare("false", strings.ToLower(val)) == 0 {
		return false
	}

	if strings.Compare("true", strings.ToLower(val)) == 0 {
		return true
	}

	return def
}

func GetInClusterK8sToken() (string, error) {
	_, err := os.Stat(SA_TOKEN_PATH)
	if err != nil {
		return "", err
	}

	dat, err := os.ReadFile(SA_TOKEN_PATH)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}