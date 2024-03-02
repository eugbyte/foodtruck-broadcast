package subscription

import (
	"encoding/json"
	"testing"
	"time"

	assrt "github.com/stretchr/testify/assert"
)

type MockSubRepo struct {
	Subs []Subscription
}

func NewMockSubRepo() *MockSubRepo {
	return &MockSubRepo{
		Subs: make([]Subscription, 0),
	}
}

func (m *MockSubRepo) Open() {}
func (m *MockSubRepo) Put(sub Subscription) error {
	m.Subs = append(m.Subs, sub)
	return nil
}

func TestHandler(t *testing.T) {
	var assert = assrt.New(t)
	mockSubRepo := NewMockSubRepo()

	handler := New(mockSubRepo)

	timestamp := time.Now().Unix()

	var subscription = Subscription{
		Endpoint: "www.abc.com",
		Geohash:  "abcde",
		LastSend: timestamp,
	}

	byts, err := json.Marshal(subscription)
	assert.Nil(err)
	req := Request{
		Body: string(byts),
	}

	resp, err := handler.Handle(req)
	assert.Nil(err)

	sub := mockSubRepo.Subs[0]

	assert.Equal(subscription, sub)

	var expRespmap map[string]string
	err = json.Unmarshal([]byte(resp.Body), &expRespmap)
	assert.Nil(err)
	assert.NotNil(expRespmap)

	assert.Equal(expRespmap["message"], "subscription saved")
}
