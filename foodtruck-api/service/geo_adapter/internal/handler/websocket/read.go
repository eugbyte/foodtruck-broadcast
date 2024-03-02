package websocket

import (
	"encoding/json"
	"errors"
	"foodtruck/pkg/model"
	"strings"
)

// listens for messages from the ws conn, and transfers them to the read channel.
func (h *Handler) readPump() {
	for {
		byts, err := h.conn.Read()
		if err != nil {
			h.fatal <- err
			return
		}

		var geoInfo model.GeoInfo
		if err := json.Unmarshal(byts, &geoInfo); err != nil {
			logger.Info(err)
			continue
		}

		geoInfo.VendorID = strings.TrimSpace(geoInfo.VendorID)
		if geoInfo.VendorID == "" {
			h.fatal <- errors.New("bad request - vendorID is empty")
			return
		}

		h.wsRead <- geoInfo
	}
}
