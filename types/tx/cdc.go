package tx

import (
	"github.com/bianjieai/irita/modules/record"
	"github.com/bianjieai/irita/modules/service"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/modules/incubator/nft"
	"github.com/tendermint/go-amino"
)

var (
	Cdc          *amino.Codec
	moduleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		service.AppModuleBasic{},
		nft.AppModuleBasic{},
		record.AppModuleBasic{},
	)
)

func init() {
	var cdc = codec.New()
	moduleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	Cdc = cdc
}
