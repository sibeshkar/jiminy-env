package shared

import (
	"github.com/sibeshkar/jiminy-env/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of Env that talks over RPC.
type GRPCClient struct{ client proto.EnvClient }

func (m *GRPCClient) Init(key string) (string, error) {
	resp, err := m.client.Init(context.Background(), &proto.Request{
		EnvId: key,
	})
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

func (m *GRPCClient) Launch(key string) (string, error) {
	resp, err := m.client.Launch(context.Background(), &proto.Request{
		EnvId: key,
	})
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

func (m *GRPCClient) Reset(key string) (string, error) {
	resp, err := m.client.Reset(context.Background(), &proto.Request{
		EnvId: key,
	})
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

func (m *GRPCClient) Close(key string) (string, error) {
	resp, err := m.client.Close(context.Background(), &proto.Request{
		EnvId: key,
	})
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

func (m *GRPCClient) GetReward() (float32, bool, error) {
	resp, err := m.client.GetReward(context.Background(), &proto.Empty{})
	if err != nil {
		return 0, false, err
	}

	return resp.Reward, resp.Done, nil
}

// func (m *GRPCClient) GetEnvObservation(key string) (string, []byte, error) {
// 	resp, err := m.client.GetEnvObservation(context.Background(), &proto.Request{
// 		EnvId: key,
// 	})
// 	if err != nil {
// 		return "none", []byte(""), err
// 	}

// 	return resp.Type, resp.Obs, nil
// }

func (m *GRPCClient) GetEnvObs(key string) (string, []byte, error) {
	resp, err := m.client.GetEnvObs(context.Background(), &proto.Request{
		EnvId: key,
	})
	if err != nil {
		return "none", []byte(""), err
	}

	return resp.Type, resp.Info, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Env
}

//Pre-launch stuff initialized here
func (m *GRPCServer) Init(
	ctx context.Context,
	req *proto.Request) (*proto.Response, error) {
	v, err := m.Impl.Init(req.EnvId)
	return &proto.Response{Status: v}, err
}

func (m *GRPCServer) Launch(
	ctx context.Context,
	req *proto.Request) (*proto.Response, error) {
	v, err := m.Impl.Launch(req.EnvId)
	return &proto.Response{Status: v}, err
}

func (m *GRPCServer) Reset(
	ctx context.Context,
	req *proto.Request) (*proto.Response, error) {
	v, err := m.Impl.Reset(req.EnvId)
	return &proto.Response{Status: v}, err
}

func (m *GRPCServer) Close(
	ctx context.Context,
	req *proto.Request) (*proto.Response, error) {
	v, err := m.Impl.Close(req.EnvId)
	return &proto.Response{Status: v}, err
}

func (m *GRPCServer) GetReward(
	ctx context.Context,
	req *proto.Empty) (*proto.Reward, error) {
	r, d, err := m.Impl.GetReward()
	return &proto.Reward{Reward: r, Done: d}, err
}

// func (m *GRPCServer) GetEnvObservation(
// 	ctx context.Context,
// 	req *proto.Request) (*proto.Observation, error) {
// 	t, obs, err := m.Impl.GetEnvObservation(req.EnvId)
// 	return &proto.Observation{Type: t, Obs: obs}, err
// }

func (m *GRPCServer) GetEnvObs(
	ctx context.Context,
	req *proto.Request) (*proto.Obs, error) {
	t, info, err := m.Impl.GetEnvObs(req.EnvId)
	return &proto.Obs{Type: t, Info: info}, err
}
