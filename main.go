package main

import (
	"flag"
	"github.com/timurkash/ws-chat/routers"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func init() {
	flag.Parse()
}

func main() {
	http.HandleFunc("/", routers.ServeHome)
	http.HandleFunc("/ws", routers.ServeWs)
	log.Println("listen & serve on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
