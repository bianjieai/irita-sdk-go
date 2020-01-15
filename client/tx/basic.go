package tx

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bianjieai/irita-sdk-go/client/lcd"
	"github.com/bianjieai/irita-sdk-go/client/rpc"
	"github.com/bianjieai/irita-sdk-go/client/types"
	"github.com/bianjieai/irita-sdk-go/keys"
	commontypes "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/util/constant"
	iritaConfig "gitlab.bianjie.ai/irita/irita/config"
)

type TxClient interface {
	SendToken(receiver string, coins []types.Coin, memo string, commit bool) (types.BroadcastTxResult, error)
	PostServiceRequest(request ServiceRequest, memo string, commit bool) (types.BroadcastTxResult, error)
	PostServiceResponse(response ServiceResponse, memo string, commit bool) (types.BroadcastTxResult, error)
	MintNFT(req NFTMintReq, memo string, commit bool) (types.BroadcastTxResult, error)
	EditNFT(req NFTEditReq, memo string, commit bool) (types.BroadcastTxResult, error)
	TransferNFT(req NFTTransferReq, memo string, commit bool) (types.BroadcastTxResult, error)
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
	iritaConfig.SetNetworkType(network)

	config := sdk.GetConfig()
	iritaConf := iritaConfig.GetConfig()
	config.SetBech32PrefixForAccount(iritaConf.GetBech32AccountAddrPrefix(), iritaConf.GetBech32AccountPubPrefix())
	config.SetBech32PrefixForValidator(iritaConf.GetBech32ValidatorAddrPrefix(), iritaConf.GetBech32ValidatorPubPrefix())
	config.SetBech32PrefixForConsensusNode(iritaConf.GetBech32ConsensusAddrPrefix(), iritaConf.GetBech32ConsensusPubPrefix())
	config.Seal()

	return &client{
		chainId:    chainId,
		keyManager: keyManager,
		liteClient: liteClient,
		rpcClient:  rpcClient,
	}, nil
}
