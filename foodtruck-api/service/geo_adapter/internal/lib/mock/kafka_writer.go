package mock

import (
	"encoding/json"
	"fmt"
	debug "foodtruck/pkg/logger"
)

var logger = debug.Logger

const queueSize = 2

type MockWriter struct {
	msgCh chan string
}

func NewMockWriter() *MockWriter {
	return &MockWriter{
		msgCh: make(chan string, queueSize),
	}
}

func (m *MockWriter) Close() error { return nil }
func (m *MockWriter) Write(msg any) error {
	byts, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling:%v", err)
	}

	if len(m.msgCh) >= queueSize {
		logger.Info("clearing chan")
		<-m.msgCh
	}

	logger.Info(len(m.msgCh))

	m.msgCh <- string(byts)
	return nil
}

func (m *MockWriter) MsgCh() chan string {
	return m.msgCh
}
