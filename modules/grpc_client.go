package modules

import (
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

type grpcClient struct {
	clientConnSingleton *grpc.ClientConn
}

func NewGRPCClient(url string) types.GRPCClient {
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	clientConn, err := grpc.Dial(url, dialOpts...)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	return &grpcClient{clientConnSingleton: clientConn}
}

func (g grpcClient) GenConn() (*grpc.ClientConn, error) {
	return g.clientConnSingleton, nil
}
