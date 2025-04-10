package domain

type SoundSensor struct {
	ID        int    `json:"id"`
	SensorID  string `json:"sensor_id"`
	RuidoDB   int    `json:"nivel"`
	Timestamp string `json:"timestamp"`
}

func (s SoundSensor) GetSensorID() string {
	return s.SensorID
}

func (s SoundSensor) GetTimestamp() string {
	return s.Timestamp
}