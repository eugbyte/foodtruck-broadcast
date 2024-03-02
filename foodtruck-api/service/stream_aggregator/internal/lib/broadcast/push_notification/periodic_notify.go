package notify

import (
	"encoding/json"
	"errors"
	"foodtruck/pkg/model"
	"strings"
	"sync"
	"time"

	debug "foodtruck/pkg/logger"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mmcloughlin/geohash"
)

var logger = debug.Logger

type GeoHash = string
type VendorID = string

type NotificationAPI interface {
	Post(payload map[string]any) (respBody []byte, err error)
}

type PeriodicNotifier struct {
	ch        chan []byte
	notifyAPI NotificationAPI
}

func NewPeriodicNotifier(notifyAPILib NotificationAPI) *PeriodicNotifier {
	return &PeriodicNotifier{
		ch:        make(chan []byte),
		notifyAPI: notifyAPILib,
	}
}

func (n *PeriodicNotifier) Channel() chan []byte { return n.ch }

// Blocking operation that calls the notification API periodically.
func (n *PeriodicNotifier) PeriodicNotify(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	geohashes := make(map[GeoHash]mapset.Set[VendorID])

	for {
		select {
		case <-ticker.C:
			if err := n.createNotifications(geohashes); err != nil {
				logger.Error(err)
			}
			// reset the geohash map
			geohashes = make(map[string]mapset.Set[string])
		case byts := <-n.ch:
			var geoInfo model.GeoInfo
			if err := json.Unmarshal(byts, &geoInfo); err != nil {
				logger.Error(err)
				continue
			}
			addToMap(geohashes, geoInfo)
		}
	}
}

// Create notifications per geohash concurrently via the NotificationAPI.
func (n *PeriodicNotifier) createNotifications(geohashes map[GeoHash]mapset.Set[string]) error {
	var errs = make([]string, 0)
	var wg sync.WaitGroup

	for hash, set := range geohashes {
		wg.Add(1)
		payload := map[string]any{
			"geohash":   hash,
			"vendorIDs": set.ToSlice(),
		}

		// Make a post request to the web notification producer repo
		go func() {
			defer wg.Done()
			if _, err := n.notifyAPI.Post(payload); err != nil {
				errs = append(errs, err.Error())
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ".\n"))
	}
	return nil
}

// As a side effect, unmarshals the msg as a geohash-coordinate key-value pair, and add it to the map.
func addToMap(geohashes map[GeoHash]mapset.Set[string], geoInfo model.GeoInfo) {
	hash := geohash.Encode(geoInfo.Lat, geoInfo.Lng)
	if _, ok := geohashes[hash]; !ok {
		geohashes[hash] = mapset.NewSet[VendorID]()
	}
	geohashes[hash].Add(geoInfo.VendorID)
}
