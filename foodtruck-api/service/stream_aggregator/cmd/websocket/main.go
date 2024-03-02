package main

import (
	"encoding/json"
	"flag"
	debug "foodtruck/pkg/logger"
	kafreader "foodtruck/pkg/queue/kafkalib/reader"
	ws_handler "foodtruck/service/stream_aggregator/internal/handler/websocket"
	notify "foodtruck/service/stream_aggregator/internal/lib/broadcast/push_notification"
	pub "foodtruck/service/stream_aggregator/internal/lib/broadcast/queue_publisher"
	wshub "foodtruck/service/stream_aggregator/internal/lib/broadcast/ws_hub"
	"foodtruck/service/stream_aggregator/internal/lib/config"
	"foodtruck/service/stream_aggregator/internal/lib/mock"
	"log"
	"net/http"
	"time"
)

var logger = debug.Logger

var addr = flag.String("addr", ":8080", "http service address")
var topic = config.KAFKA_TOPIC
var kafkaPort = config.KAFKA_PORT

func main() {
	logger.Info("stream_aggregator started")
	// kafka reader
	var reader = kafreader.New(topic, "cg-1", kafkaPort)

	// singleton hub to register ws connections
	var hub = wshub.New()
	// notification rest api lib
	var notifyAPILib = mock.NewMockNotifyAPI()
	// notification service to create notification
	var notifier = notify.NewPeriodicNotifier(notifyAPILib)

	// broadcast dequeued kafka msg to all registered ws conns in the hub
	go pub.PublishMsg(reader, hub.BroadcastCh(), notifier.Channel())
	go hub.Broadcast()
	go notifier.PeriodicNotify(time.Minute * 1)

	http.HandleFunc("/customer", func(w http.ResponseWriter, r *http.Request) {
		wshandler := ws_handler.New(hub) // stateful to each ws request
		wshandler.Handle(w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"msg": "Healthy"}); err != nil {
			logger.Info(err)
		}
	})

	logger.Info("listening to ", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
