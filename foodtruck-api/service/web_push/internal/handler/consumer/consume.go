package consumer

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"
)

// Notify all users within the specified location of the vendors nearby,
// and update the DB of the latest time notified.
func (h *handler) notify(geohash string, vendorIDs []string) error {
	timestamp := time.Now().Add(time.Minute * -5)
	logger.Info(geohash, vendorIDs)

	subscriptions, err := h.subrepo.GetAllBefore(geohash, timestamp)
	if err != nil {
		return fmt.Errorf("error querying DB: %v", err)
	}
	if len(subscriptions) == 0 {
		return errors.New("no subscriptions found")
	}

	// top 3 vendor IDs
	var topIDs []string = lo.Slice(vendorIDs, 0, 3)
	var msg string = fmt.Sprintf("foodtrucks %s ... are in your area.", strings.Join(topIDs, ", "))
	logger.Info("subscription msg: ", msg)

	if err := h.sendMsg(msg, subscriptions); err != nil {
		logger.Error(err)
	}

	// range create a copy of the struct, rather than value
	// https://yourbasic.org/golang/gotcha-range-copy-array/
	for idx := range subscriptions {
		subscriptions[idx].LastSend = time.Now().Unix()
	}
	return h.batchUpdate(subscriptions)
}

// Send the notification message to all subscribers.
func (h *handler) sendMsg(text string, subscribers []Subscription) error {
	var wg sync.WaitGroup

	var err error

	for _, sub := range subscribers {
		wg.Add(1)

		msg := NotificationMsg{
			Title: "foodtruck",
			Body:  text,
		}

		go func(s Subscription) {
			defer wg.Done()
			if _err := h.webpusher.Send(msg, s.Endpoint, s.Auth, s.P256dh, 60); err != nil {
				logger.Error(_err)
				err = _err
			}
		}(sub)
	}

	wg.Wait()
	logger.Info("finished sending web push message")
	return err
}

// batch update the subscriptions, updating the latest time the notification was sent.
func (h *handler) batchUpdate(subscriptions []Subscription) error {
	if err := h.subrepo.BatchPut(subscriptions); err != nil {
		return fmt.Errorf("batchUpdateError: %v", err)
	}
	return nil
}
