package mq

import (
	"github.com/golang/glog"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/xuelang-group/suanpan-go-sdk/config"
)

const (
	Redis = "redis"
)

type Mq interface {
	SendMessage(queue string, data interface{}, maxLen int64, trimImmediately bool)	string
	SubscribeQueue(queue, group, consumer string) <-chan interface{}
}

type EnvMq struct {
	MqType	string	`mapstructure:"--mq-type" default:"redis"`
}

func GetMq() Mq {
	argsMap := config.GetArgs()
	var envMq EnvMq
	mapstructure.Decode(argsMap, &envMq)
	defaults.SetDefaults(envMq)
	switch envMq.MqType {
	case Redis:
		var redisMq RedisMq
		mapstructure.Decode(argsMap, redisMq)
		defaults.SetDefaults(redisMq)
		return &redisMq
	default:
		glog.Errorf("Unsupported mq type: %s", envMq.MqType)
		return nil
	}
}