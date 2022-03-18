package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

func MapDefault(m map[string]string, key string, defaultValue string) string {
	res := m[key]
	if res == "" {
		res = defaultValue
	}
	return res
}

func EnvDefault(envKey string, defaultValue string) string {
	res := os.Getenv(envKey)
	if res == "" {
		res = defaultValue
	}
	return res
}

func EnvRequired(envKey string) string {
	res := os.Getenv(envKey)
	if res == "" {
		logrus.Errorf("%s is empty", envKey)
	}
	return res
}