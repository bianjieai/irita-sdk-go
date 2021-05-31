package legacy

import (
	"github.com/bianjieai/irita-sdk-go/codec"
	cryptocodec "github.com/bianjieai/irita-sdk-go/crypto/codec"
	"github.com/bianjieai/irita-sdk-go/modules/bank"
	"github.com/bianjieai/irita-sdk-go/modules/token"
)

// Cdc defines a global generic sealed amino codec to be used throughout sdk. It
// has all Tendermint crypto and evidence types registered.
//
// TODO: Deprecated - remove this global.
var Cdc *codec.LegacyAmino

func init() {
	Cdc = codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(Cdc)
	codec.RegisterEvidences(Cdc)
	bank.RegisterLegacyAminoCodec(Cdc)
	token.RegisterLegacyAminoCodec(Cdc)
	Cdc.Seal()
}
