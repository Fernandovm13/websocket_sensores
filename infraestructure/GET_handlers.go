package infrastructure

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type SensorData struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Light       int     `json:"light"`
	Sound       int     `json:"sound"`
	Air         int     `json:"air"`
	CO2         float64 `json:"co2"`
	Timestamp   string  `json:"timestamp"`
}

type Sensor struct {
	Connections      map[*websocket.Conn]bool
	ConnectionsMutex *sync.Mutex
	Messages         []SensorData
	MessagesMutex    *sync.Mutex
}

const (
	temperatureThreshold = 30.0  // Umbral de temperatura para anomalías
	humidityThreshold    = 70.0  // Umbral de humedad para anomalías
	lightThreshold       = 800   // Umbral de luz para anomalías
	noiseThreshold       = 75    // Umbral de ruido para anomalías
	airThreshold         = 1000  // Umbral de calidad del aire para anomalías
)

func NewSensor() *Sensor {
	return &Sensor{
		Connections:      make(map[*websocket.Conn]bool),
		ConnectionsMutex: &sync.Mutex{},
		Messages:         []SensorData{},
		MessagesMutex:    &sync.Mutex{},
	}
}

// Función para detectar anomalías en los datos del sensor
func detectAnomalies(sensorName string, data SensorData) {
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

	default:
		log.Printf("❗ [%s] Anomalía detectada en los datos del sensor (sin filtro): %+v", sensorName, data)
	}
}

func handleWebSocket(ctx *gin.Context, sensorName string, sensor *Sensor) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("❌ [%s] Error al actualizar a WebSocket: %v", sensorName, err)
		return
	}

	sensor.ConnectionsMutex.Lock()
	sensor.Connections[conn] = true
	sensor.ConnectionsMutex.Unlock()

	log.Printf("✅ [%s] Conexión establecida desde: %s", sensorName, conn.RemoteAddr())

	defer func() {
		sensor.ConnectionsMutex.Lock()
		delete(sensor.Connections, conn)
		sensor.ConnectionsMutex.Unlock()
		conn.Close()
		log.Printf("🔌 [%s] Conexión cerrada: %s", sensorName, conn.RemoteAddr())
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("⚠️ [%s] Error de lectura: %v", sensorName, err)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			continue
		}

		var sensorData SensorData
		if err := json.Unmarshal(p, &sensorData); err != nil {
			log.Printf("❗ [%s] Error al parsear el mensaje: %v", sensorName, err)
			continue
		}

		// Log solo los datos del sensor específico
		switch sensorName {
		case "Air":
			log.Printf("📩 [%s] Datos recibidos: Air: %d, Timestamp: %s", sensorName, sensorData.Air, sensorData.Timestamp)
		case "Temperature":
			log.Printf("📩 [%s] Datos recibidos: Temperature: %.2f, Timestamp: %s", sensorName, sensorData.Temperature, sensorData.Timestamp)
		case "Humidity":
			log.Printf("📩 [%s] Datos recibidos: Humidity: %.2f, Timestamp: %s", sensorName, sensorData.Humidity, sensorData.Timestamp)
		case "Light":
			log.Printf("📩 [%s] Datos recibidos: Light: %d, Timestamp: %s", sensorName, sensorData.Light, sensorData.Timestamp)
		case "Noise":
			log.Printf("📩 [%s] Datos recibidos: Sound: %d, Timestamp: %s", sensorName, sensorData.Sound, sensorData.Timestamp)
		}

		// Detectar anomalías antes de guardar los datos
		detectAnomalies(sensorName, sensorData)

		// Guardar los datos sin necesidad de filtrado extra
		sensor.MessagesMutex.Lock()
		sensor.Messages = append(sensor.Messages, sensorData)
		if len(sensor.Messages) > 100 {
			sensor.Messages = sensor.Messages[1:]
		}
		sensor.MessagesMutex.Unlock()

		// Enviar el mensaje de vuelta al WebSocket
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Printf("⚠️ [%s] Error al enviar respuesta: %v", sensorName, err)
			break
		}
	}
}

func listMessages(ctx *gin.Context, sensorName string, sensor *Sensor) {
	sensor.MessagesMutex.Lock()
	defer sensor.MessagesMutex.Unlock()

	history := sensor.Messages
	if history == nil {
		history = []SensorData{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sensor":   sensorName,
		"messages": history,
		"count":    len(history),
		"status":   "success",
	})

	log.Printf("📦 [%s] Historial solicitado. Total: %d", sensorName, len(history))
}


// ====== TEMPERATURE ======
var temperatureSensor = NewSensor()

func HandleWSTemperature(ctx *gin.Context) {
	handleWebSocket(ctx, "Temperature", temperatureSensor)
}

func ListTemperatureMessages(ctx *gin.Context) {
	listMessages(ctx, "Temperature", temperatureSensor)
}

// ====== NOISE ======
var noiseSensor = NewSensor()

func HandleWSNoise(ctx *gin.Context) {
	handleWebSocket(ctx, "Noise", noiseSensor)
}

func ListNoiseMessages(ctx *gin.Context) {
	listMessages(ctx, "Noise", noiseSensor)
}

// ====== LIGHT ======
var lightSensor = NewSensor()

func HandleWSLight(ctx *gin.Context) {
	handleWebSocket(ctx, "Light", lightSensor)
}

func ListLightMessages(ctx *gin.Context) {
	listMessages(ctx, "Light", lightSensor)
}

// ====== AIR ======
var airSensor = NewSensor()

func HandleWSAir(ctx *gin.Context) {
	handleWebSocket(ctx, "Air", airSensor)
}

func ListAirMessages(ctx *gin.Context) {
	listMessages(ctx, "Air", airSensor)
}
