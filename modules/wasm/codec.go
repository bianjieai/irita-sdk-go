package wasm

import (
	"github.com/bianjieai/irita-sdk-go/v2/codec"
	"github.com/bianjieai/irita-sdk-go/v2/codec/types"
	cryptocodec "github.com/bianjieai/irita-sdk-go/v2/crypto/codec"
	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
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
