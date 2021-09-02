package legacy

import (
	"github.com/bianjieai/irita-sdk-go/v2/codec"
	cryptocodec "github.com/bianjieai/irita-sdk-go/v2/crypto/codec"
)

// Cdc defines a global generic sealed Amino codec to be used throughout sdk. It
// has all Tendermint crypto and evidence types registered.
//
// TODO: Deprecated - remove this global.
var Cdc *codec.LegacyAmino

func init() {
	Cdc = codec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(Cdc)
	codec.RegisterEvidences(Cdc)
	Cdc.Seal()
}
