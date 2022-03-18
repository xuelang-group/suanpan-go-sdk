package config

import (
	"os"
	"sync"

	"github.com/xuelang-group/suanpan-go-sdk/util"
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

func buildEnv() *Env {
	return &Env{
		SpParam: os.Getenv("SP_PARAM"),
		SpNodeId: util.EnvRequired("SP_NODE_ID"),
		SpNodeGroup: util.EnvDefault("SP_NODE_GROUP","default"),
		SpDebug: os.Getenv("SP_DEBUG"),
		SpHost: util.EnvRequired("SP_HOST"),
		SpHostTls: util.EnvDefault("SP_HOST_TLS", "false"),
		SpOs: util.EnvDefault("SP_OS", "kubernetes"),
		SpPort: util.EnvDefault("SP_PORT", "7000"),
		SpUserId: util.EnvRequired("SP_USER_ID"),
		SpAppId: util.EnvRequired("SP_APP_ID"),
		SpAccessSecret: util.EnvRequired("SP_ACCESS_SECRET"),
		SpUserIdHeaderField: util.EnvDefault("SP_USER_ID_HEADER_FIELD", "x-sp-user-id"),
		SpUserSignatureHeaderField: util.EnvDefault("SP_USER_SIGNATURE_HEADER_FIELD", "x-sp-signature"),
		SpUserSignVersionHeaderField: util.EnvDefault("SP_USER_SIGN_VERSION_HEADER_FIELD", "x-sp-sign-version"),
		SpLogkitUri: os.Getenv("SP_LOGKIT_URI"),
		SpLogkitNamespace: util.EnvDefault("SP_LOGKIT_NAMESPACE", "/logkit"),
		SpLogkitPath: os.Getenv("SP_LOGKIT_PATH"),
		SpLogkitEventsAppend: util.EnvDefault("SP_LOGKIT_EVENTS_APPEND", "append"),
		SpLogkitLogsLevel: util.EnvDefault("SP_LOGKIT_LOGS_LEVEL", "info"),
	}
}

func GetEnv() *Env {
	envOnce.Do(func() {
		e = buildEnv()
	})

	return e
}
