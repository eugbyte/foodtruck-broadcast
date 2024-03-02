package websocket

import (
	debug "foodtruck/pkg/logger"
	"foodtruck/pkg/model"
	ws "foodtruck/pkg/websocket/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var logger = debug.Logger

type WSServer interface {
	Read() (msg []byte, err error)
	Open(w http.ResponseWriter, r *http.Request) error
	Close() error
}

type Queuer interface {
	Write(msg any) error
}

type Handler struct {
	queue  Queuer
	conn   WSServer
	sig    chan os.Signal
	fatal  chan error
	wsRead chan model.GeoInfo
}

func New(queue Queuer) *Handler {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return &Handler{
		queue:  queue,
		sig:    sigs,
		fatal:  make(chan error),
		wsRead: make(chan model.GeoInfo),
		conn:   ws.New(),
	}
}

// Entry point to handler.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) error {
	return h.serve(w, r)
}
