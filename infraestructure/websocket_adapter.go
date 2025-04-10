package infrastructure

import (
	"encoding/json"
	"log"
	"net/http"
	"weebsocket/application"
	"weebsocket/domain"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketAdapter struct {
	WebSocketPort application.WebSocketPort
	Connections  map[*websocket.Conn]bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewWebSocketAdapter(service application.WebSocketPort) *WebSocketAdapter {
	return &WebSocketAdapter{
		WebSocketPort: service,
		Connections:   make(map[*websocket.Conn]bool),
	}
}

func (wa *WebSocketAdapter) HandleWebSocket(ctx *gin.Context, sensorName string) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("❌ [%s] Error al actualizar a WebSocket: %v", sensorName, err)
		return
	}
	defer conn.Close()

	// Esperar mensaje de WebSocket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("⚠️ [%s] Error de lectura: %v", sensorName, err)
			break
		}

		switch sensorName {
		case "temperature":
			var data domain.TemperatureHumidity
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Printf("❗ [%s] Error al parsear el mensaje: %v", sensorName, err)
				continue
			}
			wa.WebSocketPort.HandleSensorData(sensorName, data)
		case "humidity":
			var data domain.TemperatureHumidity
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Printf("❗ [%s] Error al parsear el mensaje: %v", sensorName, err)
				continue
			}
			wa.WebSocketPort.HandleSensorData(sensorName, data)
		case "light":
			var data domain.Light
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Printf("❗ [%s] Error al parsear el mensaje: %v", sensorName, err)
				continue
			}
			wa.WebSocketPort.HandleSensorData(sensorName, data)
		case "sound":
			var data domain.SoundSensor
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Printf("❗ [%s] Error al parsear el mensaje: %v", sensorName, err)
				continue
			}
			wa.WebSocketPort.HandleSensorData(sensorName, data)
		case "air":
			var data domain.AirQualitySensor
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Printf("❗ [%s] Error al parsear el mensaje: %v", sensorName, err)
				continue
			}
			wa.WebSocketPort.HandleSensorData(sensorName, data)
		default:
			log.Printf("❗ Sensor no reconocido: %s", sensorName)
		}

		// Responder con un mensaje de éxito
		response := map[string]string{"status": "success", "message": "Datos recibidos correctamente"}
		responsePayload, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, responsePayload)
	}
}