package chatBoard

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

type ChatRoom struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}
type UserMessage struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}
type BroadcastMessage struct {
	UserMessage
	UserIp    string `json:"userIp"`
	Timestamp string `json:"timestamp"`
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (room *ChatRoom) run() {
	for {
		select {
		case client := <-room.register:
			room.clients[client] = true
		case client := <-room.unregister:
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				client.Close()
			}
		case message := <-room.broadcast:
			for client := range room.clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					fmt.Println("Error broadcasting message:", err)
					client.Close()
					delete(room.clients, client)
				}
			}
		}
	}
}

func ChatBoard_routerGroup_init(router *gin.Engine) {
	room := NewChatRoom()
	go room.run()
	chatBoardApiGroup := router.Group("/api/chatBoard")
	chatBoardApiGroup.GET("/chat", func(c *gin.Context) {
		conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("Error upgrading to websocket:", err)
			return
		}
		defer conn.Close()
		room.register <- conn
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("Error: %v\n", err)
				}
				break
			}
			var receivedMessage UserMessage
			if err := json.Unmarshal(msg, &receivedMessage); err != nil {
				fmt.Println("Error unmarshaling message:", err)
				continue
			}
			bm := BroadcastMessage{
				UserMessage: receivedMessage,
				UserIp:      conn.RemoteAddr().String(),
				Timestamp:   strconv.FormatInt(time.Now().Unix(), 10)}
			jsonMsg, err := json.Marshal(bm)
			if err != nil {
				fmt.Println("Error marshaling message:", err)
				continue
			}
			room.broadcast <- jsonMsg
		}
		room.register <- conn
	})
}
