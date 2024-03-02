package model

type GeoInfo struct {
	VendorID string `json:"vendorID"`
	// latitude: https://developer.mozilla.org/en-US/docs/Web/API/GeolocationCoordinates/latitude
	Lat float64 `json:"lat"`
	// longitude: https://developer.mozilla.org/en-US/docs/Web/API/GeolocationCoordinates/longitude
	Lng float64 `json:"lng"`
	// speed, m/s: https://developer.mozilla.org/en-US/docs/Web/API/GeolocationCoordinates/speed
	Speed *float64 `json:"speed"`
	// heading from true north https://developer.mozilla.org/en-US/docs/Web/API/GeolocationCoordinates/heading
	Heading *float64 `json:"heading"`
}
