package modules

import (
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"sync"
)

var clientConnSingleton *grpc.ClientConn
var once sync.Once

type grpcClient struct {
}

func NewGRPCClient(url string) types.GRPCClient {
	once.Do(func() {
		dialOpts := []grpc.DialOption{
			grpc.WithInsecure(),
		}
		clientConn, err := grpc.Dial(url, dialOpts...)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}
		clientConnSingleton = clientConn
	})

	return &grpcClient{}
}

func (g grpcClient) GenConn() (*grpc.ClientConn, error) {

	return clientConnSingleton, nil
}
