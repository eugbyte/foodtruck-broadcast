package consumer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/samber/lo"
	assrt "github.com/stretchr/testify/assert"
)

type MockWebPusher struct {
	Msgs []NotificationMsg
}

func NewMockWebPusher() *MockWebPusher {
	return &MockWebPusher{
		Msgs: make([]NotificationMsg, 0),
	}
}

func (w *MockWebPusher) Send(msg NotificationMsg, endpoint string, auth string, p256dh string, ttl int) error {
	w.Msgs = append(w.Msgs, msg)
	return nil
}

type MockSubRepo struct {
	Geohash string
	Subs    []Subscription
}

func NewMockSubRepo() *MockSubRepo {
	return &MockSubRepo{
		Subs: make([]Subscription, 0),
	}
}

var sub1 = Subscription{
	Endpoint: "abc",
	Geohash:  "123",
	LastSend: time.Now().Unix(),
}
var sub2 = Subscription{
	Endpoint: "def",
	Geohash:  "456",
	LastSend: time.Now().Unix(),
}

func (m *MockSubRepo) GetAllBefore(geohash string, before time.Time) ([]Subscription, error) {
	m.Geohash = geohash
	return []Subscription{sub1, sub2}, nil
}

func (m *MockSubRepo) GetAll(geohash string) ([]Subscription, error) {
	m.Geohash = geohash
	return []Subscription{sub1, sub2}, nil
}

func (m *MockSubRepo) BatchPut(subs []Subscription) error {
	m.Subs = subs
	return nil
}

func TestHandler(t *testing.T) {
	var assert = assrt.New(t)

	mockRepo := NewMockSubRepo()
	mockWebpusher := NewMockWebPusher()

	handler := New(mockWebpusher, mockRepo)

	vendorIDs := []string{"vendor1", "vendor2"}
	byts, err := json.Marshal(vendorIDs)
	assert.Nil(err)

	mockMsg := events.SQSMessage{
		MessageAttributes: map[string]events.SQSMessageAttribute{
			"geohash": {
				StringValue: lo.ToPtr("mock_geohash"),
			},
		},
		Body: string(byts),
	}

	event := MsgEvent{
		Records: []events.SQSMessage{mockMsg},
	}
	err = handler.Handle(event)
	assert.Nil(err)

	assert.Equal(mockRepo.Geohash, "mock_geohash")
}

func TestSendMsg(t *testing.T) {
	var assert = assrt.New(t)

	mockRepo := NewMockSubRepo()
	mockWebpusher := NewMockWebPusher()

	handler := New(mockWebpusher, mockRepo)

	err := handler.sendMsg("hello", []Subscription{sub1, sub2})
	assert.Nil(err)

	assert.Equal(len(mockWebpusher.Msgs), 2)
}
