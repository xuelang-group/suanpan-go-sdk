package mq

import (
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/util"
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
	envMq := newEnvMq(argsMap)
	switch envMq.MqType {
	case Redis:
		return newRedisMq(argsMap)
	default:
		log.Errorf("Unsupported mq type: %s", envMq.MqType)
		return nil
	}
}

func newEnvMq(argsMap map[string]string) *EnvMq {
	return &EnvMq{
		MqType: util.MapDefault(argsMap, "--mq-type", "redis"),
	}
}