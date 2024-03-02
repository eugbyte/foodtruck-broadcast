package main

import (
	"foodtruck/service/web_push/internal/handler/consumer"
	subrepo "foodtruck/service/web_push/internal/lib/subscription"

	webpush "foodtruck/pkg/notification/web_push"
	"foodtruck/service/web_push/internal/lib/config"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	webpusher := webpush.New(config.VAPID_PRIVATE_KEY, config.VAPID_PUBLIC_KEY, config.VAPID_EMAIL)
	subRepo := subrepo.New("subscription", "ap-southeast-1", config.DBConn())
	subRepo.Open()
	consumerHandler := consumer.New(webpusher, subRepo)
	lambda.Start(consumerHandler.Handle)
}
