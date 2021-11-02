package mq

type QueueMessage struct {
	ID    string
	Data  map[string]string
	Queue string
	Group string
}
