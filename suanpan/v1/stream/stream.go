package stream

import (
	"regexp"
	"strconv"
	"sync"

	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/mq"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/util"
)

const (
	InputPattern = `^in(\d+)$`
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
	MessageQueue                   mq.Mq
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
	defaults.SetDefaults(&envStream)

	sendQueue := envStream.StreamSendQueue
	if sendQueue == "" {
		log.Warn("StreamSendQueue is empty")
		sendQueue = "mq-" + config.GetEnv().SpNodeId
	}

	maxLen, err := strconv.ParseInt(envStream.StreamSendQueueMaxLength, 10, 64)
	if err != nil {
		log.Errorf("StreamSendQueueMaxLength is not a valid int64 value: %s", envStream.StreamSendQueueMaxLength)
		maxLen = 1000
	}

	trimImmediately, err := strconv.ParseBool(envStream.StreamSendQueueTrimImmediately)
	if err != nil {
		log.Errorf("StreamSendQueueTrimImmediately is not a valid bool value: %s", envStream.StreamSendQueueTrimImmediately)
		trimImmediately = false
	}

	q := mq.New(config.GetArgs())

	return &Stream{
		StreamSendQueue:                sendQueue,
		StreamRecvQueue:                envStream.StreamRecvQueue,
		StreamSendQueueMaxLength:       maxLen,
		StreamSendQueueTrimImmediately: trimImmediately,
		MessageQueue:                   q,
	}
}

func Subscribe() <-chan Request {
	s := getStream()
	return s.subscribe()
}

func (r *Request) Send(data map[string]string) string {
	return r.SendSuccess(data)
}

func (r *Request) SendOutput(i int, data string) string {
	return r.Send(map[string]string{
		OutputDataPrefix + strconv.Itoa(i): data,
	})
}

func (r *Request) SendSuccess(data map[string]string) string {
	data["success"] = "true"
	return r.send(data)
}

func (r *Request) SendFailure(data map[string]string) string {
	data["success"] = "false"
	return r.send(data)
}

func Send(data map[string]string) string {
	return SendSuccess(data)
}

func SendOutput(i int, data string) string {
	return Send(map[string]string{
		OutputDataPrefix + strconv.Itoa(i): data,
	})
}

func SendSuccess(data map[string]string) string {
	s := getStream()
	data["success"] = "true"
	data["request_id"] = util.GenerateUUID()
	data["node_id"] = config.GetEnv().SpNodeId
	return s.send(data)
}

func SendFailure(data map[string]string) string {
	s := getStream()
	data["success"] = "false"
	data["request_id"] = util.GenerateUUID()
	data["node_id"] = config.GetEnv().SpNodeId
	return s.send(data)
}

func (r *Request) send(data map[string]string) string {
	s := getStream()
	data["request_id"] = r.ID
	data["extra"] = r.Extra
	data["node_id"] = config.GetEnv().SpNodeId
	return s.send(data)
}

func (s *Stream) send(data map[string]string) string {
	return s.MessageQueue.SendMessage(s.StreamSendQueue, data, s.StreamSendQueueMaxLength, s.StreamSendQueueTrimImmediately)
}

func (s *Stream) subscribe() <-chan Request {
	q := mq.New(config.GetArgs())
	group := config.GetEnv().SpNodeGroup
	consumer := config.GetEnv().SpNodeId

	reqs := make(chan Request)

	go func() {
		for msg := range q.SubscribeQueue(s.StreamRecvQueue, group, consumer) {
			var req Request
			mapstructure.Decode(msg, &req)
			req.Input = make(map[string]string)
			for k, v := range msg {
				match, err := regexp.MatchString(InputPattern, k)
				if err != nil {
					log.Errorf("Message regex match error: %v", err)
				}
				if match {
					req.Input[k] = v.(string)
				}
			}
			reqs <- req
		}
	}()

	return reqs
}

func (r *Request) InputData(i int) string {
	return r.Input[InputDataPrefix+strconv.Itoa(i)]
}
