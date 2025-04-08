package main

import (
	"log"
	"os"
	infrastructure "weebsocket/infraestructure"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Configurar Gin
	r := gin.Default()
	
	// Middleware para CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Configurar rutas
	infrastructure.SetWebSocketRoutes(r)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Cambiado a 8081 para coincidir con tu frontend
	}

	log.Printf("Servidor iniciado en :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
