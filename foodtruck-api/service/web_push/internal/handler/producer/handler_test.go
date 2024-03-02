package producer

import (
	"encoding/json"
	"net/http"
	"testing"

	assrt "github.com/stretchr/testify/assert"
)

type MockMsg struct {
	Msg      string
	Metadata map[string]string
}

type MockMsgQ struct {
	Msgs []MockMsg
}

func NewMockMsgQ() *MockMsgQ {
	return &MockMsgQ{
		Msgs: make([]MockMsg, 0),
	}
}

func (m *MockMsgQ) Open() error { return nil }
func (m *MockMsgQ) Enqueue(msg string, metadata map[string]string) error {
	m.Msgs = append(m.Msgs, MockMsg{
		Msg:      msg,
		Metadata: metadata,
	})
	return nil
}

func TestEnqueue(t *testing.T) {
	var assert = assrt.New(t)
	mockQ := NewMockMsgQ()

	vendorIDs := []string{"abc", "def"}

	handler := New(mockQ)
	err := handler.Enqueue("mock_geohash", vendorIDs)
	assert.Nil(err)

	msg := mockQ.Msgs[0]
	assert.Equal(msg.Metadata["geohash"], "mock_geohash")

	var actVendorIDs []string
	err = json.Unmarshal([]byte(msg.Msg), &actVendorIDs)
	assert.Nil(err)
	assert.Equal(actVendorIDs, vendorIDs)
}

func TestHandler(t *testing.T) {
	var assert = assrt.New(t)
	mockQ := NewMockMsgQ()

	vendorIDs := []string{"abc", "def"}
	var geoInfo = GeoInfo{
		VendorIDs: vendorIDs,
	}
	byts, err := json.Marshal(geoInfo)
	assert.Nil(err)

	req := Request{
		PathParameters: map[string]string{
			"geohash": "mock_geohash",
		},
		Body: string(byts),
	}

	handler := New(mockQ)
	resp, err := handler.Handle(req)
	assert.Nil(err)
	assert.Equal(resp.StatusCode, http.StatusOK)
}

func TestHandlerValidation(t *testing.T) {
	var assert = assrt.New(t)
	mockQ := NewMockMsgQ()

	vendorIDs := []string{"abc", "def"}
	var geoInfo = GeoInfo{
		VendorIDs: vendorIDs,
	}
	byts, err := json.Marshal(geoInfo)
	assert.Nil(err)

	req := Request{
		PathParameters: map[string]string{
			"geohash": "",
		},
		Body: string(byts),
	}

	handler := New(mockQ)
	resp, err := handler.Handle(req)
	assert.Nil(err)
	assert.Equal(resp.StatusCode, http.StatusBadRequest)

	var errMsg map[string]string
	err = json.Unmarshal([]byte(resp.Body), &errMsg)
	assert.Nil(err)
	assert.Equal(errMsg["message"], "Geohash is not specified")
}
