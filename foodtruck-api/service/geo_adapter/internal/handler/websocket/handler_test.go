package websocket

import (
	"encoding/json"
	"foodtruck/pkg/model"
	"foodtruck/service/geoadapter/internal/lib/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	assrt "github.com/blend/go-sdk/assert"
	ws "github.com/gorilla/websocket"
)

func TestRead(t *testing.T) {
	assert := assrt.New(t)

	mockQueuer := mock.NewMockWriter()
	handler := New(mockQueuer)

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := handler.Handle(w, r)
			assert.Nil(err)
		}))
	defer server.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the server
	client, _, err := ws.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer client.Close()

	// Send message to server, read response and check to see if it's what we expect.
	var geoInfo = model.GeoInfo{
		VendorID: "abc",
		Lat:      1.29,
		Lng:      103.85,
	}

	byts, err := json.Marshal(geoInfo)
	assert.Nil(err)

	if err := client.WriteMessage(ws.TextMessage, byts); err != nil {
		assert.Nil(err)
	}

	msg := <-mockQueuer.MsgCh()
	assert.Equal(msg, string(byts), "message should be enqueued")
}
