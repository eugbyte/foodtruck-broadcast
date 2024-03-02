package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"syscall"
	"testing"
	"time"

	"foodtruck/pkg/model"
	wshub "foodtruck/service/stream_aggregator/internal/lib/broadcast/ws_hub"

	assrt "github.com/blend/go-sdk/assert"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/mmcloughlin/geohash"
	"github.com/samber/lo"
)

func NewMockWSClient(s *httptest.Server) (*websocket.Conn, error) {
	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	client, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	return client, err
}

func TestRead(t *testing.T) {
	assert := assrt.New(t)

	hub := wshub.New()
	handler := New(hub)

	s := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer s.Close()

	client, err := NewMockWSClient(s)
	assert.Nil(err)
	defer client.Close()

	box := geohash.Box{
		MinLat: 1,
		MaxLat: 2,
		MinLng: 3,
		MaxLng: 4,
	}

	payloadBody := make(map[string]any)
	if err := mapstructure.Decode(box, &payloadBody); err != nil {
		t.Fatalf("%v", err)
	}

	logger.Info(payloadBody)

	var payload Payload = Payload{
		Action: "bounding_box",
		Data:   payloadBody,
	}

	byts, err := json.Marshal(payload)
	assert.Nil(err)

	// Send message to server, read response and check to see if it's what we expect.
	if err := client.WriteMessage(websocket.TextMessage, byts); err != nil {
		t.Fatalf("%v", err)
	}

	for handler.boundingBox.MaxLat != 2 {
		continue
	}

	handler.sig <- syscall.SIGINT
	assert.Equal(*handler.boundingBox, box, "boundingBoxes should be equal")
}

func TestWrite(t *testing.T) {
	assert := assrt.New(t)

	hub := wshub.New()
	go hub.Broadcast()
	handler := New(hub)

	s := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer s.Close()

	client, err := NewMockWSClient(s)
	assert.Nil(err)
	defer client.Close()

	geoInfo := model.GeoInfo{
		VendorID: "abc",
		Lat:      1,
		Lng:      2,
	}
	byts, err := json.Marshal(geoInfo)
	assert.Nil(err)

	go func() {
		for {
			hub.BroadcastCh() <- byts
		}
	}()

	_, msg, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	var resp Payload
	err = json.Unmarshal(msg, &resp)
	assert.Nil(err)

	logger.Info("resp: ", resp)

	var actGeoInfo model.GeoInfo
	if err := mapstructure.Decode(resp.Data, &actGeoInfo); err != nil {
		t.Fatalf("%v", err)
	}

	handler.sig <- syscall.SIGINT
	assert.Equal(geoInfo, actGeoInfo, "geoInfo should be equal")
}

func TestHeartbeat(t *testing.T) {
	assert := assrt.New(t)

	hub := wshub.New()
	go hub.Broadcast()
	handler := New(hub)
	handler.pingTick = time.NewTicker(2 * time.Second)

	s := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer s.Close()

	client, err := NewMockWSClient(s)
	assert.Nil(err)
	defer client.Close()

	msg := lo.ToPtr("")
	client.SetPingHandler(func(appData string) error {
		_msg := "pinged " + appData
		logger.Info("pinged")
		logger.Info(appData)
		msg = &_msg
		return nil
	})

	go func() {
		_, _, _ = client.ReadMessage()
	}()

	for *msg == "" {
		continue
	}

	logger.Info(string(*msg))
	handler.sig <- syscall.SIGINT
	assert.NotNil(msg, "client should receive ping msg from ws server")
}
