package config

import (
	"foodtruck/pkg/config"
)

var STAGE string
var KAFKA_PORT string
var KAFKA_TOPIC string

func init() {
	STAGE = config.GetEnv("STAGE", "dev")
	KAFKA_TOPIC = "foodtruck"

	var kafkaPorts = map[string]string{
		"dev": "localhost:9092",
		"int": "kafka:29092",
	}
	KAFKA_PORT = kafkaPorts[STAGE]
}
