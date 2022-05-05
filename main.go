package main

import (
	"flag"
	"github.com/timurkash/ws-chat/hub"
	"github.com/timurkash/ws-chat/routers"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	hub.Init()
	http.HandleFunc("/", routers.ServeHome)
	http.HandleFunc("/ws", routers.ServeWs)
	log.Println("listen & serve on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
