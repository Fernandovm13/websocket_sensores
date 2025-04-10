package domain

type TemperatureHumidity struct {
	SensorID    string  `json:"sensor_id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Timestamp   string  `json:"timestamp"`
}

// Implementa los métodos de la interfaz SensorData
func (t TemperatureHumidity) GetSensorID() string {
	return t.SensorID
}

func (t TemperatureHumidity) GetTimestamp() string {
	return t.Timestamp
}