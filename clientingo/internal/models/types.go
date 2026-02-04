package models

// DeviceData represents the telemetry data for a single device
type DeviceData struct {
	SN          string `json:"sn"`
	Power       string `json:"power"`
	Status      string `json:"status"`
	LastUpdated string `json:"last_update"`
}

// APIRequest represents the request payload
type APIRequest struct {
	SNList []string `json:"sn_list"`
}

// APIResponse represents the response from the API
type APIResponse struct {
	Data  []DeviceData `json:"data"`
	Error string       `json:"error,omitempty"`
}
