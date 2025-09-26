package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketMessage represents a WebSocket message structure
type WebSocketMessage struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp string                 `json:"timestamp"`
	Error     string                 `json:"error,omitempty"`
}

// WebSocketConnection represents a WebSocket connection with metadata
type WebSocketConnection struct {
	Conn     *websocket.Conn
	UserID   string
	UserType string
	LastPing time.Time
	Closed   int32 // Atomic flag for connection state
}

// IWebSocketManager defines the interface for WebSocket connection management
type IWebSocketManager interface {
	AddConnection(userID, userType string, conn *websocket.Conn)
	RemoveConnection(userID, userType string)
	SendMessage(userID, userType string, message WebSocketMessage) error
	BroadcastToType(userType string, message WebSocketMessage)
	BroadcastToUser(userType, userID string, message WebSocketMessage)
	StartPingPong(conn *WebSocketConnection)
	GetConnectionCount() int
	GetConnectionsByType(userType string) []*WebSocketConnection
	GetConnection(userID, userType string) *WebSocketConnection
	IsConnected(userID, userType string) bool
}

// WebSocketManager manages WebSocket connections
type WebSocketManager struct {
	connections     sync.Map
	connectionCount int64 // Atomic counter
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() IWebSocketManager {
	return &WebSocketManager{}
}

// AddConnection adds a new WebSocket connection
func (wm *WebSocketManager) AddConnection(userID, userType string, conn *websocket.Conn) {
	connectionID := userType + ":" + userID
	connection := &WebSocketConnection{
		Conn:     conn,
		UserID:   userID,
		UserType: userType,
		LastPing: time.Now(),
		Closed:   0, // Atomic flag, 0 = open
	}

	wm.connections.Store(connectionID, connection)
	atomic.AddInt64(&wm.connectionCount, 1)
	log.Printf("WebSocket connection added: %s", connectionID)
}

// RemoveConnection removes a WebSocket connection
func (wm *WebSocketManager) RemoveConnection(userID, userType string) {
	connectionID := userType + ":" + userID
	if connInterface, exists := wm.connections.LoadAndDelete(connectionID); exists {
		conn := connInterface.(*WebSocketConnection)
		atomic.StoreInt32(&conn.Closed, 1)
		atomic.AddInt64(&wm.connectionCount, -1)
		log.Printf("WebSocket connection removed: %s", connectionID)
	}
}

// SendMessage sends a message to a specific user
func (wm *WebSocketManager) SendMessage(userID, userType string, message WebSocketMessage) error {
	connectionID := userType + ":" + userID
	connInterface, exists := wm.connections.Load(connectionID)
	if !exists {
		return nil // Connection doesn't exist, silently ignore
	}

	conn := connInterface.(*WebSocketConnection)

	// Check if connection is closed using atomic operation
	if atomic.LoadInt32(&conn.Closed) == 1 {
		return nil
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v", err)
		return err
	}

	// Double-check if connection is still open
	if atomic.LoadInt32(&conn.Closed) == 1 {
		return nil
	}

	conn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err := conn.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		log.Printf("Failed to send WebSocket message to %s: %v", connectionID, err)
		atomic.StoreInt32(&conn.Closed, 1)
		return err
	}

	return nil
}

// BroadcastToType sends a message to all connections of a specific type
func (wm *WebSocketManager) BroadcastToType(userType string, message WebSocketMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal broadcast message: %v", err)
		return
	}

	wm.connections.Range(func(key, value interface{}) bool {
		connectionID := key.(string)
		conn := value.(*WebSocketConnection)

		if conn.UserType == userType && atomic.LoadInt32(&conn.Closed) == 0 {
			conn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				log.Printf("Failed to broadcast to %s: %v", connectionID, err)
				atomic.StoreInt32(&conn.Closed, 1)
			}
		}
		return true // Continue iteration
	})
}

// BroadcastToUser sends a message to a specific user (alias for SendMessage for consistency)
func (wm *WebSocketManager) BroadcastToUser(userType, userID string, message WebSocketMessage) {
	wm.SendMessage(userID, userType, message)
}

// StartPingPong starts ping-pong mechanism for connection health
func (wm *WebSocketManager) StartPingPong(conn *WebSocketConnection) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if atomic.LoadInt32(&conn.Closed) == 1 {
			return
		}

		conn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := conn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("Ping failed for %s:%s: %v", conn.UserType, conn.UserID, err)
			atomic.StoreInt32(&conn.Closed, 1)
			return
		}
	}
}

// GetConnectionCount returns the total number of active WebSocket connections
func (wm *WebSocketManager) GetConnectionCount() int {
	return int(atomic.LoadInt64(&wm.connectionCount))
}

// GetConnectionsByType returns a slice of connections for a specific user type
func (wm *WebSocketManager) GetConnectionsByType(userType string) []*WebSocketConnection {
	var filtered []*WebSocketConnection
	wm.connections.Range(func(key, value interface{}) bool {
		conn := value.(*WebSocketConnection)
		if conn.UserType == userType && atomic.LoadInt32(&conn.Closed) == 0 {
			filtered = append(filtered, conn)
		}
		return true // Continue iteration
	})
	return filtered
}

// GetConnection returns a specific connection
func (wm *WebSocketManager) GetConnection(userID, userType string) *WebSocketConnection {
	connectionID := userType + ":" + userID
	if connInterface, exists := wm.connections.Load(connectionID); exists {
		return connInterface.(*WebSocketConnection)
	}
	return nil
}

// IsConnected checks if a user is connected
func (wm *WebSocketManager) IsConnected(userID, userType string) bool {
	conn := wm.GetConnection(userID, userType)
	if conn == nil {
		return false
	}

	return atomic.LoadInt32(&conn.Closed) == 0
}

// WebSocket configuration constants
const (
	WebSocketPingInterval   = 30 * time.Second
	WebSocketWriteTimeout   = 10 * time.Second
	WebSocketReadTimeout    = 10 * time.Second
	WebSocketPongTimeout    = 60 * time.Second
	WebSocketMaxMessageSize = 1024
)

// WebSocket upgrader configuration
var WebSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Configure based on your security requirements
}

// CreateWebSocketMessage creates a WebSocket message with current timestamp
func CreateWebSocketMessage(messageType string, data map[string]interface{}) WebSocketMessage {
	return WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// CreateWebSocketErrorMessage creates a WebSocket error message
func CreateWebSocketErrorMessage(messageType string, errorMsg string, data map[string]interface{}) WebSocketMessage {
	if data == nil {
		data = make(map[string]interface{})
	}

	return WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
		Error:     errorMsg,
	}
}

// CreateConnectionEstablishedMessage creates a connection established message
func CreateConnectionEstablishedMessage(userID, userType string, channel string) WebSocketMessage {
	return CreateWebSocketMessage("connection_established", map[string]interface{}{
		"userId":   userID,
		"userType": userType,
		"channel":  channel,
	})
}

// CreatePingMessage creates a ping message
func CreatePingMessage() WebSocketMessage {
	return CreateWebSocketMessage("ping", map[string]interface{}{
		"timestamp": time.Now().Unix(),
	})
}

// CreatePongMessage creates a pong message
func CreatePongMessage() WebSocketMessage {
	return CreateWebSocketMessage("pong", map[string]interface{}{
		"timestamp": time.Now().Unix(),
	})
}

// CreateSystemMessage creates a system message
func CreateSystemMessage(message string) WebSocketMessage {
	return CreateWebSocketMessage("system_message", map[string]interface{}{
		"message": message,
	})
}

// WebSocket message type constants
const (
	MessageTypeConnectionEstablished = "connection_established"
	MessageTypePing                  = "ping"
	MessageTypePong                  = "pong"
	MessageTypeSystemMessage         = "system_message"
	MessageTypeError                 = "error"
	MessageTypeBiddingStarted        = "bidding_started"
	MessageTypeBidReceived           = "bid_received"
	MessageTypeBiddingEnded          = "bidding_ended"
	MessageTypeDriverAssigned        = "driver_assigned"
	MessageTypeTimerUpdate           = "timer_update"
	MessageTypeTripNotification      = "trip_notification"
	MessageTypeTripStatusUpdate      = "trip_status_update"
	MessageTypeDriverLocation        = "driver_location_update"
	MessageTypeNoDriverFound         = "no_driver_found"
)

// WebSocket user type constants
const (
	UserTypeDriver = "driver"
	UserTypeRider  = "rider"
	UserTypeAdmin  = "admin"
)

// WebSocket statistics
type WebSocketStats struct {
	TotalConnections  int    `json:"totalConnections"`
	DriverConnections int    `json:"driverConnections"`
	RiderConnections  int    `json:"riderConnections"`
	AdminConnections  int    `json:"adminConnections"`
	Timestamp         string `json:"timestamp"`
}

// GetWebSocketStats returns WebSocket connection statistics
func GetWebSocketStats(manager IWebSocketManager) WebSocketStats {
	driverConns := manager.GetConnectionsByType(UserTypeDriver)
	riderConns := manager.GetConnectionsByType(UserTypeRider)
	adminConns := manager.GetConnectionsByType(UserTypeAdmin)

	return WebSocketStats{
		TotalConnections:  manager.GetConnectionCount(),
		DriverConnections: len(driverConns),
		RiderConnections:  len(riderConns),
		AdminConnections:  len(adminConns),
		Timestamp:         time.Now().Format(time.RFC3339),
	}
}

// WebSocket connection health check
type ConnectionHealth struct {
	UserID     string    `json:"userId"`
	UserType   string    `json:"userType"`
	LastPing   time.Time `json:"lastPing"`
	IsHealthy  bool      `json:"isHealthy"`
	Connection string    `json:"connection"`
}

// GetConnectionHealth returns health information for a connection
func GetConnectionHealth(manager IWebSocketManager, userID, userType string) ConnectionHealth {
	conn := manager.GetConnection(userID, userType)

	if conn == nil {
		return ConnectionHealth{
			UserID:     userID,
			UserType:   userType,
			IsHealthy:  false,
			Connection: "disconnected",
		}
	}

	return ConnectionHealth{
		UserID:     userID,
		UserType:   userType,
		LastPing:   conn.LastPing,
		IsHealthy:  atomic.LoadInt32(&conn.Closed) == 0,
		Connection: "connected",
	}
}

// BroadcastToMultipleUsers sends a message to multiple users
func BroadcastToMultipleUsers(manager IWebSocketManager, userType string, userIDs []string, message WebSocketMessage) {
	for _, userID := range userIDs {
		manager.SendMessage(userID, userType, message)
	}
}

// BroadcastToAllUsers sends a message to all connected users
func BroadcastToAllUsers(manager IWebSocketManager, message WebSocketMessage) {
	manager.BroadcastToType(UserTypeDriver, message)
	manager.BroadcastToType(UserTypeRider, message)
	manager.BroadcastToType(UserTypeAdmin, message)
}
