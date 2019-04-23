package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/shared"
)

// Here is a real implementation of Env that writes to a local file with
// the key name and the contents are the value of the key.
type Env struct{}

func (Env) Launch(key string) (string, error) {
	var err error
	err = nil
	return "env is launched:" + key, err
}

func (Env) Reset(key string) (string, error) {
	var err error
	err = nil
	return "env is reset:" + key, err
}

func (Env) Close(key string) (string, error) {
	var err error
	err = nil
	return "env is closed:" + key, err
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
