package rpc

import (
	itypes "github.com/irisnet/sdk-go/client/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

type RPCClient interface {
	BroadcastTx(broadcastType string, tx types.Tx) (itypes.BroadcastTxResult, error)
	GetStatus() (ResultStatus, error)
	GetTx(hash string) (ResultTx, error)
}

type client struct {
	rpc *rpcclient.HTTP
}

func NewClient(nodeUrl string) RPCClient {
	rpc := rpcclient.NewHTTP(nodeUrl, "/websocket")
	return &client{rpc: rpc}
}
