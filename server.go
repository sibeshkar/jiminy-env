package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Headers struct {
	Sent_at           string `json:"sent_at"`
	Message_id        string `json:"message_id"`
	Parent_message_id string `json:"parent_message_id"`
	Episode_id        string `json:"episode_id"`
}

type Body struct {
	Env_id string `json:"env_id"`
}

type Message struct {
	Method  string  `json:"method"`
	Headers Headers `json:"headers"`
	Body    Body    `json:"body"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		//_, msg, err := conn.ReadMessage()
		m := Message{}
		err := conn.ReadJSON(&m)
		if err != nil {
			return
		}
		if m.Method == "v0.env.launch" {
			go Launch(&m)
		} else if m.Method == "v0.env.reset" {
			go Reset(&m)
		} else if m.Method == "v0.env.close" {
			go Close(&m)
		}

	}
}

func Launch(m *Message) {
	fmt.Printf("launch message received: %s\n", m.Body.Env_id)
}

func Reset(m *Message) {
	fmt.Printf("reset message received: %s\n", m.Body.Env_id)
}

func Close(m *Message) {
	fmt.Printf("close message received: %s\n", m.Body.Env_id)
}
