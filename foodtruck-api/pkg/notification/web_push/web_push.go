package webpush

import (
	"encoding/json"
	"errors"
	debug "foodtruck/pkg/logger"
	"io"

	webpush "github.com/SherClockHolmes/webpush-go"
)

var logger = debug.Logger

type Message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type webPusher struct {
	vapidPrivateKey string
	VapidPublicKey  string
	SenderEmail     string
}

func New(VAPIDPrivateKey string, VAPIDPublicKey string, senderEmail string) *webPusher {
	wp := webPusher{
		vapidPrivateKey: VAPIDPrivateKey,
		VapidPublicKey:  VAPIDPublicKey,
		SenderEmail:     senderEmail,
	}
	return &wp
}

// Send the notification to the browser.
func (wp *webPusher) Send(msg Message, endpoint string, auth string, p256dh string, ttl int) error {
	byts, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	sub := &webpush.Subscription{
		Endpoint: endpoint,
		Keys: webpush.Keys{
			Auth:   auth,
			P256dh: p256dh,
		},
	}

	// Send Notification
	resp, err := webpush.SendNotification(byts, sub, &webpush.Options{
		Subscriber:      wp.SenderEmail,
		VAPIDPublicKey:  wp.VapidPublicKey,
		VAPIDPrivateKey: wp.vapidPrivateKey,
		TTL:             ttl, // seconds
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// https://web.dev/push-notifications-web-push-protocol/#response-from-push-service
	logger.Info(resp.Status, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return errors.New(string(body))
	}

	return nil
}
