package slashing

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// Client define a group of interface for wasm module
type Client interface {
	sdk.Module
}
