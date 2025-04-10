package domain

type SensorData interface {
	GetSensorID() string
	GetTimestamp() string
}
