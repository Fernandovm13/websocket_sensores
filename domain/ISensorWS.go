package domain

import "github.com/gin-gonic/gin"

type IsensorWS interface {
	HandleWebSocket(ctx *gin.Context, sensorName string)
	ListMessages(ctx *gin.Context, sensorName string)
}

type AnomalyDetectionPort interface {
	Detect(sensorName string, data SensorData)
}