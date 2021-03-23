package perm

import (
	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	cryptocodec "github.com/bianjieai/irita-sdk-go/crypto/codec"
	sdk "github.com/bianjieai/irita-sdk-go/types"
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
		&MsgAssignRoles{},
		&MsgUnassignRoles{},
		&MsgBlockAccount{},
		&MsgUnblockAccount{},
	)
}
