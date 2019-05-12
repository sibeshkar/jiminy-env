package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
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
	Sent_at         int64   `json:"sent_at"`
	MessageId       string  `json:"message_id"`
	ParentMessageId string  `json:"parent_message_id"`
	EpisodeId       float32 `json:"episode_id"`
}

type Body struct {
	EnvId     string  `json:"env_id"`
	EnvStatus string  `json:"env_status"`
	Fps       float32 `json:"fps"`
	Reward    float32 `json:"reward"`
	Done      bool    `json:"done"`
	Obs       string  `json:"observation"`
	ObsType   string  `json:"observation_type"`
	Info      string  `json:"info"`
	InfoType  string  `json:"info_type"`
	Message   string  `json:"message"`
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

	var lastdone bool

	for {

		switch agent_conn.envState.EnvStatus {
		case "launching":
			fmt.Println("Env is still launching")
			time.Sleep(100 * time.Microsecond)
		case "resetting":
			fmt.Println("Env is resetting to task")
			time.Sleep(100 * time.Microsecond)
		case "running":
			reward, done, _ := env.GetReward()

			agent_conn.SendEnvReward(reward, done)
			agent_conn.SendEnvObservation()
			if done && lastdone != true {
				go agent_conn.Reset()
				fmt.Println("Environment is resetting to task again")
				agent_conn.envState.SetEpisodeId(agent_conn.envState.GetEpisodeId() + 1)
			}
			fmt.Println("Environment is running")
			lastdone = done
		}

		time.Sleep(time.Duration(1000/agent_conn.envState.Fps) * time.Millisecond)
	}

}

func pluginRPC(pluginObj *shared.PluginConfig) shared.Env {
	//log.SetOutput(ioutil.Discard)

	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		//Level:  hclog.Trace,
		Name: "plugin",
	})

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", pluginObj.BinaryFile),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Logger: logger,
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

	envState, err := NewEnvState("", "launching", 1, 20)
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
			c.InitLaunch(&m)
		} else if m.Method == "v0.env.reset" {
			c.InitReset(&m)
		} else if m.Method == "v0.env.close" {
			c.Close(&m)
		}

	}

}

func (c *AgentConn) InitLaunch(m *Message) {
	fmt.Printf("launch message received: %s\n", m.Body.EnvId)
	c.envState.SetEnvStatus("launching")

	result, err := env.Launch(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during launch: \n", err)
	}
	c.envState.SetEnvId(m.Body.EnvId)
	fmt.Println("from binary received: ", result)
}

func (c *AgentConn) SendLaunchReply(parent_message_id string, err error) {

	//include this inside everytime
	// func (c *AgentConn) SendLaunchError() {

	// }

}

func (c *AgentConn) InitReset(m *Message) {
	fmt.Println("reset message received: %s\n", m.Body.EnvId)
	c.envState.SetEnvStatus("resetting")

	result, err := env.Reset(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during reset: \n", err)
	} else {
		c.envState.SetEnvStatus("running")
		c.envState.SetEnvId(m.Body.EnvId)
		fmt.Println("from binary received: ", result)

	}

	c.SendResetReply(m.Headers.MessageId, err)
}

func (c *AgentConn) SendResetReply(parent_message_id string, err error) {

	if err != nil {
		method := "v0.reply.error"

		headers := Headers{
			Sent_at:         time.Now().Unix(),
			EpisodeId:       c.envState.GetEpisodeId(),
			ParentMessageId: parent_message_id,
		}

		body := Body{
			Message: err.Error(),
		}

		m := Message{
			Method:  method,
			Headers: headers,
			Body:    body,
		}

		c.SendMessage(m)

	} else {

		method := "v0.reply.env.reset"

		headers := Headers{
			Sent_at:         time.Now().Unix(),
			EpisodeId:       c.envState.GetEpisodeId(),
			ParentMessageId: parent_message_id,
		}

		body := Body{}

		m := Message{
			Method:  method,
			Headers: headers,
			Body:    body,
		}

		c.SendMessage(m)

	}

}

func (c *AgentConn) Reset() {

	c.envState.SetEnvStatus("resetting")
	c.SendEnvDescribe()

	_, err := env.Reset(c.envState.GetEnvId())
	if err != nil {
		fmt.Println("error during reset: \n", err)
	}

	c.envState.SetEnvStatus("running")

}

func (c *AgentConn) Close(m *Message) {
	fmt.Println("close message received: %s\n", m.Body.EnvId)
	result, err := env.Close(m.Body.EnvId)
	if err != nil {
		fmt.Println("error during close: \n", err)
	}
	fmt.Println("from binary received: ", result)
}

//Send a message to connected Agent
func (c *AgentConn) SendMessage(m Message) error {
	err := c.ws.WriteJSON(&m)
	if err != nil {
		log.Println("write:", err)
	}
	return err
}

func (c *AgentConn) SendEnvReward(reward float32, done bool) {

	method := "v0.env.reward"

	headers := Headers{
		Sent_at:   time.Now().Unix(),
		EpisodeId: c.envState.GetEpisodeId(),
	}

	body := Body{
		EnvId:  c.envState.GetEnvId(),
		Reward: reward,
		Done:   done,
	}

	m := Message{
		Method:  method,
		Headers: headers,
		Body:    body,
	}

	c.SendMessage(m)

}

//Required function of AgentConn
//
// func (c *AgentConn) SendEnvText(text string) {

// }
//Protobuf method GetEnvObservation (sent once every 1/fps)
func (c *AgentConn) SendEnvObservation() {

	t, obs, err := env.GetEnvObservation(c.envState.EnvId)
	if err != nil {
		fmt.Println(err)
	}
	var observation string
	if t == "image" {
		observation = base64.StdEncoding.EncodeToString(obs)
	} else {
		observation = string(obs)
	}

	method := "v0.env.observation"

	headers := Headers{
		Sent_at:   time.Now().Unix(),
		EpisodeId: c.envState.GetEpisodeId(),
	}

	body := Body{
		Obs:     observation,
		ObsType: t,
	}

	m := Message{
		Method:  method,
		Headers: headers,
		Body:    body,
	}

	c.SendMessage(m)

}

func (c *AgentConn) SendEnvInfo() {

	t, info, err := env.GetEnvInfo(c.envState.EnvId)
	if err != nil {
		fmt.Println(err)
	}

	method := "v0.env.info"

	headers := Headers{
		Sent_at:   time.Now().Unix(),
		EpisodeId: c.envState.GetEpisodeId(),
	}

	body := Body{
		Info:     base64.StdEncoding.EncodeToString(info),
		InfoType: t,
	}

	m := Message{
		Method:  method,
		Headers: headers,
		Body:    body,
	}

	c.SendMessage(m)

}

func (c *AgentConn) SendEnvDescribe() {

	method := "v0.env.describe"

	headers := Headers{
		Sent_at:   time.Now().Unix(),
		EpisodeId: c.envState.GetEpisodeId(),
	}

	body := Body{
		EnvId:     c.envState.GetEnvId(),
		EnvStatus: c.envState.GetEnvStatus(),
		Fps:       c.envState.GetFps(),
	}

	m := Message{
		Method:  method,
		Headers: headers,
		Body:    body,
	}

	c.SendMessage(m)

}

func DummyObs() {
	t, obs, err := env.GetEnvObservation("test")
	if err != nil {
		fmt.Println(err)
	}

	// img, _, err := image.Decode(bytes.NewReader(obs))
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// //save the imgByte to file
	// out, err := os.Create("./QRImg.png")

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// err = png.Encode(out, img)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	//obsString := base64.StdEncoding.EncodeToString(obs)
	var observation string
	if t == "image" {
		observation = base64.StdEncoding.EncodeToString(obs)
	} else {
		observation = string(obs)
	}
	fmt.Println("The type is ", t)
	fmt.Println("The obs is ", observation)
}

func DummyInfo() {

	t, info, err := env.GetEnvInfo("test")
	if err != nil {
		fmt.Println(err)
	}

	infoString := base64.StdEncoding.EncodeToString(info)
	fmt.Println("The type is ", t)
	fmt.Println("The info is ", infoString)

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
