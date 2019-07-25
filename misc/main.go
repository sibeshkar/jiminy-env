package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/shared"
	log "github.com/sirupsen/logrus"
)

var EnvPlugin shared.PluginConfig

func main() {
	pluginLink := "sibeshkar/wob-v1"
	var env shared.Env
	EnvPlugin = shared.CreatePluginConfig(pluginLink)
	env, _ = pluginRPC(&EnvPlugin)

	//TODO: client.Kill() before ending process, otherwise there are zombie plugin processes
	go env.Init(pluginLink)
	env.Launch(pluginLink)

	time.Sleep(3 * time.Second)
	for {
		t, obs, err := env.GetEnvInfo("agent_conn.envState.EnvId")

		log.Infof("The type is %v, the obs is %v, error is %v:", t, obs, err)
	}
	//env.Reset(pluginLink + "/ClickButton")

}

func pluginRPC(pluginObj *shared.PluginConfig) (shared.Env, *plugin.Client) {
	//log.SetOutput(ioutil.Discard)

	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		Level:  hclog.Trace, //Uncomment this line to get more detailed plugin Trace errors
		Name:   "plugin",
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
