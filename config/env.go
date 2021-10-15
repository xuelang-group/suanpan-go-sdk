package config

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-defaults"
)

var (
	e = &Env{}
	envOnce sync.Once
)

type Env struct{
	SpParam		string		`envconfig:"SP_PARAM"`
	SpNodeId	string		`envconfig:"SP_NODE_ID" validate:"required"`
	SpNodeGroup	string		`envconfig:"SP_NODE_GROUP" default:"default"`
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
		glog.Errorf("Decode env variables failed: %v", err)
	}

	defaults.SetDefaults(e)
	validate := validator.New()
	err = validate.Struct(e)
	if err != nil {
		glog.Errorf("Validate env variables failed: %v", err)
	}

	return e
}