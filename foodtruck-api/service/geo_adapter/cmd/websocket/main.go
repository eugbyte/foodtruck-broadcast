package main

import (
	"encoding/json"
	"flag"
	"fmt"
	debug "foodtruck/pkg/logger"
	kafreader "foodtruck/pkg/queue/kafkalib/reader"
	kafwriter "foodtruck/pkg/queue/kafkalib/writer"
	ws_handler "foodtruck/service/geoadapter/internal/handler/websocket"
	"foodtruck/service/geoadapter/internal/lib/config"
	"log"
	"net/http"
)

var logger = debug.Logger

var topic = config.KAFKA_TOPIC
var kafkaPort = config.KAFKA_PORT
var addr = flag.String("addr", ":6000", "http service address")

func main() {
	logger.Info("started geo_adapter")
	logger.Info(fmt.Sprintf("topic: %s. kafkaPort: %s", topic, kafkaPort))

	var reader = kafreader.New(topic, "cg-1", kafkaPort)
	defer reader.Close()

	var writer = kafwriter.New(topic, kafkaPort)
	defer writer.Close()

	if err := writer.CreateTopic(); err != nil {
		logger.Info(err.Error())
	}

	http.HandleFunc("/vendor", func(w http.ResponseWriter, r *http.Request) {
		wshandler := ws_handler.New(writer) // stateful to each ws request
		err := wshandler.Handle(w, r)
		panic(err)
	})

	http.HandleFunc("/test-reader", func(w http.ResponseWriter, r *http.Request) {
		msg, err := reader.Read()
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(msg.Value); err != nil {
			logger.Error(err)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"msg": "Healthy"}); err != nil {
			logger.Info(err)
			w.Header().Set("Content-Type", "application/json")
		}
	})

	logger.Info("listening to ", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
