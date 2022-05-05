package app

import (
	"github.com/timurkash/ws-chat/hub"
	"github.com/timurkash/ws-chat/routers"
	"log"
	"net/http"
)

func RunApp(addr *string) error {
	hub.Init()
	http.HandleFunc("/", routers.ServeHome)
	http.HandleFunc("/ws", routers.ServeWs)
	log.Println("listen & serve on", *addr)
	return http.ListenAndServe(*addr, nil)
}
