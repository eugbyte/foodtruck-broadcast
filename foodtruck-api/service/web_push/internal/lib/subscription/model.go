package subscription

type Subscription struct {
	Endpoint   string `json:"endpoint" dynamo:"endpoint,hash"`                       // hash key
	Geohash    string `json:"geohash" dynamo:"geohash" index:"geohashIndex,hash"`    // GSI
	LastSend   int64  `json:"lastSend" dynamo:"lastSend" index:"geohashIndex,range"` // GSI
	Expiration int64  `json:"expiration" dynamo:"expiration"`
	Auth       string `json:"auth" dynamo:"auth"`
	P256dh     string `json:"p256dh" dynamo:"p256dh"` // public key
	OptIn      bool   `json:"optIn" dynamo:"optIn"`
}
