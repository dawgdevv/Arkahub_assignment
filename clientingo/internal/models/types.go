package models

type DeviceData struct {
	SN          string `json:"sn"`
	Power       string `json:"power"`
	Status      string `json:"status"`
	LastUpdated string `json:"last_update"`
}

type APIRequest struct {
	SNList []string `json:"sn_list"`
}

type APIResponse struct {
	Data  []DeviceData `json:"data"`
	Error string       `json:"error,omitempty"`
}
