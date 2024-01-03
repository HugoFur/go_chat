package main

import (
	"go_chat/internal/chat"
	"log"
	"net/http"
)

func main() {
	server := chat.NewServer()
	go server.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.HandleWebSocket(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
