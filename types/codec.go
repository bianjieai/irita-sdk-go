package types

import (
	"github.com/bianjieai/irita-sdk-go/v2/codec"
	"github.com/bianjieai/irita-sdk-go/v2/codec/types"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Marshaler
	TxConfig          TxConfig
	Amino             *codec.LegacyAmino
}
