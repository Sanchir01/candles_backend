package app

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

var errUnknownType = fmt.Errorf("unknown type")

var flushTimeout = 5000

func NewProducer(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("kafka.Producer: %w", err)
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Key:   nil,
	}
	kafkaChan := make(chan kafka.Event, 1)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("kafka.Producer.Produce: %w", err)
	}
	e := <-kafkaChan

	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
