package types

import (
	"google.golang.org/grpc"
)

func GenGRPCConn(remote string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(remote, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
