package modules

import (
	"google.golang.org/grpc"
)

type grpcClient struct {
	url string
}

func NewGRPCClient(url string) grpcClient {
	return grpcClient{
		url: url,
	}
}

func (g grpcClient) GenConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(g.url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
