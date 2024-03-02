package ws

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	debug "foodtruck/pkg/logger"
	wshub "foodtruck/service/stream_aggregator/internal/lib/broadcast/ws_hub"

	"github.com/mmcloughlin/geohash"
)

type WSServer = wshub.WSServer

// Payload to and from client ws.
type Payload struct {
	Action string `json:"action"`
	Data   any    `json:"data"` // default type for objects is `map[string]any`
}

type Hub interface {
	Add(wsID string, wsServer WSServer)
	BroadcastCh() (broadcastChan chan []byte)
	Conns() map[string]WSServer
	Remove(wsID string)
}

type handler struct {
	hub         Hub // stateful singleton
	pingTick    *time.Ticker
	fatal       chan error
	sig         chan os.Signal
	mu          sync.Mutex
	boundingBox *geohash.Box
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var logger = debug.Logger

// To be used as a transient object, rather than a singleton.
// Unlike singleton objects, transient objects are different across different requests.
//
// Since WS are long lived connections, being used as a singleton will result in the different requests concurrently modifying the state of the struct.
func New(hub Hub) *handler {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return &handler{
		hub:      hub,
		pingTick: time.NewTicker(pingPeriod),
		fatal:    make(chan error),
		sig:      sigs,
		boundingBox: &geohash.Box{
			MinLat: -90,
			MaxLat: 90,
			MinLng: -180,
			MaxLng: 180,
		},
	}
}

// Entry point to handler.
func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.serve(w, r)
}

// Set the bounding box in a thread safe manner.
func (h *handler) setBoundBox(boundingBox *geohash.Box) {
	defer h.mu.Unlock()
	h.mu.Lock()
	h.boundingBox = boundingBox
}

// Get the bounding box in a thread safe manner.
func (h *handler) boundBox() (boundingBox *geohash.Box) {
	defer h.mu.Unlock()
	h.mu.Lock()
	return h.boundingBox
}
