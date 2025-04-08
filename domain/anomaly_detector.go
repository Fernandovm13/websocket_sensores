package domain

import "log"

const (
	temperatureThreshold = 30.0
	humidityThreshold    = 70.0
	lightThreshold       = 800
	noiseThreshold       = 75
	airThreshold         = 1000
)

type AnomalyDetector struct{}

func (ad *AnomalyDetector) Detect(sensorName string, data SensorData) {
	switch sensorName {
	case "Temperature":
		if data.Temperature > temperatureThreshold {
			log.Printf("❗ [%s] Anomalía detectada: Temperatura demasiado alta (%.2f°C)!", sensorName, data.Temperature)
		} else if data.Temperature < 0 {
			log.Printf("❗ [%s] Anomalía detectada: Temperatura demasiado baja (%.2f°C)!", sensorName, data.Temperature)
		}

	case "Humidity":
		if data.Humidity > humidityThreshold {
			log.Printf("❗ [%s] Anomalía detectada: Humedad demasiado alta (%.2f%%)!", sensorName, data.Humidity)
		}

	case "Light":
		if data.Light > lightThreshold {
			log.Printf("❗ [%s] Anomalía detectada: Luz demasiado alta (%d lux)!", sensorName, data.Light)
		}

	case "Noise":
		if data.Sound > noiseThreshold {
			log.Printf("❗ [%s] Anomalía detectada: Ruido demasiado alto (%d dB)!", sensorName, data.Sound)
		}

	case "Air":
		if data.Air > airThreshold {
			log.Printf("❗ [%s] Anomalía detectada: Calidad del aire demasiado baja (%d)!", sensorName, data.Air)
		}
	}
}
