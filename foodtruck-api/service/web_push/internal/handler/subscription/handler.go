package subscription

import (
	"encoding/json"
	subrepo "foodtruck/service/web_push/internal/lib/subscription"
	"net/http"

	debug "foodtruck/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
)

var logger = debug.Logger

type Subscription = subrepo.Subscription
type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

var decodeErrResp = Response{
	Headers:    map[string]string{"Content-Type": "application/json"},
	StatusCode: http.StatusInternalServerError,
	Body:       `{"message":"Failed to unmarshal struct"}`,
}

type SubRepo interface {
	Open()
	Put(sub Subscription) error
}

type handler struct {
	subRepo SubRepo
}

func New(subRepo SubRepo) *handler {
	return &handler{
		subRepo: subRepo,
	}
}

func (h *handler) Handle(req Request) (resp Response, err error) {
	var subscription Subscription
	if err := json.Unmarshal([]byte(req.Body), &subscription); err != nil {
		errResp := decodeErrResp
		errResp.StatusCode = http.StatusBadRequest
		return errResp, nil
	}

	if err = h.subscribe(subscription); err != nil {
		logger.Error(err)
		return Response{
			Headers:    map[string]string{"Content-Type": "application/json"},
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"Failed to update DB"}`,
		}, nil
	}
	logger.Info("update success")

	byts, err := json.Marshal(map[string]string{
		"message": "subscription saved",
	})
	if err != nil {
		return decodeErrResp, nil
	}

	return Response{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
		},
		StatusCode: http.StatusOK,
		Body:       string(byts),
	}, nil
}
