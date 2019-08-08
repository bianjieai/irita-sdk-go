package tx

import (
	"fmt"
	sdktypes "github.com/irisnet/irishub/types"
	"github.com/irisnet/sdk-go/client/lcd"
	"github.com/irisnet/sdk-go/client/rpc"
	"github.com/irisnet/sdk-go/client/types"
	"github.com/irisnet/sdk-go/keys"
	commontypes "github.com/irisnet/sdk-go/types"
	"github.com/irisnet/sdk-go/util/constant"
)

type TxClient interface {
	SendToken(receiver string, coins []types.Coin, memo string, commit bool) (types.BroadcastTxResult, error)
}

type client struct {
	chainId    string
	keyManager keys.KeyManager
	liteClient lcd.LiteClient
	rpcClient  rpc.RPCClient
}

func NewClient(chainId string, networkType commontypes.NetworkType, keyManager keys.KeyManager,
	liteClient lcd.LiteClient, rpcClient rpc.RPCClient) (TxClient, error) {
	var (
		network string
	)
	switch networkType {
	case commontypes.Mainnet:
		network = constant.NetworkTypeMainnet
	case commontypes.Testnet:
		network = constant.NetworkTypeTestnet
	default:
		return &client{}, fmt.Errorf("invalid networktype, %d", networkType)
	}
	sdktypes.SetNetworkType(network)

	return &client{
		chainId:    chainId,
		keyManager: keyManager,
		liteClient: liteClient,
		rpcClient:  rpcClient,
	}, nil
}
