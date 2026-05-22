package payload

type Status struct {
	ApiVersion   string  `json:"api_version"`
	ConnectionId float64 `json:"connection_id"`
	System       string  `json:"system"`
	Version      string  `json:"version"`
}
