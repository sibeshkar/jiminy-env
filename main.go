package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/shared"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var env shared.Env
var envStatus EnvStatus

var agent_conn AgentConn

type AgentConn struct {
	ws  *websocket.Conn
	fps int64
}

type Headers struct {
	Sent_at         int64  `json:"sent_at"`
	MessageId       string `json:"message_id"`
	ParentMessageId string `json:"parent_message_id"`
	EpisodeId       string `json:"episode_id"`
}

type Body struct {
	Env_id string  `json:"env_id"`
	Reward float32 `json:"reward"`
	Done   bool    `json:"done"`
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

	fmt.Println("Connected:")
	agent_conn := AgentConn{
		ws:  conn,
		fps: 600,
	}

	go agent_conn.OnMessage()
	go agent_conn.Environment()
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

func (c *AgentConn) Environment() error {
	for {
		c.SendReward(GetReward())
		time.Sleep(time.Duration(1/c.fps) * time.Second)
	}

}

func (c *AgentConn) SendReward(m *Message) error {
	err := c.ws.WriteJSON(&m)
	if err != nil {
		log.Println("write:", err)
	}
	return err
}

//Random function to generate get reward from environment
func GetReward() *Message {
	//env.GetReward()string
	reward, _ := env.GetReward()

	method := "v0.env.reward"

	headers := Headers{
		Sent_at: time.Now().Unix(),
	}

	body := Body{
		Env_id: "wob.mini.TicTacToe",
		Reward: reward,
	}

	m := Message{
		Method:  method,
		Headers: headers,
		Body:    body,
	}

	return &m
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
