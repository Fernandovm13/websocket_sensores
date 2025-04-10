package application

import (
	"fmt"
	"log"
	"weebsocket/domain"
)

type WebSocketService struct{}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{}
}

func (ws *WebSocketService) HandleSensorData(sensorName string, data domain.SensorData) error {
	switch v := data.(type) {
	case domain.TemperatureHumidity:
		log.Printf("Procesando datos de temperatura y humedad: %+v", v)
	case domain.Light:
		log.Printf("Procesando datos de luz: %+v", v)
	case domain.SoundSensor:
		log.Printf("Procesando datos de sonido: %+v", v)
	case domain.AirQualitySensor:
		log.Printf("Procesando datos de calidad del aire: %+v", v)
	default:
		return fmt.Errorf("tipo de datos desconocido: %T", v)
	}
	return nil
}

func (ws *WebSocketService) BroadcastMessage(message []byte) error {
    log.Printf("Broadcasting message: %s", string(message))
    return nil
}