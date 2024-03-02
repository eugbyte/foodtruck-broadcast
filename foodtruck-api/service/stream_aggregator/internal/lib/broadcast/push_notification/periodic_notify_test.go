package notify

import (
	pub "foodtruck/service/stream_aggregator/internal/lib/broadcast/queue_publisher"
	"foodtruck/service/stream_aggregator/internal/lib/mock"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	assrt "github.com/stretchr/testify/assert"
)

type MockRestClient struct {
	Ch chan map[string]any
}

func NewMockRestClient() *MockRestClient {
	return &MockRestClient{
		Ch: make(chan map[string]any, 100),
	}
}

func (m *MockRestClient) Post(payload map[string]any) (respBody []byte, err error) {
	m.Ch <- payload
	return []byte{}, nil
}

func TestNotify(t *testing.T) {
	var assert = assrt.New(t)

	reader := mock.NewMockReader()
	var notifier = NewPeriodicNotifier(NewMockRestClient())

	go pub.PublishMsg(reader, notifier.Channel())

	msg1 := <-notifier.Channel()
	assert.NotNil(msg1)

}

func TestCreateNotifications(t *testing.T) {
	var assert = assrt.New(t)
	restClient := NewMockRestClient()
	var notifier = NewPeriodicNotifier(restClient)

	geohashes := make(map[string]mapset.Set[string])
	geohashes["123"] = mapset.NewSet[string]("abc", "def")
	geohashes["456"] = mapset.NewSet[string]("abc", "def")

	err := notifier.createNotifications(geohashes)
	assert.NoError(err)
	assert.Equal(len(restClient.Ch), 2)

	payload := <-restClient.Ch
	logger.Info(payload)

	assert.True(payload["geohash"] == "123" || payload["geohash"] == "456")
}
