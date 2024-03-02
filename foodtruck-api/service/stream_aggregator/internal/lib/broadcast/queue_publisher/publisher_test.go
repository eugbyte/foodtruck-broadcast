package kafkapublisher

import (
	wslib "foodtruck/pkg/websocket/server"
	wshub "foodtruck/service/stream_aggregator/internal/lib/broadcast/ws_hub"
	"foodtruck/service/stream_aggregator/internal/lib/mock"
	"testing"

	assrt "github.com/blend/go-sdk/assert"
)

func TestBroadcast(t *testing.T) {
	var assert = assrt.New(t)

	hub := wshub.New()
	conn1 := wslib.New()
	conn2 := wslib.New()
	reader := mock.NewMockReader()

	hub.Add("1", conn1)
	hub.Add("2", conn2)

	go PublishMsg(reader, hub.BroadcastCh())
	go hub.Broadcast()

	msg1 := <-conn1.Channel()
	msg2 := <-conn2.Channel()

	assert.NotNil(msg1)
	assert.NotNil(msg2)

	assert.Equal(string(msg1), string(msg2), "same message should be broadcasted to all subscribed ws conns")
}
