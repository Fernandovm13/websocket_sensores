package infrastructure

import (
	"github.com/gin-gonic/gin"
)

// SetWebSocketRoutes configura las rutas WebSocket en el servidor Gin
func SetWebSocketRoutes(engine *gin.Engine, webSocketAdapter *WebSocketAdapter) {
	wsGroup := engine.Group("/ws")
	{
		// Rutas para WebSocket (handshake)
		wsGroup.GET("/handshake/temperature", func(c *gin.Context) {
			webSocketAdapter.HandleWebSocket(c, "temperature")
		})
		wsGroup.GET("/handshake/humidity", func(c *gin.Context) {
			webSocketAdapter.HandleWebSocket(c, "humidity")
		})
		wsGroup.GET("/handshake/light", func(c *gin.Context) {
			webSocketAdapter.HandleWebSocket(c, "light")
		})
		wsGroup.GET("/handshake/sound", func(c *gin.Context) {
			webSocketAdapter.HandleWebSocket(c, "sound")
		})
		wsGroup.GET("/handshake/air", func(c *gin.Context) {
			webSocketAdapter.HandleWebSocket(c, "air")
		})
	}
}
