package wasm

import (
	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	cryptocodec "github.com/bianjieai/irita-sdk-go/crypto/codec"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	// ModuleName define the module name
	ModuleName = "wasm"
)

var (
	amino = codec.NewLegacyAmino()
	// ModuleCdc define the codec for wasm module
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterInterfaces regisger the implement of the msg interface for InterfaceRegistry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgStoreCode{},
		&MsgInstantiateContract{},
		&MsgExecuteContract{},
		&MsgMigrateContract{},
		&MsgUpdateAdmin{},
		&MsgClearAdmin{},
	)
}
