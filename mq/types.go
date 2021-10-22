package mq

type QueueMessage struct {
	ID    string
	Data  map[string]interface{}
	Queue string
	Group string
}
