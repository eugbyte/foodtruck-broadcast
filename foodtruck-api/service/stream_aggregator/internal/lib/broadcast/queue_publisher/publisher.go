package kafkapublisher

import (
	kafkalib "foodtruck/pkg/queue/kafkalib/reader"

	debug "foodtruck/pkg/logger"
)

type Message = kafkalib.Message

type qReader interface {
	Read() (Message, error)
	Close() error
}

var logger = debug.Logger

// Blocking operation that opens the kafka connections, dequeues the message from kafka, and as a side effect,
// publish copies of the message to all subscribers
func PublishMsg(reader qReader, subscribers ...chan []byte) {
	defer reader.Close()

	for {
		msg, err := reader.Read()
		if err != nil {
			logger.Info(err)
			continue
		}
		for _, subscriber := range subscribers {
			subscriber <- msg.Value
		}
	}
}
