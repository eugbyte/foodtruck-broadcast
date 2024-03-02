package main

import (
	"flag"
	debug "foodtruck/pkg/logger"
	"foodtruck/pkg/model"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samber/lo"
)

var addr = flag.String("addr", "localhost:6000", "http service address")
var logger = debug.Logger

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/vendor"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			return
		case <-ticker.C:
			lng := 103.8198 + randomFloat(0.01, 0.09) // (round to nearest)
			lat := 1.3521 + randomFloat(0.001, 0.009) // (round to nearest)

			geoInfo := model.GeoInfo{
				VendorID: strconv.Itoa(randInt(1, 3)),
				Lat:      lat,
				Lng:      lng,
				Speed:    lo.ToPtr(500.0),
			}
			if err := conn.WriteJSON(geoInfo); err != nil {
				logger.Error(err)
			}
		}
	}
}

func randomFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
