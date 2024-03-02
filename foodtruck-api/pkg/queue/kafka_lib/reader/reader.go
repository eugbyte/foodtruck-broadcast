package reader

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	reader *kafka.Reader
}

func New(
	topic string,
	consumerGroup string,
	conn ...string,
) *Reader {
	return &Reader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  conn,
			Topic:    topic,
			GroupID:  consumerGroup,
			MaxBytes: 10e6, // 10MB
		}),
	}
}

type Message = kafka.Message

// Blocks until a message becomes available, or an error occurs
func (k *Reader) Read() (Message, error) {
	if k.reader == nil {
		panic("reader is not initialized. Remember to call Open().")
	}
	return k.reader.ReadMessage(context.TODO())
}

func (k *Reader) Close() error {
	if k.reader == nil {
		return errors.New("kafka reader is not initialized")
	}
	return k.reader.Close()
}

func (k *Reader) SetKafkaOffset(startTime time.Time) error {
	return (*k.reader).SetOffsetAt(context.TODO(), startTime) //won't work if consumer group is set.
}
