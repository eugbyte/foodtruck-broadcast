package mock

import (
	"encoding/json"
	"foodtruck/pkg/model"
	"strconv"
	"time"

	kafkalib "foodtruck/pkg/queue/kafkalib/reader"
	"math/rand"

	"github.com/samber/lo"
	"github.com/segmentio/kafka-go"
)

type MockReader struct {
}

func NewMockReader() *MockReader {
	return &MockReader{}
}

func (m *MockReader) Open()                            {}
func (m *MockReader) Close() error                     { return nil }
func (m *MockReader) SetKafkaOffset(t time.Time) error { return nil }

func (m *MockReader) Read() (kafkalib.Message, error) {
	time.Sleep(2 * time.Second)

	lng := 103.8198 + randomFloat(0.01, 0.09) // (round to nearest)
	lat := 1.3521 + randomFloat(0.001, 0.009) // (round to nearest)

	geoInfo := model.GeoInfo{
		VendorID: strconv.Itoa(randInt(1, 3)),
		Lat:      lat,
		Lng:      lng,
		Speed:    lo.ToPtr(500.0),
	}
	byts, err := json.Marshal(geoInfo)
	return kafka.Message{
		Value: byts,
	}, err
}

func randomFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
