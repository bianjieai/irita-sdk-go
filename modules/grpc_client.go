package modules

import (
	"sync"
	"time"

	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var clientConnSingleton *grpc.ClientConn
var once sync.Once

type grpcClient struct {
}

func NewGRPCClient(url string) types.GRPCClient {
	once.Do(func() {
		var kacp = keepalive.ClientParameters{
			Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		}

		dialOpts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithKeepaliveParams(kacp),
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
