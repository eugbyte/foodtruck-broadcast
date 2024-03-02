package main

import (
	sqslib "foodtruck/pkg/queue/sqs"
	"foodtruck/service/web_push/internal/handler/producer"

	"foodtruck/service/web_push/internal/lib/config"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	msgQ := sqslib.New("foodtruck-geo-notify", "ap-southeast-1", config.QConn())
	if err := msgQ.Open(); err != nil {
		panic(err)
	}
	producerHandler := producer.New(msgQ)

	lambda.Start(producerHandler.Handle)
}
