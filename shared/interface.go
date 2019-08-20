// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/sibeshkar/jiminy-env/proto"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"env_grpc": &EnvGRPCPlugin{},
}

// Env is the interface that we're exposing as a plugin.
type Env interface {
	Init(key string, record bool) (string, error)
	Launch(key string) (string, error)
	Reset(key string) (string, error)
	Close(key string) (string, error)
	GetReward() (float32, bool, error)
	//GetEnvObservation(key string) (string, []byte, error)
	GetEnvObs(key string) (string, []byte, error)
	DoAction(action []byte) (string, error)
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type EnvGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Env
}

func (p *EnvGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterEnvServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *EnvGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewEnvClient(c)}, nil
}
