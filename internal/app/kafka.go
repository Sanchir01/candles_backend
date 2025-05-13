package app

import (
	"context"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	producer *kafka.Writer
}

var errUnknownType = fmt.Errorf("unknown type")

var flushTimeout = 5000

func ensureTopicExists(brokers []string, topic string, partitions int, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial kafka broker: %w", err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		},
	}

	return conn.CreateTopics(topicConfigs...)
}
func NewProducer(address []string, topic string) (*Producer, error) {
	if len(address) == 0 {
		return nil, fmt.Errorf("no Kafka brokers provided")
	}
	err := ensureTopicExists(address, topic, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}
	writer := &kafka.Writer{
		Addr:     kafka.TCP(address...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		producer: writer,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, message string) error {

	if err := p.producer.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	}); err != nil {
		return err
	}
	if err := p.producer.Close(); err != nil {
		return err
	}
	return nil
}
