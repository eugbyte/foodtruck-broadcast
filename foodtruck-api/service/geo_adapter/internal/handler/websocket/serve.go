package websocket

import (
	"foodtruck/pkg/model"
	"net/http"
)

func (h *Handler) serve(w http.ResponseWriter, r *http.Request) error {
	if err := h.conn.Open(w, r); err != nil {
		logger.Info("failed to upgrade connection")
		return err
	}
	defer h.conn.Close()

	logger.Info("ws conn established")

	go h.readPump()

	for {
		select {
		case <-h.sig:
			logger.Info("interupt sig")
			return nil
		case err := <-h.fatal:
			return err
		case geoInfo := <-h.wsRead:
			logger.Info(geoInfo)
			if err := h.enqueue(geoInfo); err != nil {
				return err
			}
			logger.Info("successfully enqueued")
		}
	}
}

func (h *Handler) enqueue(geoInfo model.GeoInfo) error {
	return h.queue.Write(geoInfo)
}
