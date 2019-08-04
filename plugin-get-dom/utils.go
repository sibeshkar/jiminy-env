package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"encoding/json"

	// "github.com/golang/protobuf/proto"
	jiminyProtos "github.com/prannayk/jiminy-protos/protos"

	"github.com/sibeshkar/jiminy-env/shared"
)

var TypeAssertionToJSONFailedError = errors.New("Type Assertion to JSON object / list failed")
var DomTagPropertyModeFailure = errors.New("DOM tag property mode failure")

//this file all ancilliary helper functions that don't go in main.go

//safely execute script in browser

func ExitOnInterrupt(cmd *exec.Cmd) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		if err := cmd.Process.Kill(); err != nil {
			log.Fatalf("Cmd: %v was not able to exit: %v", cmd, err)
		}
	}()
}
func safe_execute(script string, args []interface{}) (interface{}, error) {
	reply, err := wd.ExecuteScript(script, args)
	return reply, err
}

//Serve static files in the Init() call
func serve_static(link string) {

	//staticfilesPath = pluginConfig.Directory + "static"
	fmt.Println("starting server...")

	// server = &http.Server{
	// 	Handler: http.FileServer(http.Dir("/home/sibesh/.jiminy/plugins/sibeshkar/wob-v0/static")),
	// }
	// listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 3000})

	// if err != nil {
	// 	log.Fatal("error creating listener")
	// }

	// server.Serve(listener)

	http.Handle("/", http.FileServer(http.Dir(shared.UserHomeDir()+"/"+".jiminy/plugins/"+link+"/static")))
	http.ListenAndServe(":3000", nil)

}

func dom_to_betadom(json interface{}) ([]*jiminyProtos.DomObject, error) {
	var dom_object_list []*jiminyProtos.DomObject
	interface_list := list.New()

	interface_list.PushBack(json)

	for interface_list.Len() > 0 {
		json_element := interface_list.Front()
		json_mapv := json_element.Value
		json_object, ok := json_mapv.(map[string]interface{})
		if !ok {
			return nil, TypeAssertionToJSONFailedError
		}

		if _, ok := json_object["tag"]; !ok {
			return nil, DomTagPropertyModeFailure
		}

		tag, ok := json_object["tag"].(string)
		if !ok {
			return nil, DomTagPropertyModeFailure
		}

		id, ok := json_object["id"].(string)
		if !ok {
			return nil, DomTagPropertyModeFailure
		}

		if json_object["children"] == nil {
			continue
		}
		childrenList, ok := json_object["children"].([]interface{})
		if !ok {
			fmt.Printf("%+v\n", json_mapv)
			return nil, TypeAssertionToJSONFailedError
		}

		for _, child := range childrenList {
			interface_list.PushBack(child)
		}

		if tagInObjectList(tag, json_object["text"]) {
			bb, err := getBoundingBox(json_object)
			if err != nil {
				return nil, err
			}
			dom_object := &jiminyProtos.DomObject{
				Type:            getType(tag, id),
				Content:         getStringOrNil(json_object["text"]),
				DescriptionText: string("Empty"),
				Focused:         json_object["focused"] != nil,
				BoundingBox:     bb,
			}
			dom_object_list = append(dom_object_list, dom_object)
		}
		interface_list.Remove(interface_list.Front())
	}
	return dom_object_list, nil
}

func getBoundingBox(json_object map[string]interface{}) (*jiminyProtos.BoundingBox, error) {
	if _, ok := json_object["left"]; !ok {
		return nil, DomTagPropertyModeFailure
	}
	if _, ok := json_object["height"]; !ok {
		return nil, DomTagPropertyModeFailure
	}
	if _, ok := json_object["top"]; !ok {
		return nil, DomTagPropertyModeFailure
	}
	if _, ok := json_object["width"]; !ok {
		return nil, DomTagPropertyModeFailure
	}
	return &jiminyProtos.BoundingBox{
		X1: int32(json_object["top"].(float64)),
		Y1: int32(json_object["left"].(float64)),
		X2: int32(json_object["top"].(float64) + json_object["height"].(float64)),
		Y2: int32(json_object["left"].(float64) + json_object["width"].(float64)),
	}, nil
}

func getStringOrNil(str interface{}) string {
	if str == nil {
		return ""
	}
	return str.(string)
}

func getType(tag string, id string) string {
	if tag == "DIV" || tag == "SPAN" || tag == "P" {
		if id == "query" { return "query" }
		return "text"
	}
	if tag == "INPUT_text" {
		return "input"
	}
	if tag == "LABEL" {
		return "checkbox"
	}
	return "click"
}

func tagInObjectList(tag string, text interface{}) bool {
	if tag == "DIV" {
		return text != nil
	}
	if tag == "SPAN" {
		return true
	}
	if tag == "P" {
		return true
	}
	if tag == "INPUT_text" {
		return true
	}
	if tag == "LABEL" {
		return true
	}
	if tag == "BUTTON" {
		return true
	}
	return false
}

func process_dom(json_input interface{}) ([]byte, error) {
	betadom_object_list, err := dom_to_betadom(json_input)
	if err != nil {
		return nil, err
	}
	instance := &jiminyProtos.DomObjectInstance{
		NumObjects: int32(len(betadom_object_list)),
		Objects:    betadom_object_list,
	}

	instance_marshal, err := json.Marshal(instance)
	if err != nil {
		return nil, err
	}
	return instance_marshal, nil
}
