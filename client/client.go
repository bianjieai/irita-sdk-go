package client

import (
	"github.com/bianjieai/irita-sdk-go/client/basic"
	"github.com/bianjieai/irita-sdk-go/client/lcd"
	"github.com/bianjieai/irita-sdk-go/client/rpc"
	"github.com/bianjieai/irita-sdk-go/client/tx"
	"github.com/bianjieai/irita-sdk-go/keys"
	"github.com/bianjieai/irita-sdk-go/types"
)

type iritaClient struct {
	basic.HttpClient
	lcd.LiteClient
	rpc.RPCClient
	tx.TxClient
}

type IRITAClient interface {
	basic.HttpClient
	lcd.LiteClient
	rpc.RPCClient
	tx.TxClient
}

func NewIRITAClient(baseUrl, nodeUrl string, networkType types.NetworkType, km keys.KeyManager) (IRITAClient, error) {
	var (
		ic iritaClient
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

	ic = iritaClient{
		HttpClient: basicClient,
		LiteClient: liteClient,
		RPCClient:  rpcClient,
		TxClient:   txClient,
	}

	return ic, nil
}
