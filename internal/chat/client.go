package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client representa um cliente conectado ao servidor de chat
type Client struct {
	socket *websocket.Conn // Conexão WebSocket do cliente
	send   chan Message    // Canal para enviar mensagens
	server *Server         // Referência ao servidor do chat
}

// NewClient cria um novo cliente
func NewClient(socket *websocket.Conn, server *Server) *Client {
	return &Client{
		socket: socket,
		send:   make(chan Message),
		server: server,
	}
}

func (c *Client) Read() {
	defer func() {
		c.server.unregister <- c
		c.socket.Close()
	}()

	for {
		var msg Message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			log.Printf("Erro na leitura: %v", err)
			break
		}
		c.server.broadcast <- msg
	}
}

func (c *Client) Write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// O servidor fechou o canal
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.socket.WriteJSON(message)
			if err != nil {
				log.Printf("Erro na escrita: %v", err)
				return
			}
		}
	}
}
