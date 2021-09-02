package nft

import (
	"github.com/bianjieai/irita-sdk-go/v2/codec"
	"github.com/bianjieai/irita-sdk-go/v2/codec/types"
	cryptocodec "github.com/bianjieai/irita-sdk-go/v2/crypto/codec"
	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgIssueDenom{},
		&MsgMintNFT{},
		&MsgEditNFT{},
		&MsgTransferNFT{},
		&MsgBurnNFT{},
	)
}
