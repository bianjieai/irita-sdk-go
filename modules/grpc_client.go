package modules

import (
	"sync"

	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

var clientConnSingleton *grpc.ClientConn
var once sync.Once

type grpcClient struct {
}

func NewGRPCClient(url string, opts ...grpc.DialOption) types.GRPCClient {
	once.Do(func() {

		if opts == nil || len(opts) <= 0 {
			opts = []grpc.DialOption{
				grpc.WithInsecure(),
			}
		}

		clientConn, err := grpc.Dial(url, opts...)
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
