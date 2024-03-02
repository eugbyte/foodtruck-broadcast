package ws

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/mmcloughlin/geohash"
	"github.com/pkg/errors"
)

// Receive messages from the ws client, and as a side effect, modifies the bounding box.
func (h *handler) readPump(wsID string) {
	conn := h.hub.Conns()[wsID]

	for {
		byts, err := conn.Read()
		if err != nil {
			h.fatal <- errors.Wrap(err, "read error")
			return
		}
		logger.Info(string(byts))

		var payload Payload
		if err := json.Unmarshal(byts, &payload); err != nil {
			logger.Error(err)
			continue
		}

		if payload.Action == "bounding_box" {
			var box geohash.Box
			if err := mapstructure.Decode(payload.Data, &box); err != nil {
				logger.Error(err)
				continue
			}

			h.setBoundBox(&box)
			logger.Infof("%s: message received: %v", wsID, *h.boundBox())
		}
	}

}
