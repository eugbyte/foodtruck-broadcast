package ws

import (
	"time"

	"github.com/pkg/errors"
)

func (h *handler) heartbeat(wsID string) {
	// Whenever the client pongs back to the server in response to a ping from a server, extend the read deadline
	h.onClientPong(wsID)

	for range h.pingTick.C {
		h.pingClient(wsID)
	}
}

// Ping the client to check whether the connection is alive. For heatbeat purpose, non-blocking.
func (h *handler) pingClient(wsID string) {
	conn := h.hub.Conns()[wsID]

	if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		h.fatal <- errors.Wrap(err, "write deadline error")
		return
	}
	logger.Infof("%s: pinging...", wsID)
	if err := conn.Ping(); err != nil {
		h.fatal <- errors.Wrap(err, "ping error error")
		return
	}
}

// Whenever the ws client sends a message or pongs the server, extend the read deadline. For heatbeat purpose, non-blocking.
func (h *handler) onClientPong(wsID string) {
	conn := h.hub.Conns()[wsID]

	// Whenever the client pongs back to the server in response to a ping from a server, extend the read deadline
	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		h.fatal <- errors.Wrap(err, "read deadline error")
		return
	}

	conn.OnClientPong(func(appData string) error {
		logger.Infof("%s: pong received", wsID)
		if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			h.fatal <- errors.Wrap(err, "onPong error")
		}
		return nil
	})
}
