package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

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
	record_val      = false
	recordingDir    string
	recorder        *Recorder
)

// Here is a real implementation of Env that writes to a local file with
// the key name and the contents are the value of the key.
type Env struct{}

//Init function contains all the ancilliary services, like static file servers, VNC servers, OBS websocket etc that are initialized
//before the actual environment runtime (say a browser) starts. Important : they are background services that need to run concurrently to the main runtime.
//Can also run pre-flight checks here.
func (Env) Init(key string, record bool) (string, error) {
	//http://127.0.0.1:3000/miniwob/bisect-angle.html
	//cmd := exec.Command("sh", "-c", shared.UserHomeDir()+"/"+".jiminy/plugins/"+key+"/vendor/boxware-tigervnc", "&")
	//cmd.Start()
	record_val = record
	var err error
	if os.Getenv("DISPLAY") == "" {
		os.Setenv("DISPLAY", ":0")
		cmd := exec.Command("/bin/bash", "-c", shared.UserHomeDir()+"/"+".jiminy/plugins/"+key+"/vendor/boxware-setpasswd", "&")
		err = cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		ExitOnInterrupt(cmd)
		cmd = exec.Command("/bin/bash", "-c", shared.UserHomeDir()+"/"+".jiminy/plugins/"+key+"/vendor/boxware-tigervnc", "&")
		err = cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		ExitOnInterrupt(cmd)
		cmd = exec.Command("/bin/bash", "-c", shared.UserHomeDir()+"/"+".jiminy/plugins/"+key+"/vendor/noVNC/utils/launch.sh", "--vnc", "localhost:5900", "&")
		err = cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		ExitOnInterrupt(cmd)
	}

	if record_val {
		timeCurrent := strconv.FormatInt(time.Now().Unix(), 10)
		recordingDir := shared.UserHomeDir() + "/" + ".jiminy/plugins/" + key + "/recordings/recording_" + timeCurrent + "/"

		os.MkdirAll(recordingDir, os.ModePerm)

		proxy := NewVncProxy("", recordingDir, ":5901", "boxware", "localhost", "boxware", "5900", "dummyDesk")

		go proxy.StartListening()

		recorder, err = NewRecorder(recordingDir)
		if err != nil {
			fmt.Print("Error creating recorder %v", err)
		}

		//go recorder.StartListening()

	}

	serve_static(key)
	return "env is initialized:" + key, err
}

//Launch function contains the code to launch and handle the main environment runtime(say a browser).
func (Env) Launch(key string) (string, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

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
	for os.Getenv("DISPLAY") == "" {
		time.Sleep(100 * time.Millisecond)

	}
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

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

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

	//TODO: fix this
	if record_val {
		recorder.writeToDisk()
		recorder.NewBatch()
		recorder.AddRewardtoBatch(reward, done, "{}")
	}
	return reward, done, err
}

// func (Env) GetEnvObservation(key string) (string, []byte, error) {

// 	//obs, err := wd.Screenshot()
// 	// fmt.Println(source)
// 	// source, err := wd.PageSource()

// 	// //obs, err := base64.StdEncoding.DecodeString(source)

// 	// obs := []byte(source)
// 	reply, err := safe_execute("return core.getDOMInfo();", nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	obs, err := json.Marshal(reply)

// 	return "dom", obs, nil

// }

//Extra observation to get from episode
func (Env) GetEnvObs(key string) (string, []byte, error) {

	// instruction, err := safe_execute("return document.querySelector('#query').textContent", nil)

	// info := []byte(instruction.(string))

	reply, err := safe_execute("return core.getDOMInfo();", nil)
	if err != nil {
		fmt.Println(err)
	}

	obs, err := process_dom(reply)
	obsString := bytes.NewBuffer(obs).String()
	if record_val {
		recorder.AddObstoBatch("dom", obsString)
	}

	return "dom", obs, err

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
