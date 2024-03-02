package ws

import (
	"encoding/json"
	"foodtruck/pkg/model"
	"time"

	geohashlib "github.com/mmcloughlin/geohash"
	"github.com/pkg/errors"
)

// Send messages received from the hub to the ws clients
func (h *handler) writePump(wsID string) {
	conn := h.hub.Conns()[wsID]

	for msg := range conn.Channel() {
		var geoInfo model.GeoInfo
		if err := json.Unmarshal(msg, &geoInfo); err != nil {
			logger.Error(err)
			continue
		}

		// if user is not close to vendor, ignore
		// if !h.boundBox().Contains(geoInfo.Latitude, geoInfo.Longitude) {
		// 	continue
		// }

		byts, err := json.Marshal(Payload{
			Action: "geo_info",
			Data:   geoInfo,
		})
		if err != nil {
			logger.Error(err)
			continue
		}

		if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
			h.fatal <- errors.Wrap(err, "write deadline error")
			return
		}
		if err := conn.Write(byts); err != nil && conn.IsClosedError(err) {
			h.fatal <- errors.Wrap(err, "write error")
			return
		}

		hash := geohashlib.Encode(geoInfo.Lat, geoInfo.Lng)
		_ = hash
	}
}
