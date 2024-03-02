package main

import (
	"foodtruck/service/web_push/internal/handler/subscription"
	"foodtruck/service/web_push/internal/lib/config"
	subrepo "foodtruck/service/web_push/internal/lib/subscription"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	subRepo := subrepo.New("subscription", "ap-southeast-1", config.DBConn())
	subRepo.Open()
	subHandler := subscription.New(subRepo)

	lambda.Start(subHandler.Handle)
}
