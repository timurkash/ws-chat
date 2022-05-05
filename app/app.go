package app

import (
	"github.com/timurkash/ws-chat/app/routers"
	"github.com/timurkash/ws-chat/hub"
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
