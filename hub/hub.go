package hub

import (
	"github.com/timurkash/ws-chat/ws"
)

var Hub *Struct

type Struct struct {
	Register   chan *ws.Client
	Unregister chan *ws.Client
	Broadcast  chan []byte
	Clients    map[*ws.Client]struct{}
}

func Init() {
	Hub = newHub()
}

func newHub() *Struct {
	hub := &Struct{
		Register:   make(chan *ws.Client),
		Unregister: make(chan *ws.Client),
		Broadcast:  make(chan []byte),
		Clients:    make(map[*ws.Client]struct{}),
	}
	go hub.run()
	return hub
}

func (h *Struct) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = struct{}{}
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
