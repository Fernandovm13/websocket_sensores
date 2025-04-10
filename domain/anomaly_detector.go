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

// Detecta anomalías basadas en el tipo de sensor
func (ad *AnomalyDetector) Detect(sensorName string, data SensorData) {
	switch sensorName {
	case "Temperature":
		if t, ok := data.(TemperatureHumidity); ok {
			if t.Temperature > temperatureThreshold {
				log.Printf("❗ [%s] Anomalía detectada: Temperatura demasiado alta (%.2f°C)!", sensorName, t.Temperature)
			} else if t.Temperature < 0 {
				log.Printf("❗ [%s] Anomalía detectada: Temperatura demasiado baja (%.2f°C)!", sensorName, t.Temperature)
			}
		}
	case "Humidity":
		if t, ok := data.(TemperatureHumidity); ok {
			if t.Humidity > humidityThreshold {
				log.Printf("❗ [%s] Anomalía detectada: Humedad demasiado alta (%.2f%%)!", sensorName, t.Humidity)
			}
		}
	case "Light":
		if l, ok := data.(Light); ok {
			if l.Nivel > lightThreshold {
				log.Printf("❗ [%s] Anomalía detectada: Luz demasiado alta (%.2f lux)!", sensorName, l.Nivel)
			}
		}
	case "Noise":
		if s, ok := data.(SoundSensor); ok {
			if s.RuidoDB > noiseThreshold {
				log.Printf("❗ [%s] Anomalía detectada: Ruido demasiado alto (%d dB)!", sensorName, s.RuidoDB)
			}
		}
	case "Air":
		if a, ok := data.(AirQualitySensor); ok {
			if a.Air_level > airThreshold {
				log.Printf("❗ [%s] Anomalía detectada: Calidad del aire demasiado baja (%d)!", sensorName, a.Air_level)
			}
		}
	}
}