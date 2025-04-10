package domain

type Light struct {
	ID        int     `json:"id"`
	SensorID  string  `json:"sensor_id"`
	Nivel     float64 `json:"nivel"`
	Timestamp string  `json:"timestamp"`
}

func (l Light) GetSensorID() string {
	return l.SensorID
}

func (l Light) GetTimestamp() string {
	return l.Timestamp
}