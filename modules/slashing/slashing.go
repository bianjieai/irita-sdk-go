package slashing

import (
	"github.com/bianjieai/irita-sdk-go/codec"
	codectypes "github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

var (
	_ Client = slashingClient{}
)

type slashingClient struct {
	sdk.BaseClient
}

//NewClient return the instance of slashing client
func NewClient(bc sdk.BaseClient) Client {
	return slashingClient{
		BaseClient: bc,
	}
}

//Name return the module name
func (slashing slashingClient) Name() string {
	return ModuleName
}

//RegisterCodec register the module struce to amion
func (slashing slashingClient) RegisterCodec(cdc *codec.LegacyAmino) {}

//RegisterInterfaceTypes register the module struce to InterfaceRegistry
func (slashing slashingClient) RegisterInterfaceTypes(registry codectypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}
