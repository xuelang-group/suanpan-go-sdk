package mq

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/thoas/go-funk"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/log"
)

type RedisMq struct {
	MqRedisHost string `mapstructure:"--mq-redis-host" default:"localhost"`
	MqRedisPort string `mapstructure:"--mq-redis-port" default:"6379"`
}

func (r *RedisMq) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: r.MqRedisHost + ":" + r.MqRedisPort,
	})
}

func (r *RedisMq) recvMessages(queue, group, consumer, consumeID string) []QueueMessage {
	cli := r.getClient()
	args := &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{queue, consumeID},
	}
	res, err := cli.XReadGroup(context.Background(), args).Result()
	if err != nil {
		log.Errorf("Read redis group failed: %w", err)
	}

	messages := make([]QueueMessage, 0)
	for _, x := range res {
		for _, m := range x.Messages {
			messages = append(messages, QueueMessage{
				ID:    m.ID,
				Data:  m.Values,
				Queue: x.Stream,
				Group: group,
			})
		}
	}

	lostMessages := funk.Filter(messages, func(m QueueMessage) bool {
		return m.ID != "" && m.Data == nil
	}).([]QueueMessage)
	lostMessageIDs := funk.Map(lostMessages, func(m QueueMessage) string {
		return m.ID
	}).([]string)
	if len(lostMessages) > 0 {
		cli.XAck(context.Background(), queue, group, lostMessageIDs...)
		log.Warnf("Messages have lost: %w", lostMessageIDs)
	}

	return messages
}

func (r *RedisMq) createQueue(queue, group, consumeID string) {
	cli := r.getClient()
	log.Infof("Create queue %s-%s", queue, group)
	err := cli.XGroupCreateMkStream(context.Background(), queue, group, consumeID).Err()
	if err != nil {
		log.Warnf("Create redis queue error: %w", err)
	}
}

func (r *RedisMq) SubscribeQueue(queue, group, consumer string) <-chan map[string]interface{} {
	cli := r.getClient()
	r.createQueue(queue, group, "0")
	log.Info("Subscribing message")

	msg := make(chan map[string]interface{})

	go func() {
		for {
			messages := r.recvMessages(queue, group, consumer, ">")
			for _, message := range messages {
				msg <- message.Data
				cli.XAck(context.Background(), queue, group, message.ID)
			}
		}
	}()

	return msg
}

func (r *RedisMq) SendMessage(queue string, data map[string]string, maxLen int64, trimImmediately bool) string {
	cli := r.getClient()
	args := &redis.XAddArgs{
		Stream: queue,
		Values: data,
		ID:     "*",
		MaxLen: maxLen,
		Approx: !trimImmediately,
	}
	id, err := cli.XAdd(context.Background(), args).Result()
	if err != nil {
		log.Errorf("Send redis message failed: %w", err)
	}

	return id
}
