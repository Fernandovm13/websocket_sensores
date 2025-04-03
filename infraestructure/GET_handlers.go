package infrastructure

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// En producción, restringe CheckOrigin
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// ===== SENSOR TEMPERATURE =====

var (
	activeTemperatureConnections = make(map[*websocket.Conn]bool)
	temperatureConnectionsMutex  = &sync.Mutex{}
	temperatureMessages          []int
	temperatureMessagesMutex     = &sync.Mutex{}
)

func HandleWSTemperature(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Temperature - Error al actualizar a WebSocket: %v", err)
		return
	}

	temperatureConnectionsMutex.Lock()
	activeTemperatureConnections[conn] = true
	temperatureConnectionsMutex.Unlock()

	log.Printf("Temperature - Conexión establecida desde: %s", conn.RemoteAddr())

	defer func() {
		temperatureConnectionsMutex.Lock()
		delete(activeTemperatureConnections, conn)
		temperatureConnectionsMutex.Unlock()
		conn.Close()
		log.Printf("Temperature - Conexión cerrada: %s", conn.RemoteAddr())
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Temperature - Error de lectura: %v", err)
			}
			break
		}

		messageStr := string(p)
		log.Printf("Temperature - Mensaje recibido (RAW): %s", messageStr)

		num, err := strconv.Atoi(messageStr)
		if err != nil {
			log.Printf("Temperature - Mensaje no numérico: %s", messageStr)
			conn.WriteMessage(websocket.TextMessage, []byte("ERROR: Se esperaba un número"))
			continue
		}

		temperatureMessagesMutex.Lock()
		temperatureMessages = append(temperatureMessages, num)
		if len(temperatureMessages) > 100 {
			temperatureMessages = temperatureMessages[1:]
		}
		temperatureMessagesMutex.Unlock()

		log.Printf("Temperature - Número procesado y almacenado: %d", num)

		if err := conn.WriteMessage(messageType, []byte(strconv.Itoa(num))); err != nil {
			log.Printf("Temperature - Error al enviar respuesta: %v", err)
			break
		}
	}
}

func ListTemperatureMessages(ctx *gin.Context) {
	temperatureMessagesMutex.Lock()
	defer temperatureMessagesMutex.Unlock()

	history := temperatureMessages
	if history == nil {
		history = []int{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messages": history,
		"count":    len(history),
		"status":   "success",
	})
}

// ===== SENSOR NOISE =====

var (
	activeNoiseConnections = make(map[*websocket.Conn]bool)
	noiseConnectionsMutex  = &sync.Mutex{}
	noiseMessages          []int
	noiseMessagesMutex     = &sync.Mutex{}
)

func HandleWSNoise(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Noise - Error al actualizar a WebSocket: %v", err)
		return
	}

	noiseConnectionsMutex.Lock()
	activeNoiseConnections[conn] = true
	noiseConnectionsMutex.Unlock()

	log.Printf("Noise - Conexión establecida desde: %s", conn.RemoteAddr())

	defer func() {
		noiseConnectionsMutex.Lock()
		delete(activeNoiseConnections, conn)
		noiseConnectionsMutex.Unlock()
		conn.Close()
		log.Printf("Noise - Conexión cerrada: %s", conn.RemoteAddr())
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Noise - Error de lectura: %v", err)
			}
			break
		}

		messageStr := string(p)
		log.Printf("Noise - Mensaje recibido (RAW): %s", messageStr)

		num, err := strconv.Atoi(messageStr)
		if err != nil {
			log.Printf("Noise - Mensaje no numérico: %s", messageStr)
			conn.WriteMessage(websocket.TextMessage, []byte("ERROR: Se esperaba un número"))
			continue
		}

		noiseMessagesMutex.Lock()
		noiseMessages = append(noiseMessages, num)
		if len(noiseMessages) > 100 {
			noiseMessages = noiseMessages[1:]
		}
		noiseMessagesMutex.Unlock()

		log.Printf("Noise - Número procesado y almacenado: %d", num)

		if err := conn.WriteMessage(messageType, []byte(strconv.Itoa(num))); err != nil {
			log.Printf("Noise - Error al enviar respuesta: %v", err)
			break
		}
	}
}

func ListNoiseMessages(ctx *gin.Context) {
	noiseMessagesMutex.Lock()
	defer noiseMessagesMutex.Unlock()

	history := noiseMessages
	if history == nil {
		history = []int{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messages": history,
		"count":    len(history),
		"status":   "success",
	})
}

// ===== SENSOR LIGHT =====

var (
	activeLightConnections = make(map[*websocket.Conn]bool)
	lightConnectionsMutex  = &sync.Mutex{}
	lightMessages          []int
	lightMessagesMutex     = &sync.Mutex{}
)

func HandleWSLight(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Light - Error al actualizar a WebSocket: %v", err)
		return
	}

	lightConnectionsMutex.Lock()
	activeLightConnections[conn] = true
	lightConnectionsMutex.Unlock()

	log.Printf("Light - Conexión establecida desde: %s", conn.RemoteAddr())

	defer func() {
		lightConnectionsMutex.Lock()
		delete(activeLightConnections, conn)
		lightConnectionsMutex.Unlock()
		conn.Close()
		log.Printf("Light - Conexión cerrada: %s", conn.RemoteAddr())
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Light - Error de lectura: %v", err)
			}
			break
		}

		messageStr := string(p)
		log.Printf("Light - Mensaje recibido (RAW): %s", messageStr)

		num, err := strconv.Atoi(messageStr)
		if err != nil {
			log.Printf("Light - Mensaje no numérico: %s", messageStr)
			conn.WriteMessage(websocket.TextMessage, []byte("ERROR: Se esperaba un número"))
			continue
		}

		lightMessagesMutex.Lock()
		lightMessages = append(lightMessages, num)
		if len(lightMessages) > 100 {
			lightMessages = lightMessages[1:]
		}
		lightMessagesMutex.Unlock()

		log.Printf("Light - Número procesado y almacenado: %d", num)

		if err := conn.WriteMessage(messageType, []byte(strconv.Itoa(num))); err != nil {
			log.Printf("Light - Error al enviar respuesta: %v", err)
			break
		}
	}
}

func ListLightMessages(ctx *gin.Context) {
	lightMessagesMutex.Lock()
	defer lightMessagesMutex.Unlock()

	history := lightMessages
	if history == nil {
		history = []int{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messages": history,
		"count":    len(history),
		"status":   "success",
	})
}

// ===== SENSOR AIR (Calidad del aire) =====

var (
	activeAirConnections = make(map[*websocket.Conn]bool)
	airConnectionsMutex  = &sync.Mutex{}
	airMessages          []int
	airMessagesMutex     = &sync.Mutex{}
)

func HandleWSAir(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Air - Error al actualizar a WebSocket: %v", err)
		return
	}

	airConnectionsMutex.Lock()
	activeAirConnections[conn] = true
	airConnectionsMutex.Unlock()

	log.Printf("Air - Conexión establecida desde: %s", conn.RemoteAddr())

	defer func() {
		airConnectionsMutex.Lock()
		delete(activeAirConnections, conn)
		airConnectionsMutex.Unlock()
		conn.Close()
		log.Printf("Air - Conexión cerrada: %s", conn.RemoteAddr())
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Air - Error de lectura: %v", err)
			}
			break
		}

		messageStr := string(p)
		log.Printf("Air - Mensaje recibido (RAW): %s", messageStr)

		num, err := strconv.Atoi(messageStr)
		if err != nil {
			log.Printf("Air - Mensaje no numérico: %s", messageStr)
			conn.WriteMessage(websocket.TextMessage, []byte("ERROR: Se esperaba un número"))
			continue
		}

		airMessagesMutex.Lock()
		airMessages = append(airMessages, num)
		if len(airMessages) > 100 {
			airMessages = airMessages[1:]
		}
		airMessagesMutex.Unlock()

		log.Printf("Air - Número procesado y almacenado: %d", num)

		if err := conn.WriteMessage(messageType, []byte(strconv.Itoa(num))); err != nil {
			log.Printf("Air - Error al enviar respuesta: %v", err)
			break
		}
	}
}

func ListAirMessages(ctx *gin.Context) {
	airMessagesMutex.Lock()
	defer airMessagesMutex.Unlock()

	history := airMessages
	if history == nil {
		history = []int{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"messages": history,
		"count":    len(history),
		"status":   "success",
	})
}
