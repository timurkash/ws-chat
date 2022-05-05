package routers

import (
	"github.com/gorilla/websocket"
	"github.com/timurkash/ws-chat/hub"
	"github.com/timurkash/ws-chat/hub/ws"
	"log"
	"net/http"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 11024

	sendBufferSize = 256
)

func ServeWs(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &ws.Client{
		Conn: conn,
		Send: make(chan []byte, sendBufferSize),
		Register: func(client *ws.Client) {
			hub.Hub.Unregister <- client
		},
		Broadcast: func(message []byte) {
			hub.Hub.Broadcast <- message
		},
	}
	hub.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "app/routers/home.html")
}
