package websocket

import (
	"encoding/json"
	"fmt"
	"instant-messenger-backend/models"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var roomID string
var server *Server

type ChatRoom struct {
	participants map[*websocket.Conn]struct{}
	roomID       string
	roomMu       sync.Mutex
}

type Server struct {
	rooms    map[string]*ChatRoom
	roomsMu  sync.Mutex
	upgrader websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		rooms:    make(map[string]*ChatRoom),
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

// GetServerInstance returns the current server instance (if initialized)
func GetServerInstance() *Server {
	return server
}

func (s *Server) createRoom(roomID string) {
	s.roomsMu.Lock()
	defer s.roomsMu.Unlock()

	s.rooms[roomID] = &ChatRoom{
		participants: make(map[*websocket.Conn]struct{}),
		roomID:       roomID,
	}
}

func (s *Server) joinRoom(roomID string, conn *websocket.Conn) error {
	s.roomsMu.Lock()
	defer s.roomsMu.Unlock()

	room, ok := s.rooms[roomID]
	if !ok {
		return fmt.Errorf("room %s not found", roomID)
	}

	room.roomMu.Lock()
	defer room.roomMu.Unlock()

	if len(room.participants) >= 2 {
		return fmt.Errorf("room %s is full", roomID)
	}

	room.participants[conn] = struct{}{}

	return nil
}

func (s *Server) broadcastRoom(roomID string, msg []byte) error {
	fmt.Println("Log 1 - Broadcasting message to room", roomID)
	s.roomsMu.Lock()
	defer s.roomsMu.Unlock()

	room, ok := s.rooms[roomID]
	if !ok {
		return fmt.Errorf("room %s not found", roomID)
	}

	fmt.Println("Log 2 - Broadcasting message to room", roomID)

	fmt.Println("Broadcasting Message", string(msg))

	room.roomMu.Lock()
	defer room.roomMu.Unlock()

	for participant := range room.participants {
		if err := participant.WriteMessage(websocket.TextMessage, msg); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	roomID = r.URL.Query().Get("room")
	if roomID == "" {
		roomID = "default"
	}

	if _, ok := s.rooms[roomID]; !ok {
		s.createRoom(roomID)
	}

	if err := s.joinRoom(roomID, conn); err != nil {
		log.Println(err)
		return
	}

	fmt.Println("a user connected")

	// for {
	// 	_, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println(err)
	// 		s.roomsMu.Lock()
	// 		delete(s.rooms[roomID].participants, conn)
	// 		s.roomsMu.Unlock()
	// 		return
	// 	}

	// 	if err := s.broadcastRoom(roomID, msg); err != nil {
	// 		log.Println("Error broadcasting message:", err)
	// 	}
	// }
}

// SendChatToClient broadcasts the message to the specified room
// This function can be called without needing a server instance
func SendChatToClient(messageJSON []byte) error {
    // Decode JSON message
    var message models.Message
    err := json.Unmarshal(messageJSON, &message)
    if err != nil {
        return fmt.Errorf("error unmarshalling message: %w", err)
    }

    if server == nil {
        return fmt.Errorf("server is not initialized")
    }

    // Convert message.Content to a byte slice
    messageContent := []byte(message.Content)

    // Broadcast the message using websocket.TextMessage
    server.broadcastRoom(roomID, messageContent)

    return nil
}

func InitGorillaWebsocket() {
	server = NewServer()

	http.HandleFunc("/ws", server.handleWebSocket)

	fmt.Println("Gorilla Websocket Server is listening on port 8181")
	log.Fatal(http.ListenAndServe(":8181", nil))
}