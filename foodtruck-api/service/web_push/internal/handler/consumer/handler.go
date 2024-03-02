package consumer

import (
	"encoding/json"
	errs "errors"
	debug "foodtruck/pkg/logger"
	webpush "foodtruck/pkg/notification/web_push"
	sub "foodtruck/service/web_push/internal/handler/subscription"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

var logger = debug.Logger

type NotificationMsg = webpush.Message
type MsgEvent = events.SQSEvent
type Subscription = sub.Subscription

type Webpusher interface {
	// Send the notification to the browser.
	Send(msg NotificationMsg, endpoint string, auth string, p256dh string, ttl int) error
}

type SubRepo interface {
	GetAll(geohash string) ([]Subscription, error)
	// Get all users within a location before a certain time.
	GetAllBefore(geohash string, before time.Time) ([]Subscription, error)
	// Batch update.
	BatchPut(subs []Subscription) error
}

type handler struct {
	webpusher Webpusher
	subrepo   SubRepo
}

func New(webpusher Webpusher, subrepo SubRepo) *handler {
	return &handler{
		webpusher: webpusher,
		subrepo:   subrepo,
	}
}

func (h *handler) Handle(event MsgEvent) error {
	logger.Info("SQS event detected")
	if len(event.Records) == 0 {
		return errs.New("no msg in queue")
	}

	// aggregated error
	var aggrErr = make([]string, 0)

	for _, msg := range event.Records {
		geohash, ok := msg.MessageAttributes["geohash"]
		if !ok {
			logger.Error(errs.New("geohash is empty"))
			continue
		}

		var vendorIDs []string
		if err := json.Unmarshal([]byte(msg.Body), &vendorIDs); err != nil {
			return err
		}

		logger.Info(vendorIDs)
		if err := h.notify(*geohash.StringValue, vendorIDs); err != nil {
			logger.Info("consume error: ", err.Error())
			aggrErr = append(aggrErr, err.Error())
		}
	}

	if len(aggrErr) > 0 {
		logger.Info("aggregate error detected")
		logger.Info(strings.Join(aggrErr, ". "))
		return errs.New(strings.Join(aggrErr, ". "))
	}
	logger.Info("consumer success")
	return nil
}
