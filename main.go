package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/transport/ws"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func subHandler(sock mangos.Socket) {
	count := 0
	for {
		msg := fmt.Sprintf("PUB #%d %s", count, time.Now().String())
		if err := sock.Send([]byte(msg)); err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
		count++
	}
}

func main() {
	fmt.Println("awesome")

	sock, err := pub.NewSocket()
	if err != nil {
		panic(err)
	}

	listener, err := sock.NewListener("ws://localhost:3000/sub", nil)
	if err != nil {
		panic(err)
	}

	h, err := listener.GetOption(ws.OptionWebSocketHandler)
	if err != nil {
		panic(err)
	}

	http.Handle("/sub", h.(http.Handler))
	listener.Listen()

	go subHandler(sock)

	http.HandleFunc("/", handler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
