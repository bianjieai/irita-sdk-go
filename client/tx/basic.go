package tx

import (
	"fmt"
	"github.com/bianjieai/irita-sdk-go/client/lcd"
	"github.com/bianjieai/irita-sdk-go/client/rpc"
	"github.com/bianjieai/irita-sdk-go/client/types"
	"github.com/bianjieai/irita-sdk-go/keys"
	commontypes "github.com/bianjieai/irita-sdk-go/types"
	iConfig "github.com/bianjieai/irita-sdk-go/types/config"
	"github.com/bianjieai/irita-sdk-go/util/constant"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TxClient interface {
	SendToken(receiver string, coins []types.Coin, memo string, commit bool) (types.BroadcastTxResult, error)
	PostServiceRequest(request ServiceRequest, memo string, commit bool) (types.BroadcastTxResult, error)
	PostServiceResponse(response ServiceResponse, memo string, commit bool) (types.BroadcastTxResult, error)
	MintNFT(req NFTMintReq, memo string, commit bool) (types.BroadcastTxResult, error)
	EditNFT(req NFTEditReq, memo string, commit bool) (types.BroadcastTxResult, error)
	TransferNFT(req NFTTransferReq, memo string, commit bool) (types.BroadcastTxResult, error)
	CreateRecord(req RecordCreateReq, memo string, commit bool) (types.BroadcastTxResult, error)
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

	config := sdk.GetConfig()
	addrConf := iConfig.GetIritaAddrPrefixConfig(network)
	config.SetBech32PrefixForAccount(addrConf.Conf.GetBech32AccountAddrPrefix(), addrConf.Conf.GetBech32AccountPubPrefix())
	config.SetBech32PrefixForValidator(addrConf.Conf.GetBech32ValidatorAddrPrefix(), addrConf.GetBech32ValidatorPubPrefix())
	config.SetBech32PrefixForConsensusNode(addrConf.Conf.GetBech32ConsensusAddrPrefix(), addrConf.Conf.GetBech32ConsensusPubPrefix())
	config.Seal()

	return &client{
		chainId:    chainId,
		keyManager: keyManager,
		liteClient: liteClient,
		rpcClient:  rpcClient,
	}, nil
}
