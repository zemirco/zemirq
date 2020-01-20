package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/protocol/rep"
	"nanomsg.org/go/mangos/v2/transport/ws"
)

func reqHandler(sock mangos.Socket) {
	fmt.Println("repHandler running.")
	var err error
	var msg []byte

	for {
		fmt.Println("waiting for REQ")
		// Could also use sock.RecvMsg to get header
		msg, _ = sock.Recv()
		fmt.Println("parsing REQ")
		if string(msg) == "DATE" { // no need to terminate
			fmt.Println("RECEIVED DATE REQUEST")
			d := date()
			fmt.Printf("SENDING DATE %s\n", d)
			err = sock.Send([]byte(d))
			if err != nil {
				die("can't send reply: %s", err.Error())
			}
		} else if string(msg) == "GREET" { // no need to terminate
			fmt.Println("RECEIVED GREET REQUEST")
			d := "Hello, requester!"
			fmt.Printf("SENDING GREET %s\n", d)
			err = sock.Send([]byte(d))
			if err != nil {
				die("can't send reply: %s", err.Error())
			}
		}
	}
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func subHandler(sock mangos.Socket) {
	fmt.Println("subHandler running.")
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
	fmt.Println("awesome, let's go!")

	sock, err := pub.NewSocket()
	if err != nil {
		panic(err)
	}

	sublistener, err := sock.NewListener("ws://127.0.0.1:8081/sub", nil)
	if err != nil {
		panic(err)
	}

	subh, err := sublistener.GetOption(ws.OptionWebSocketHandler)
	if err != nil {
		panic(err)
	}

	http.Handle("/sub", subh.(http.Handler))
	sublistener.Listen()
	go subHandler(sock)

	repsock, err := rep.NewSocket()
	if err != nil {
		panic(err)
	}

	replistener, err := repsock.NewListener("ws://127.0.0.1:8081/req", nil)
	if err != nil {
		panic(err)
	}

	reph, err := replistener.GetOption(ws.OptionWebSocketHandler)
	if err != nil {
		panic(err)
	}

	http.Handle("/req", reph.(http.Handler))
	replistener.Listen()
	go reqHandler(repsock)

	http.Handle("/", http.FileServer(http.Dir("public")))

	e := http.ListenAndServe(":8081", nil)
	die("Http server died: %v", e)
}
