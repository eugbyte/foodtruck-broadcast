package producer

import "encoding/json"

// Enqueue message to message queue.
//
// Loose coupling between the lambda and business logic in case, e.g. change cloud provider.
func (h *handler) Enqueue(geohash string, vendorIDs []string) (err error) {
	metadata := map[string]string{
		"geohash": geohash,
	}

	byts, err := json.Marshal(vendorIDs)
	if err != nil {
		return err
	}

	return h.msgQ.Enqueue(string(byts), metadata)
}
