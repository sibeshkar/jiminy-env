package shared

import (
	"github.com/sibeshkar/jiminy-env/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of Env that talks over RPC.
type GRPCClient struct{ client proto.EnvClient }

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

func (m *GRPCClient) GetReward() (float32, error) {
	resp, err := m.client.GetReward(context.Background(), &proto.Empty{})
	if err != nil {
		return 0, err
	}

	return resp.Reward, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Env
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
	v, err := m.Impl.GetReward()
	return &proto.Reward{Reward: v}, err
}
