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

// WebSocketAdapter es el adaptador que maneja la infraestructura de WebSocket
type WebSocketAdapter struct {
	WebSocketPort application.WebSocketPort
	Connections  map[*websocket.Conn]bool
}

// NewWebSocketAdapter crea un nuevo WebSocketAdapter
func NewWebSocketAdapter(service application.WebSocketPort) *WebSocketAdapter {
	return &WebSocketAdapter{
		WebSocketPort: service,
		Connections:   make(map[*websocket.Conn]bool),
	}
}

// HandleWebSocket maneja la conexión WebSocket para un sensor específico
func (wa *WebSocketAdapter) HandleWebSocket(ctx *gin.Context, sensorName string) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("❌ [%s] Error al actualizar a WebSocket: %v", sensorName, err)
		return
	}

	// Añadir la conexión a la lista de conexiones activas
	wa.Connections[conn] = true

	// Leer los mensajes enviados por el cliente
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("⚠️ [%s] Error de lectura: %v", sensorName, err)
			break
		}

		// Deserializar los datos en SensorData
		var sensorData domain.SensorData
		if err := json.Unmarshal(p, &sensorData); err != nil {
			log.Printf("❗ [%s] Error al parsear el mensaje: %v", sensorName, err)
			continue
		}

		// Delegar la lógica de negocio a la capa de aplicación
		if err := wa.WebSocketPort.HandleSensorData(sensorName, sensorData); err != nil {
			log.Printf("❌ [%s] Error al procesar datos del sensor: %v", sensorName, err)
			continue
		}

		// Responder al cliente con un mensaje de éxito
		response := map[string]string{
			"status":  "success",
			"message": "Datos recibidos correctamente",
		}
		responsePayload, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, responsePayload)
	}
}

// Función para enviar un mensaje a todos los clientes conectados
func (wa *WebSocketAdapter) BroadcastMessage(message []byte) {
	for conn := range wa.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("❌ Error enviando mensaje a cliente: %v", err)
			conn.Close()
			delete(wa.Connections, conn)
		}
	}
}
