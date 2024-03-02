// Used to publish and subscribe securely to/from a kafka queue.
// Primarily for sensor integrations while utilising the sidecar pattern.
package writer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Interface to help with mocking Kafka for unit tests
type KafkaPublisher interface {
	PublishMessage(string, interface{}) error
}

// Contains the type of message (E.g. POS) and the contents of a Kafka message
type KafkaMessage struct {
	Type string
	Body []byte
}

// Create an instance of a kafka writer
type Writer struct {
	writer *kafka.Writer
}

// Instatiates a new kafka writer/publisher with a
// logger, writer, connection properties, and whether or not to use SSL.
// N.B. Always use SSL in a integration/production environment
func New(topic string, conn ...string) *Writer {
	return &Writer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(conn...),
			Topic:                  topic,
			AllowAutoTopicCreation: true,
		},
	}
}

func (k *Writer) Close() error {
	return k.writer.Close()
}

func (k *Writer) CreateTopic() error {
	// to create topics when auto.create.topics.enable='true'
	conn, err := kafka.Dial("tcp", k.writer.Addr.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             k.writer.Topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}

// Method used to publish a message to an instance of a kafka queue.
func (k *Writer) Write(message any) error {
	if message == nil {
		return errors.New("message is nil")
	}

	byts, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling:%v", err)
	}

	msg := kafka.Message{
		Key:   []byte(uuid.New().String()),
		Time:  time.Now(),
		Value: byts,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return k.writer.WriteMessages(ctx, msg)
}
