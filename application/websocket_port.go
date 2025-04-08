package application

import "weebsocket/domain"

type WebSocketPort interface {
	HandleSensorData(sensorName string, data domain.SensorData) error
	
	BroadcastMessage(message []byte) error
}
