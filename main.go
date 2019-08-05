package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/shared"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var EnvPlugin shared.PluginConfig
var defaultFps float32 = 20

//jiminy run sibeshkar/wob-v0/TicTacToe
//jiminy install <folder>
//jiminy zip <folder>
//jiminy pull sibeshkar/wob-v0/TicTacToe
//jiminy push sibeshkar/wob-v0/TicTacToe

var upgrader = websocket.Upgrader{}

var env shared.Env
var client_conf ClientConfig

var agent_conn AgentConn

type ClientConfig struct {
	client *plugin.Client
	init   bool
}

type AgentConn struct {
	wsLock   sync.Mutex
	ws       *websocket.Conn
	envState *EnvState
}

type Headers struct {
	Sent_at         float64 `json:"sent_at"`
	MessageId       int32   `json:"message_id"`
	ParentMessageId int32   `json:"parent_message_id"`
	EpisodeId       int64   `json:"episode_id"`
}

type Body struct {
	EnvId     string  `json:"env_id"`
	EnvStatus string  `json:"env_state"` //TODO : Fix inconsistency later
	Fps       float32 `json:"fps"`
	Reward    float32 `json:"reward"`
	Done      bool    `json:"done"`
	Obs       string  `json:"observation"`
	ObsType   string  `json:"observation_type"`
	Info      string  `json:"info"`
	InfoType  string  `json:"info_type"`
	Message   string  `json:"message"`
	Seed      int64   `json:"seed"`
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
			Usage:   "Run an installed environment. e.g. {empty} or sibeshkar/wob-v0 or sibeshkar/wob-v0/TicTacToe",
			Action: func(c *cli.Context) error {

				if len(c.Args().First()) != 0 {
					s := strings.Split(c.Args().First(), "/")
					if len(s) == 2 {

						RunPlugin(c.Args().First())

					} else if len(s) == 3 {
						RunPluginTask(c.Args().First())
					}

				} else {
					log.Info("Running websocket server...")
					RunEmpty()
				}
				return nil
			},
		},
		{
			Name:    "install",
			Aliases: []string{"c"},
			Usage:   "Install env plugin from directory or zip file.",
			Action: func(c *cli.Context) error {
				log.Info("the directory is ", c.Args().First())
				var format []string = strings.Split(c.Args().First(), ".")
				if format[len(format)-1] == "zip" {
					shared.InstallFromArchive(c.Args().First())
				} else {
					shared.Install(c.Args().First())
				}

				return nil
			},
		},
		{
			Name:    "zip",
			Aliases: []string{"c"},
			Usage:   "Zip plugin folder according to config.json to create zip archive inside",
			Action: func(c *cli.Context) error {
				log.Info("the directory is ", c.Args().First())

				shared.CreateArchive(c.Args().First())

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func RunPlugin(pluginLink string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	log.Info("Running environment: ", pluginLink)

	EnvPlugin = shared.CreatePluginConfig(pluginLink)
	env, client_conf.client = pluginRPC(&EnvPlugin)
	client_conf.init = true
	//TODO: client.Kill() before ending process, otherwise there are zombie plugin processes
	go env.Init(pluginLink)
	//time.Sleep(2 * time.Second)
	env.Launch(pluginLink)
	log.Info("Running websocket server...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":15900", nil))

}

func RunPluginTask(pluginLink string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	s := strings.Split(pluginLink, "/")

	envPluginLink := s[0] + "/" + s[1]
	log.Info("Running environment: ", envPluginLink)

	EnvPlugin = shared.CreatePluginConfig(envPluginLink)
	env, client_conf.client = pluginRPC(&EnvPlugin)
	client_conf.init = true
	//TODO: client.Kill() before ending process, otherwise there are zombie plugin processes
	go env.Init(envPluginLink)
	//time.Sleep(2 * time.Second)

	for {
		_, err := env.Launch(envPluginLink)
		if err != nil {
			log.Errorf("Error while launching %v ", err)
			break
		}
		//time.Sleep(2 * time.Second)

		log.Info("Resetting to task: ", pluginLink)
		_, err = env.Reset(pluginLink)
		if err != nil {
			log.Errorf("Error while resetting %v ", err)
		} else {
			break
		}
		//time.Sleep(1 * time.Second)

	}

	log.Info("Running websocket server...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":15900", nil))

}

func RunEmpty() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":15900", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {

	agent_conn, err := NewAgentConn(w, r)
	if err != nil {
		log.Println(err)
	}
	log.Info("Connected to ", agent_conn.ws.RemoteAddr())

	go agent_conn.OnMessage()

	var lastdone bool = true

	t0 := time.Now()
	var n float32 = 0.0

mainloop:
	for {
		n += 1.0
		t := t0.Add(time.Duration(n*1000/agent_conn.envState.Fps) * time.Millisecond)
		dt := t.Sub(time.Now())
		//log.Infof("n:%v, t: %v, t0: %v, dt: %v", n, t, t0, dt)

		if dt > 0 {
			if dt >= time.Duration(1*time.Second) {
				log.Infof("Sleeping for %v", dt)
			}
			//log.Infof("Sleeping for %v", dt)
			time.Sleep(dt)
		} else {
			log.Infof("Rewarder falling behind %v", dt)

		}
		switch agent_conn.envState.GetEnvStatus() {
		case "launching":
			log.Info("Env is still launching")
			//time.Sleep(20 * time.Microsecond)
		case "resetting":
			log.Info("Env is resetting to task")
			//time.Sleep(20 * time.Microsecond)
		case "running":
			reward, done, _ := env.GetReward()

			//log.Infof("reward:%v, done: %v", reward, done)

			if reward != 0.0 {
				if !done {
					log.Infof("Done is not true")
					done = true
				}
				log.Infof("The reward is %v, the done is %v", reward, done)
				if err := agent_conn.SendEnvReward(reward, done, `{}`); err != nil {
					log.Error(err)
					break mainloop
				}
			}

			// if err := agent_conn.SendEnvObservation(); err != nil {
			// 	log.Error(err)
			// }

			if err := agent_conn.SendEnvObservation(); err != nil {
				log.Error(err)
				break mainloop
			}

			// t, obs, err := env.GetEnvObs("agent_conn.envState.EnvId")

			// log.Infof("The type is %v, the obs is %v, error is %v:", t, obs, err)

			if done != lastdone {
				if done == true || lastdone == false {
					currEpisode := agent_conn.envState.GetEpisodeId()
					agent_conn.envState.SetEpisodeId(currEpisode + 1)
					log.Info("Environment is resetting to task again")
					agent_conn.Reset()

				} else if done == false || lastdone == true {
					currEpisode := agent_conn.envState.GetEpisodeId()
					log.Info("Episode ID is ", currEpisode)
					agent_conn.SendEnvDescribe()
					n = 0.0
					t0 = time.Now()
				}

			}

			lastdone = done
		}

		//time.Sleep(time.Duration(1000/agent_conn.envState.Fps) * time.Millisecond)
	}

}

func pluginRPC(pluginObj *shared.PluginConfig) (shared.Env, *plugin.Client) {
	//log.SetOutput(ioutil.Discard)

	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		//Level:  hclog.Trace, //Uncomment this line to get more detailed plugin Trace errors
		Name: "plugin",
	})

	// We're a host. Start by launching the plugin process.
	cfg := &plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", pluginObj.BinaryFile),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Logger: logger,
	}
	client := plugin.NewClient(cfg)

	// go func() {
	// 	c := make(chan os.Signal, 1)
	// 	signal.Notify(c, os.Interrupt)

	// 	<-c
	// 	err := cfg.Cmd.Process.Kill()
	// 	if err != nil {
	// 		log.Fatalf("Process exit failed: %v", err)
	// 	}
	// }()
	// defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Info("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("env_grpc")
	if err != nil {
		log.Info("Error:", err.Error())
		os.Exit(1)
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	return raw.(shared.Env), client
}

//Create new AgentConnection
func NewAgentConn(w http.ResponseWriter, r *http.Request) (AgentConn, error) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	envState, err := NewEnvState("", "launching", 1, defaultFps)
	agent_conn := AgentConn{
		ws:       conn,
		envState: envState,
	}

	if err != nil {
		log.Println(err)
	}

	return agent_conn, err
}

//Concurrent process to process incoming websocket messages
func (c *AgentConn) OnMessage() error {
	for {
		//_, msg, err := conn.ReadMessage()
		m := Message{}
		err := c.ws.ReadJSON(&m)
		if err != nil {
			return err
		}
		if m.Method == "v0.env.launch" {
			if client_conf.init {
				client_conf.client.Kill()
				client_conf.init = false
			}

			EnvPlugin = shared.CreatePluginConfig(m.Body.EnvId)
			env, client_conf.client = pluginRPC(&EnvPlugin)
			client_conf.init = true
			defer client_conf.client.Kill()
			c.InitLaunch(&m)
		} else if m.Method == "v0.env.reset" {
			c.InitReset(&m)
		} else if m.Method == "v0.env.close" {
			c.Close(&m)
		}

	}

}

//What is called on Initial Launch
func (c *AgentConn) InitLaunch(m *Message) {
	log.Info("launch message received: %s\n", m.Body.EnvId)
	c.envState.SetEnvStatus("launching")
	go env.Init(m.Body.EnvId)
	//time.Sleep(2 * time.Second)
	result, err := env.Launch(m.Body.EnvId)
	if err != nil {
		log.Info("error during launch: \n", err)
	}
	c.envState.SetEnvId(m.Body.EnvId)
	log.Info("from binary received: ", result)
	c.envState.SetEnvStatus("resettng")

	c.SendLaunchReply(m.Headers.MessageId, err)
}

func (c *AgentConn) SendLaunchReply(parent_message_id int32, err error) error {

	if err != nil {
		method := "v0.reply.error"

		headers := Headers{
			Sent_at:         float64(time.Now().UnixNano() / 1000000),
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

		err = c.SendMessage(m)
		return err

	} else {

		method := "v0.reply.env.launch"

		headers := Headers{
			Sent_at:         float64(time.Now().UnixNano() / 1000000),
			EpisodeId:       c.envState.GetEpisodeId(),
			ParentMessageId: parent_message_id,
		}

		body := Body{}

		m := Message{
			Method:  method,
			Headers: headers,
			Body:    body,
		}

		err = c.SendMessage(m)
		return err

	}

}

func (c *AgentConn) InitReset(m *Message) {
	log.Info("reset message received: %s\n", m.Body.EnvId)
	c.envState.SetEnvStatus("resetting")

	result, err := env.Reset(m.Body.EnvId)
	if err != nil {
		log.Info("error during reset: \n", err)
	} else {
		c.envState.SetEnvStatus("running")
		c.envState.SetEnvId(m.Body.EnvId)
		log.Info("from binary received: ", result)

	}

	c.SendResetReply(m.Headers.MessageId, err)
}

func (c *AgentConn) SendResetReply(parent_message_id int32, err error) error {

	if err != nil {
		method := "v0.reply.error"

		headers := Headers{
			Sent_at:         float64(time.Now().UnixNano() / 1000000),
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

		err = c.SendMessage(m)
		return err

	} else {

		method := "v0.reply.env.reset"

		headers := Headers{
			Sent_at:         float64(time.Now().UnixNano() / 1000000),
			EpisodeId:       c.envState.GetEpisodeId(),
			ParentMessageId: parent_message_id,
		}

		body := Body{}

		m := Message{
			Method:  method,
			Headers: headers,
			Body:    body,
		}

		err = c.SendMessage(m)
		return err

	}

}

func (c *AgentConn) Reset() {

	c.envState.SetEnvStatus("resetting")
	c.SendEnvDescribe()

	_, err := env.Reset(c.envState.GetEnvId())
	if err != nil {
		log.Info("error during reset: \n", err)
	}

	c.envState.SetEnvStatus("running")
	//c.SendEnvDescribe()

}

func (c *AgentConn) Close(m *Message) {
	log.Info("close message received: %s\n", m.Body.EnvId)
	result, err := env.Close(m.Body.EnvId)
	if err != nil {
		log.Info("error during close: \n", err)
	}
	log.Info("from binary received: ", result)
}

//Send a message to connected Agent
func (c *AgentConn) SendMessage(m Message) error {
	c.wsLock.Lock()
	defer c.wsLock.Unlock()
	err := c.ws.WriteJSON(&m)
	if err != nil {
		log.Println("write:", err)
	}
	return err
}

func (c *AgentConn) SendEnvReward(reward float32, done bool, info string) error {

	method := "v0.env.reward"

	headers := Headers{
		Sent_at:   float64(time.Now().UnixNano() / 1000000),
		EpisodeId: c.envState.GetEpisodeId(),
	}

	body := Body{
		EnvId:  c.envState.GetEnvId(),
		Reward: reward,
		Done:   done,
		Info:   info,
	}

	m := Message{
		Method:  method,
		Headers: headers,
		Body:    body,
	}

	err := c.SendMessage(m)
	return err
}

//Required function of AgentConn
//
// func (c *AgentConn) SendEnvText(text string) {

// }
//Protobuf method GetEnvObservation (sent once every 1/fps)
func (c *AgentConn) SendEnvObservation() error {

	t, obs, err := env.GetEnvObs(c.envState.EnvId)
	if err != nil {
		log.Info(err)
	}

	observation := base64.StdEncoding.EncodeToString(obs)
	// if t == "image" {
	// 	observation := base64.StdEncoding.EncodeToString(obs)
	// } else {
	// 	observation := string(obs)
	// }

	method := "v0.env.observation"

	headers := Headers{
		Sent_at:   float64(time.Now().UnixNano() / 1000000),
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

	err = c.SendMessage(m)
	return err

}

// func (c *AgentConn) SendEnvInfo() error {

// 	t, info, err := env.GetEnvInfo(c.envState.EnvId)
// 	if err != nil {
// 		log.Info(err)
// 	}

// 	method := "v0.env.info"

// 	headers := Headers{
// 		Sent_at:   float64(time.Now().UnixNano() / 1000000),
// 		EpisodeId: c.envState.GetEpisodeId(),
// 	}

// 	body := Body{
// 		Info:     base64.StdEncoding.EncodeToString(info),
// 		InfoType: t,
// 	}

// 	m := Message{
// 		Method:  method,
// 		Headers: headers,
// 		Body:    body,
// 	}

// 	err = c.SendMessage(m)
// 	return err

// }

func (c *AgentConn) SendEnvDescribe() error {

	method := "v0.env.describe"

	headers := Headers{
		Sent_at:   float64(time.Now().UnixNano() / 1000000),
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

	err := c.SendMessage(m)
	return err

}
