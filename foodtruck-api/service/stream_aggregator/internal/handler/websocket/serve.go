package ws

import (
	wsserver "foodtruck/pkg/websocket/server"
	"net/http"

	wshub "foodtruck/service/stream_aggregator/internal/lib/broadcast/ws_hub"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mmcloughlin/geohash"
)

// Establish the websocket connection and run the process. Blocking operation.
func (h *handler) serve(w http.ResponseWriter, r *http.Request) {
	var conn wshub.WSServer = wsserver.New()
	if err := conn.Open(w, r); err != nil {
		logger.Info("failed to upgrade connection")
		return
	}

	wsID, err := gonanoid.New()
	if err != nil {
		logger.Error(err)
		return
	}

	geoHash := r.URL.Query().Get("geohash")
	logger.Infof("%s: geoHash: %s", wsID, geoHash)
	if err := geohash.Validate(geoHash); err == nil {
		box := geohash.BoundingBox(geoHash)
		h.setBoundBox(&box)
	}
	logger.Info(h.boundingBox)

	h.hub.Add(wsID, conn)
	logger.Info("registered ws conn: ", wsID)

	go h.readPump(wsID)
	go h.writePump(wsID)
	go h.heartbeat(wsID)

	for {
		select {
		case sig := <-h.sig:
			logger.Info("connection closing normally: ", sig)
			h.pingTick.Stop()  // for the ticker to be recovered by the garbage collector
			h.hub.Remove(wsID) // unregister the connection
			if err := conn.Close(); err != nil {
				logger.Error(err)
			}
			return
		case err := <-h.fatal:
			logger.Infof("%s - fatal error: %v", wsID, err)
			h.pingTick.Stop()  // for the ticker to be recovered by the garbage collector
			h.hub.Remove(wsID) // unregister the connection

			// if error is due to a closed connection, simply return. Otherwise, close the connection.
			if conn.IsClosedError(err) {
				logger.Info("connection already closed: ", err)
				return
			}
			if err := conn.Close(); err != nil {
				logger.Error(err)
			}
			return
		}
	}

}
