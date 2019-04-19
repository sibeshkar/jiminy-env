package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var i = 0

type Message struct {
	Name string
	Body string
	Time int64
}

func main() {
	http.HandleFunc("/", echo)
	http.ListenAndServe(":8080", nil)
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		// Read message from browser
		//_, msg, err := conn.ReadMessage()
		m := Message{}
		err := conn.ReadJSON(&m)
		if err != nil {
			return
		}
		// Print the message to the console
		//fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), m.Body)
		// Write message back to browser
		// if err = conn.WriteMessage(msgType, msg); err != nil {
		// 	return
		// }
	}
}
