package hub

import (
	"github.com/timurkash/ws-chat/wsclient"
)

var Hub = NewHub()

type Struct struct {
	Register   chan *wsclient.Client
	Unregister chan *wsclient.Client
	Broadcast  chan []byte
	Clients    map[*wsclient.Client]struct{}
}

func NewHub() *Struct {
	hub := &Struct{
		Register:   make(chan *wsclient.Client),
		Unregister: make(chan *wsclient.Client),
		Broadcast:  make(chan []byte),
		Clients:    make(map[*wsclient.Client]struct{}),
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
