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


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketAdapter struct {
	WebSocketPort application.WebSocketPort
	Connections  map[*websocket.Conn]bool
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

	wa.Connections[conn] = true

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("⚠️ [%s] Error de lectura: %v", sensorName, err)
			break
		}

		var sensorData domain.SensorData
		if err := json.Unmarshal(p, &sensorData); err != nil {
			log.Printf("❗ [%s] Error al parsear el mensaje: %v", sensorName, err)
			continue
		}

		if err := wa.WebSocketPort.HandleSensorData(sensorName, sensorData); err != nil {
			log.Printf("❌ [%s] Error al procesar datos del sensor: %v", sensorName, err)
			continue
		}

		response := map[string]string{
			"status":  "success",
			"message": "Datos recibidos correctamente",
		}
		responsePayload, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, responsePayload)
	}
}

func (wa *WebSocketAdapter) BroadcastMessage(message []byte) {
	for conn := range wa.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("❌ Error enviando mensaje a cliente: %v", err)
			conn.Close()
			delete(wa.Connections, conn)
		}
	}
}
