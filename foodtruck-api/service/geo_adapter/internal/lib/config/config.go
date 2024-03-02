package config

import (
	"foodtruck/pkg/config"
	debug "foodtruck/pkg/logger"
)

var STAGE string
var KAFKA_PORT string
var KAFKA_TOPIC string

var logger = debug.Logger

func init() {
	STAGE = config.GetEnv("STAGE", "dev")
	logger.Info("STAGE: ", STAGE)
	KAFKA_TOPIC = "foodtruck"

	var kafkaPorts = map[string]string{
		"dev": "localhost:9092",
		"int": "kafka:29092",
	}
	KAFKA_PORT = kafkaPorts[STAGE]
}
