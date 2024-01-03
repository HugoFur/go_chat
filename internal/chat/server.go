package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender  string `json:"sender"`  // Nome do remetente
	Content string `json:"content"` // Conteúdo da mensagem
}

type Server struct {
	clients    map[*Client]bool // conjunto de clientes conectados
	register   chan *Client     // canal para registrar novos clientes
	unregister chan *Client     // canal para remover clientes
	broadcast  chan Message     // canal para enviar mensagens para todos os clientes
}

func newClient(socket *websocket.Conn, server *Server) *Client {
	return &Client{
		socket: socket,
		send:   make(chan Message),
		server: server,
	}
}

// NewServer cria um novo servidor de chat
func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

// Run inicia o servidor e lida com canais de registro, cancelamento de registro e broadcast
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
			// Enviar mensagem para todos os clientes conectados
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

// HandleWebSocket lida com a conexão WebSocket
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao criar conexão WebSocket:", err)
		return
	}

	client := NewClient(conn, s)
	s.register <- client

	// Iniciar rotinas de leitura e escrita do cliente
	go client.Read()
	go client.Write()
}
