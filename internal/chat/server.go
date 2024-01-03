package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type Server struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
			log.Println("Novo cliente registrado")

		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
				log.Println("Cliente removido")
			}

		case message := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao criar conexão WebSocket:", err)
		return
	}

	client := NewClient(conn, s)
	s.register <- client

	go client.Read()
	go client.Write()
}
