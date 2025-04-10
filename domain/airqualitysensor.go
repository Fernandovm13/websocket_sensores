package domain

type AirQualitySensor struct {
    SensorID  string  `json:"SensorID"`
    CO2PPM    int     `json:"CO2PPM"`
    AirLevel  int     `json:"Air_level"`
    Timestamp string  `json:"Timestamp"`
}
func (a AirQualitySensor) GetSensorID() string {
	return a.SensorID
}

func (a AirQualitySensor) GetTimestamp() string {
	return a.Timestamp
}