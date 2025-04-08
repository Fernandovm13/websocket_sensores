package main

import (
	"fmt"
	"weebsocket/application"
	infrastructure "weebsocket/infraestructure"

	"github.com/gin-gonic/gin"
)

func main() {
	wsService := application.NewWebSocketService()

	wsAdapter := infrastructure.NewWebSocketAdapter(wsService)

	r := gin.Default()

	infrastructure.SetWebSocketRoutes(r, wsAdapter)

	if err := r.Run(":8084"); err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
