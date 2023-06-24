package kafka

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
)

type Consumer struct {
	ConfigMap *ckafka.ConfigMap
	Topics    []string
}

func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer {
	return &Consumer{
		ConfigMap: configMap,
		Topics:    topics,
	}
}

func (c *Consumer) Consume(msgChan chan *ckafka.Message) error {
	consume, err := ckafka.NewConsumer(c.ConfigMap)

	if err != nil {
		return errors.New("error in the new consumer in the consumer")
	}

	err = consume.SubscribeTopics(c.Topics, nil)

	if err != nil {
		return errors.New("error in the subscribe topics in the consumer")
	}

	for {
		msg, err := consume.ReadMessage(-1)

		if err == nil {
			msgChan <- msg
		}
	}
}
