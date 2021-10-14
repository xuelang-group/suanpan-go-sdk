package stream

import (
	"strconv"
	"sync"

	"github.com/golang/glog"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/mq"
)

type EnvStream struct {
	StreamSendQueue                string `mapstructure:"--stream-send-queue"`
	StreamRecvQueue                string `mapstructure:"--stream-recv-queue" default:"mq-master"`
	StreamSendQueueMaxLength       string `mapstructure:"--stream-send-queue-max-length" default:"1000"`
	StreamSendQueueTrimImmediately string `mapstructure:"--stream-send-queue-trim-immediately" default:"False"`
}

type Stream struct {
	StreamSendQueue                string
	StreamRecvQueue                string
	StreamSendQueueMaxLength       int64
	StreamSendQueueTrimImmediately bool
}

var (
	s          *Stream
	streamOnce sync.Once
)

func getStream() *Stream {
	streamOnce.Do(func() {
		s = buildStream()
	})

	return s
}

func buildStream() *Stream {
	argsMap := config.GetArgs()
	var envStream EnvStream
	mapstructure.Decode(argsMap, &envStream)
	defaults.SetDefaults(envStream)

	sendQueue := envStream.StreamSendQueue
	if sendQueue == "" {
		glog.Warning("StreamSendQueue is empty")
		sendQueue = "mq-" + config.GetEnv().SpNodeId
	}

	maxLen, err := strconv.ParseInt(envStream.StreamSendQueueMaxLength, 10, 64)
	if err != nil {
		glog.Errorf("StreamSendQueueMaxLength is not a valid int64 value: %s", envStream.StreamSendQueueMaxLength)
		maxLen = 1000
	}

	trimImmediately, err := strconv.ParseBool(envStream.StreamSendQueueTrimImmediately)
	if err != nil {
		glog.Errorf("StreamSendQueueTrimImmediately is not a valid bool value: %s", envStream.StreamSendQueueTrimImmediately)
		trimImmediately = false
	}

	return &Stream{
		StreamSendQueue:                sendQueue,
		StreamRecvQueue:                envStream.StreamRecvQueue,
		StreamSendQueueMaxLength:       maxLen,
		StreamSendQueueTrimImmediately: trimImmediately,
	}
}

func SendMessage(data interface{}) string {
	s := getStream()
	return s.sendMessage(data)
}

func SubscribeQueue() <-chan interface{} {
	s := getStream()
	return s.subscribeQueue()
}

func (s *Stream) sendMessage(data interface{}) string {
	q := mq.GetMq()
	return q.SendMessage(s.StreamSendQueue, data, s.StreamSendQueueMaxLength, s.StreamSendQueueTrimImmediately)
}

func (s *Stream) subscribeQueue() <-chan interface{} {
	q := mq.GetMq()
	group := config.GetEnv().SpNodeGroup
	consumer := config.GetEnv().SpNodeId
	return q.SubscribeQueue(s.StreamRecvQueue, group, consumer)
}