package client

import (
	"github.com/irisnet/sdk-go/client/basic"
	"github.com/irisnet/sdk-go/client/lcd"
	"github.com/irisnet/sdk-go/client/rpc"
	"github.com/irisnet/sdk-go/client/tx"
	"github.com/irisnet/sdk-go/keys"
	"github.com/irisnet/sdk-go/types"
)

type irisnetClient struct {
	basic.HttpClient
	lcd.LiteClient
	rpc.RPCClient
	tx.TxClient
}

type IRISnetClient interface {
	basic.HttpClient
	lcd.LiteClient
	rpc.RPCClient
	tx.TxClient
}

func NewIRISnetClient(baseUrl, nodeUrl string, networkType types.NetworkType, km keys.KeyManager) (IRISnetClient, error) {
	var (
		ic irisnetClient
	)
	basicClient := basic.NewClient(baseUrl)
	liteClient := lcd.NewClient(basicClient)
	rpcClient := rpc.NewClient(nodeUrl)
	status, err := rpcClient.GetStatus()
	if err != nil {
		return ic, err
	}
	chainId := status.NodeInfo.Network
	txClient, err := tx.NewClient(chainId, networkType, km, liteClient, rpcClient)
	if err != nil {
		return ic, err
	}

	ic = irisnetClient{
		HttpClient: basicClient,
		LiteClient: liteClient,
		RPCClient:  rpcClient,
		TxClient:   txClient,
	}

	return ic, nil
}
