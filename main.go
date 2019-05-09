package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/shared"
	"github.com/urfave/cli"

	"github.com/gorilla/websocket"
)

var EnvPlugin shared.PluginConfig

//jiminy run sibeshkar/wob-v0/TicTacToe
//jiminy build .
//jiminy pull sibeshkar/wob-v0/TicTacToe
//jiminy push sibeshkar/wob-v0/TicTacToe

var upgrader = websocket.Upgrader{}

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

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"a"},
			Usage:   "Run an installed environment. e.g. sibeshkar/wob-v0",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				Run(c.Args().First())
				return nil
			},
		},
		{
			Name:    "install",
			Aliases: []string{"c"},
			Usage:   "Install env plugin from directory or zip file.",
			Action: func(c *cli.Context) error {
				fmt.Println("the directory is ", c.Args().First())
				var format []string = strings.Split(c.Args().First(), ".")
				if format[len(format)-1] == "zip" {
					shared.InstallFromArchive(c.Args().First())
				} else {
					shared.Install(c.Args().First())
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func Run(pluginLink string) {

	EnvPlugin = shared.CreatePluginConfig(pluginLink)
	env = pluginRPC(&EnvPlugin)
	go env.Init(pluginLink)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":15900", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {

	agent_conn, err := NewAgentConn(w, r)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Connected to ", agent_conn.ws.RemoteAddr())

	go agent_conn.OnMessage()

	for {

		switch agent_conn.envState.EnvStatus {
		case "launching":
			fmt.Println("Env is still launching")
		case "resetting":
			fmt.Println("Env is resetting to task")
		case "running":
			reply := GetReward()
			agent_conn.SendToAgent(ConstructRewardMessage(reply))
			fmt.Println("Environment is running")
		}

		time.Sleep(time.Duration(1000/agent_conn.envState.Fps) * time.Millisecond)
	}

}

func pluginRPC(pluginObj *shared.PluginConfig) shared.Env {
	log.SetOutput(ioutil.Discard)

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", pluginObj.BinaryFile),
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

func NewAgentConn(w http.ResponseWriter, r *http.Request) (AgentConn, error) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	envState, err := NewEnvState("wob.mini.TicTacToe", "resetting", 1, 20)
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
			c.Launch(&m)
		} else if m.Method == "v0.env.reset" {
			c.Reset(&m)
		} else if m.Method == "v0.env.close" {
			c.Close(&m)
		}

	}

}

func (c *AgentConn) SendToAgent(m Message) error {
	err := c.ws.WriteJSON(&m)
	if err != nil {
		log.Println("write:", err)
	}
	return err
}

func (c *AgentConn) Launch(m *Message) {
	fmt.Printf("launch message received: %s\n", m.Body.EnvId)
	c.envState.SetEnvStatus("launching")

	result, err := env.Launch(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during launch: \n", err)
	}
	c.envState.SetEnvId(m.Body.EnvId)
	fmt.Println("from binary received: ", result)
}

func (c *AgentConn) Reset(m *Message) {
	fmt.Println("reset message received: %s\n", m.Body.EnvId)
	c.envState.SetEnvStatus("resetting")

	result, err := env.Reset(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during reset: \n", err)
	}
	c.envState.SetEnvStatus("running")
	c.envState.SetEnvId(m.Body.EnvId)
	fmt.Println("from binary received: ", result)
}

func (c *AgentConn) Close(m *Message) {
	fmt.Println("close message received: %s\n", m.Body.EnvId)
	result, err := env.Close(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during close: \n", err)
	}
	fmt.Println("from binary received: ", result)
}

// go func() {

// 	for {
// 		select {
// 		case status := <-statusChan:
// 			switch status {
// 			case "launching":
// 				fmt.Println("Runtime is launching")
// 			case "resetting":
// 				fmt.Println("Env is resetting")
// 			case "running":
// 				fmt.Println("Env is running")
// 			}
// 		case trigger := <-triggerChan:
// 			switch trigger {
// 			case pluginObj *shared.PluginConfig"launch":
// 				Launch(&m)
// 			case "reset":
// 				Reset(&m)
// 			}
// 		}
// 	}
// }()

// statusChan := make(chan string)
// triggerChan := make(chan string)
// rewardChan := make(chan Body)
// time.Sleep(5 * time.Second)

// go RewardController(statusChan, triggerChan, rewardChan)
// // // //go EnvController(status)
// // // //go agent_conn.RewardController()
// ticker := time.NewTicker(time.Duration(1000/agent_conn.envState.Fps) * time.Millisecond)

// go func() {
// 	i := 0
// 	for {
// 		select {
// 		case <-ticker.C:
// 			fmt.Println("Sending reward ", i)
// 			agent_conn.Send(ConstructRewardMessage(GetReward()))
// 			i++
// 		case state := <-statusChan:
// 			fmt.Println("Status is", state)
// 			agent_conn.envState.SetEnvStatus(state)
// 			// default:
// 			// 	fmt.Println("Default option")
// 		}

// 	}
// }()

//go agent_conn.EnvController()
//}
