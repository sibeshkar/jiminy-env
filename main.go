package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/shared"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var env shared.Env

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
	env = startRPC()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":15900", nil)
}

func startRPC() shared.Env {
	log.SetOutput(ioutil.Discard)

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv("ENV_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	// defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("env_grpc")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	return raw.(shared.Env)
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
			Launch(&m)
		} else if m.Method == "v0.env.reset" {
			Reset(&m)
		} else if m.Method == "v0.env.close" {
			Close(&m)
		}

	}
}

func Launch(m *Message) {
	fmt.Printf("launch message received: %s\n", m.Body.Env_id)
	result, err := env.Launch(m.Body.Env_id)
	if err != nil {
		fmt.Printf("error during launch: \n", err)
	}
	fmt.Printf("from binary received: %s\n", result)
}

func Reset(m *Message) {
	fmt.Printf("reset message received: %s\n", m.Body.Env_id)
	result, err := env.Reset(m.Body.Env_id)
	if err != nil {
		fmt.Printf("error during reset: \n", err)
	}
	fmt.Printf("from binary received: %s\n", result)
}

func Close(m *Message) {
	fmt.Printf("close message received: %s\n", m.Body.Env_id)
	result, err := env.Close(m.Body.Env_id)
	if err != nil {
		fmt.Printf("error during close: \n", err)
	}
	fmt.Printf("from binary received: %s\n", result)
}
