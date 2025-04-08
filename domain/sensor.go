package domain

type SensorData struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Light       int     `json:"light"`
	Sound       int     `json:"sound"`
	Air         int     `json:"air"`
	CO2         float64 `json:"co2"`
	Timestamp   string  `json:"timestamp"`
}