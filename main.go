package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/shared"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var addr = flag.String("addr", "localhost:15900", "http service address")
var env shared.Env

var agent_conn AgentConn

type AgentConn struct {
	ws       *websocket.Conn
	envState *EnvState
}

type Headers struct {
	Sent_at         int64  `json:"sent_at"`
	MessageId       string `json:"message_id"`
	ParentMessageId string `json:"parent_message_id"`
	EpisodeId       string `json:"episode_id"`
}

type Body struct {
	EnvId  string  `json:"env_id"`
	Reward float32 `json:"reward"`
	Done   bool    `json:"done"`
}

type Message struct {
	Method  string  `json:"method"`
	Headers Headers `json:"headers"`
	Body    Body    `json:"body"`
}

func main() {
	env = pluginRPC()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func pluginRPC() shared.Env {
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

	agent_conn, err := NewAgentConn(w, r)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Connected to ", agent_conn.ws.RemoteAddr())

	go agent_conn.OnMessage()
	go agent_conn.RewardController()
	//go agent_conn.EnvController()
}

func NewAgentConn(w http.ResponseWriter, r *http.Request) (AgentConn, error) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	envState, err := NewEnvState("wob.mini.TicTacToe", "resetting", 1, 60)
	agent_conn := AgentConn{
		ws:       conn,
		envState: envState,
	}

	if err != nil {
		log.Println(err)
	}

	return agent_conn, err
}

func (c *AgentConn) OnMessage() error {
	for {
		//_, msg, err := conn.ReadMessage()
		m := Message{}
		err := c.ws.ReadJSON(&m)
		if err != nil {
			return err
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
	fmt.Printf("launch message received: %s\n", m.Body.EnvId)
	result, err := env.Launch(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during launch: \n", err)
	}
	fmt.Println("from binary received: ", result)
}

func Reset(m *Message) {
	fmt.Println("reset message received: %s\n", m.Body.EnvId)
	result, err := env.Reset(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during reset: \n", err)
	}
	fmt.Println("from binary received: ", result)
}

func Close(m *Message) {
	fmt.Println("close message received: %s\n", m.Body.EnvId)
	result, err := env.Close(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during close: \n", err)
	}
	fmt.Println("from binary received: ", result)
}
