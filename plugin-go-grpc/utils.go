package main

import (
	"fmt"
	"net/http"

	"github.com/sibeshkar/jiminy-env/shared"
)

//this file all ancilliary helper functions that don't go in main.go

//safely execute script in browser
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
