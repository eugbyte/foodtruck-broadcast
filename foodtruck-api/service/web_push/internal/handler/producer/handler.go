package producer

import (
	"encoding/json"
	"net/http"

	debug "foodtruck/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
)

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

var logger = debug.Logger
var decodeErrResp = Response{
	Headers:    map[string]string{"Content-Type": "application/json"},
	StatusCode: http.StatusInternalServerError,
	Body:       `{"message":"Failed to unmarshal struct"}`,
}

type GeoInfo struct {
	VendorIDs []string `json:"vendorIDs"`
}

type MsgQer interface {
	Open() (err error)
	Enqueue(msg string, metadata map[string]string) error
}

type handler struct {
	msgQ MsgQer
}

func New(msgQ MsgQer) *handler {
	return &handler{
		msgQ: msgQ,
	}
}

func (h *handler) Handle(req Request) (resp Response, err error) {
	geohash := req.PathParameters["geohash"]
	logger.Info("geohash: ", geohash)

	if geohash == "" {
		return Response{
			Headers:    map[string]string{"Content-Type": "application/json"},
			StatusCode: http.StatusBadRequest,
			Body:       `{"message":"Geohash is not specified"}`,
		}, nil
	}

	var geoInfo GeoInfo
	if err := json.Unmarshal([]byte(req.Body), &geoInfo); err != nil {
		errResp := decodeErrResp
		errResp.StatusCode = http.StatusBadRequest
		return errResp, nil
	}

	if geoInfo.VendorIDs == nil || len(geoInfo.VendorIDs) == 0 {
		return Response{
			Headers:    map[string]string{"Content-Type": "application/json"},
			StatusCode: http.StatusBadRequest,
			Body:       `{"message":"No VendorIDs"}`,
		}, nil
	}
	logger.Info(geoInfo.VendorIDs)

	if err := h.Enqueue(geohash, geoInfo.VendorIDs); err != nil {
		return Response{
			Headers:    map[string]string{"Content-Type": "application/json"},
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"Failed to enqueue geoInfo"}`,
		}, nil
	}

	byts, err := json.Marshal(map[string]string{
		"message": "Successfully enqueued geoInfo",
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
