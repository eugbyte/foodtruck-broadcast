package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func echo(w http.ResponseWriter, r *http.Request) {
	wsServer := New()
	if err := wsServer.Open(w, r); err != nil {
		panic(err)
	}

	for {
		msg, err := wsServer.Read()
		if err != nil {
			panic(err)
		}
		newMsg := fmt.Sprintf("received: %s", string(msg))
		if err := wsServer.Write([]byte(newMsg)); err != nil {
			panic(err)
		}
	}
}

func TestWSServer(t *testing.T) {
	assert := assert.New(t)

	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	client, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer client.Close()

	// Send message to server, read response and check to see if it's what we expect.
	if err := client.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
		t.Fatalf("%v", err)
	}

	_, msg, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}
	assert.Equal("received: hello", string(msg), "messages should be equal")

}
