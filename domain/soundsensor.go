package domain

type SoundSensor struct {
    SensorID  string  `json:"SensorID"`
    RuidoDB   int     `json:"RuidoDB"`
    Timestamp string  `json:"Timestamp"`
}

func (s SoundSensor) GetSensorID() string {
	return s.SensorID
}

func (s SoundSensor) GetTimestamp() string {
	return s.Timestamp
}