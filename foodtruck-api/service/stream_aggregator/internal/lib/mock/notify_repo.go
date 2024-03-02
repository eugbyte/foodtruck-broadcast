package mock

import "encoding/json"

type MockNotifyAPI struct{}

func NewMockNotifyAPI() *MockNotifyAPI {
	return &MockNotifyAPI{}
}

func (m *MockNotifyAPI) Post(payload map[string]any) (respBody []byte, err error) {
	return json.Marshal(map[string]string{
		"message": "lambda successfully invoked",
	})
}
