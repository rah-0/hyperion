package util

import (
	"os"
)

func GetEnvKeyValue(key string) string {
	value, _ := os.LookupEnv(key)
	return value
}
