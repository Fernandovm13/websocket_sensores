package application

import (
	"fmt"
	"log"
	"weebsocket/domain"
)

type WebSocketService struct {
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{}
}

func (ws *WebSocketService) HandleSensorData(sensorName string, data domain.SensorData) error {

	if err := validateSensorData(data); err != nil {
		log.Printf("❌ [%s] Datos inválidos: %v", sensorName, err)
		return fmt.Errorf("datos inválidos para el sensor %s: %v", sensorName, err)
	}

	logSensorData(sensorName, data)

	log.Printf("🟢 [%s] Datos procesados correctamente.", sensorName)
	return nil
}

func (ws *WebSocketService) BroadcastMessage(message []byte) error {
	log.Printf("📡 Enviando mensaje de difusión: %s", string(message))
	return nil
}




func logSensorData(sensorName string, data domain.SensorData) {
	log.Printf("📊 [%s] Información del sensor:", sensorName)
	log.Printf("🌡️ Temperatura: %.2f°C", data.Temperature)
	log.Printf("💧 Humedad: %.2f%%", data.Humidity)
	log.Printf("💡 Luz: %d lux", data.Light)
	log.Printf("🔊 Sonido: %d dB", data.Sound)
	log.Printf("💨 CO2: %.2f ppm", data.CO2)
	log.Printf("🌍 Timestamp: %s", data.Timestamp)
}

func validateSensorData(data domain.SensorData) error {
	if data.Timestamp == "" {
		return fmt.Errorf("el timestamp es obligatorio")
	}
	if data.Temperature < 14 || data.Temperature > 28 {
		return fmt.Errorf("temperatura fuera de rango")
	}
	if data.Humidity < 0 || data.Humidity > 100 {
		return fmt.Errorf("humedad fuera de rango")
	}
	if data.Light < 0 || data.Light > 10000 {
		return fmt.Errorf("luz fuera de rango")
	}
	return nil
}
