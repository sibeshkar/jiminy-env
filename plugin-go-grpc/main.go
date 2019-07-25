package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/shared"
	"github.com/tebeka/selenium"
)

var (
	seleniumPath    = "vendor/selenium-server-standalone-3.141.59.jar"
	geckoDriverPath = "vendor/geckodriver"
	port            = 8080
	pluginConfig    shared.PluginConfig
	wd              selenium.WebDriver
	service         *selenium.Service
	filesPath       string
)

// Here is a real implementation of Env that writes to a local file with
// the key name and the contents are the value of the key.
type Env struct{}

//Init function contains all the ancilliary services, like static file servers, VNC servers, OBS websocket etc that are initialized
//before the actual environment runtime (say a browser) starts. Important : they are background services that need to run concurrently to the main runtime.
//Can also run pre-flight checks here.
func (Env) Init(key string) (string, error) {
	//http://127.0.0.1:3000/miniwob/bisect-angle.html
	//cmd := exec.Command("sh", "-c", shared.UserHomeDir()+"/"+".jiminy/plugins/"+key+"/vendor/boxware-tigervnc", "&")
	//cmd.Start()
	var err error
	if os.Getenv("DISPLAY") == "" {
		os.Setenv("DISPLAY", ":0")
		cmd := exec.Command("/bin/bash", "-c", shared.UserHomeDir()+"/"+".jiminy/plugins/"+key+"/vendor/boxware-setpasswd", "&")
		err = cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		cmd = exec.Command("/bin/bash", "-c", shared.UserHomeDir()+"/"+".jiminy/plugins/"+key+"/vendor/boxware-tigervnc", "&")
		err = cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	}

	recordingDir := shared.UserHomeDir() + "/" + ".jiminy/plugins/" + key + "/recordings/"

	os.MkdirAll(recordingDir, os.ModePerm)

	proxy := create_vnc_proxy("", recordingDir, ":5901", "boxware", "localhost", "boxware", "5900", "dummyDesk")

	go proxy.StartListening()

	serve_static(key)
	return "env is initialized:" + key, err
}

//Launch function contains the code to launch and handle the main environment runtime(say a browser).
func (Env) Launch(key string) (string, error) {

	pluginConfig = shared.CreatePluginConfig(key)

	var (
		seleniumPath    = pluginConfig.Directory + "vendor/selenium-server-standalone-3.141.59.jar"
		geckoDriverPath = pluginConfig.Directory + "vendor/geckodriver"
	)

	var err error

	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath),
		selenium.Output(os.Stderr),
	}

	selenium.SetDebug(true)
	service, err = selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err)
	}
	//defer service.Stop()

	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	//defer wd.Quit()
	//serve_static()
	return "env is launched:" + key, err

}

//Reset handles resetting to a certain task of the environment. Tasks lists in tasks.go.
//Handler functions can be written as to how handle each task best.
func (Env) Reset(key string) (string, error) {

	var envString []string = strings.Split(key, "/")

	current, _ := wd.CurrentURL()

	if current == taskList[envString[len(envString)-1]] {
		return "env is already reset", nil
	} else {
		if err := wd.Get(taskList[envString[len(envString)-1]]); err != nil {
			panic(err)
		}
		var reply interface{}
		var err error

		script := "return document.readyState"
		reply, err = wd.ExecuteScript(script, nil)
		if err != nil {
			panic(err)
		}

		return "env is reset now:" + reply.(string), err

	}

}

//Closing and wrapping up environment when over.
func (Env) Close(key string) (string, error) {
	wd.Quit()
	service.Stop()
	var err error
	return "env is closed:" + key, err
}

func (Env) GetReward() (float32, bool, error) {
	script := "return WOB_REWARD_GLOBAL; "
	reply, err := safe_execute(script, nil)
	safe_execute(" window.WOB_REWARD_GLOBAL = 0;", nil)

	if err != nil {
		panic(err)
	}
	reward := float32(reply.(float64))
	reply_done, err := safe_execute(" try { return WOB_DONE_GLOBAL; } catch(err) { return false; } ", nil)
	done := reply_done.(bool)
	return reward, done, err
}

// func (Env) GetEnvObservation(key string) (string, []byte, error) {

// 	//obs, err := wd.Screenshot()
// 	// fmt.Println(source)
// 	source, err := wd.PageSource()

// 	//obs, err := base64.StdEncoding.DecodeString(source)

// 	obs := []byte(source)

// 	return "dom", obs, err

// }

//Misc info to get once every episode
func (Env) GetEnvObs(key string) (string, []byte, error) {

	instruction, err := safe_execute("return document.querySelector('#query').textContent", nil)

	info := []byte(instruction.(string))

	return "dom", info, err

}

func main() {

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"env": &shared.EnvGRPCPlugin{Impl: &Env{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})

}
