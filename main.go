package main

import (
	"flag"
	"github.com/timurkash/ws-chat/app"
	"log"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	if err := app.RunApp(addr); err != nil {
		log.Fatal("RunApp: ", err)
	}
}
