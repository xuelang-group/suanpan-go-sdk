package mq

import (
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

const (
	Redis = "redis"
)

type Mq interface {
	SendMessage(queue string, data map[string]string, maxLen int64, trimImmediately bool)	string
	SubscribeQueue(queue, group, consumer string) <-chan map[string]interface{}
}

type EnvMq struct {
	MqType	string	`mapstructure:"--mq-type" default:"redis"`
}

func New(argsMap map[string]string) Mq {
	var envMq EnvMq
	mapstructure.Decode(argsMap, &envMq)
	defaults.SetDefaults(&envMq)
	switch envMq.MqType {
	case Redis:
		var redisMq RedisMq
		mapstructure.Decode(argsMap, &redisMq)
		defaults.SetDefaults(&redisMq)
		redisMq.initClient()
		return &redisMq
	default:
		log.Errorf("Unsupported mq type: %s", envMq.MqType)
		return nil
	}
}