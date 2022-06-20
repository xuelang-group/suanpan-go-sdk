package config

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-defaults"
	"github.com/sirupsen/logrus"
)

var (
	e       = &Env{}
	envOnce sync.Once
)

type Env struct {
	SpParam                      string `envconfig:"SP_PARAM"`
	SpNodeId                     string `envconfig:"SP_NODE_ID" validate:"required"`
	SpNodeGroup                  string `envconfig:"SP_NODE_GROUP" default:"default"`
	SpDebug                      string `envconfig:"SP_DEBUG"`
	SpHost                       string `envconfig:"SP_HOST" validate:"required"`
	SpHostTls                    string `envconfig:"SP_HOST_TLS" default:"false"`
	SpTermPort                   string `envconfig:"SP_TERM_PORT" default:"8002"`
	SpOs                         string `envconfig:"SP_OS" default:"kubernetes"`
	SpPort                       string `envconfig:"SP_PORT" default:"7000"`
	SpUserId                     string `envconfig:"SP_USER_ID" validate:"required"`
	SpAppId                      string `envconfig:"SP_APP_ID" validate:"required"`
	SpAccessSecret               string `envconfig:"SP_ACCESS_SECRET" validate:"required"`
	SpUserIdHeaderField          string `envconfig:"SP_USER_ID_HEADER_FIELD" default:"x-sp-user-id"`
	SpUserSignatureHeaderField   string `envconfig:"SP_USER_SIGNATURE_HEADER_FIELD" default:"x-sp-signature"`
	SpUserSignVersionHeaderField string `envconfig:"SP_USER_SIGN_VERSION_HEADER_FIELD" default:"x-sp-sign-version"`
	SpLogkitUri                  string `envconfig:"SP_LOGKIT_URI"`
	SpLogkitNamespace            string `envconfig:"SP_LOGKIT_NAMESPACE" default:"/logkit"`
	SpLogkitPath                 string `envconfig:"SP_LOGKIT_PATH"`
	SpLogkitEventsAppend         string `envconfig:"SP_LOGKIT_EVENTS_APPEND" default:"append"`
	SpLogkitLogsLevel            string `envconfig:"SP_LOGKIT_LOGS_LEVEL" default:"info"`
}

func GetEnv() *Env {
	envOnce.Do(func() {
		e = buildEnv()
	})

	return e
}

func buildEnv() *Env {
	err := envconfig.Process("config", e)
	if err != nil {
		logrus.Errorf("Decode env variables failed: %v", err)
	}

	defaults.SetDefaults(e)
	validate := validator.New()
	err = validate.Struct(e)
	if err != nil {
		logrus.Errorf("Validate env variables failed: %v", err)
	}

	return e
}
