package network

import (
	"fmt"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type GrpcAgent struct {
	name   string
	Client *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func NewGrpcAgent(cfg *GrpcCfg) (*GrpcAgent, error) {

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	ctx, cancel := context.WithTimeout(context.Background(), GrpcDialTimeout)
	defer cancel()
	client, err := grpc.DialContext(ctx, address, customDialOptions()...)
	if nil != err {
		return nil, fmt.Errorf("dial err: %s", err.Error())
	}

	ret := &GrpcAgent{
		name:   cfg.Name,
		Client: client,
		ctx:    ctx,
		cancel: cancel,
	}

	return ret, nil
}

func customDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	}
}
