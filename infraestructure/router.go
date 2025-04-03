package infrastructure

import (
	"github.com/gin-gonic/gin"
)

func SetWebSocketRoutes(engine *gin.Engine) {
	wsGroup := engine.Group("/ws")
	{
		// Rutas para el WebSocket (handshake)
		wsGroup.GET("/handshake/temperature", HandleWSTemperature)
		wsGroup.GET("/handshake/noise", HandleWSNoise)
		wsGroup.GET("/handshake/light", HandleWSLight)
		wsGroup.GET("/handshake/air", HandleWSAir)

		// Rutas para consultar el historial de mensajes de cada sensor
		wsGroup.GET("/messages/temperature", ListTemperatureMessages)
		wsGroup.GET("/messages/noise", ListNoiseMessages)
		wsGroup.GET("/messages/light", ListLightMessages)
		wsGroup.GET("/messages/air", ListAirMessages)
	}
}
