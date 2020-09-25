package sdk

import (
	"fmt"
	"github.com/bianjieai/irita-sdk-go/modules/token"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bianjieai/irita-sdk-go/codec"
	cdctypes "github.com/bianjieai/irita-sdk-go/codec/types"
	"github.com/bianjieai/irita-sdk-go/modules"
	"github.com/bianjieai/irita-sdk-go/modules/bank"
	"github.com/bianjieai/irita-sdk-go/modules/keys"
	"github.com/bianjieai/irita-sdk-go/types"
)

type IRITAClient struct {
	logger            log.Logger
	cdc               *codec.LegacyAmino
	appCodec          codec.Marshaler
	interfaceRegistry cdctypes.InterfaceRegistry
	moduleManager     map[string]types.Module

	types.BaseClient
	Bank  bank.BankI
	Token token.TokenI
	Key   keys.KeyI
}

func NewIRITAClient(cfg types.ClientConfig) IRITAClient {
	//create cdc for encoding and decoding
	cdc := types.NewCodec()

	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	//create a instance of baseClient
	baseClient := modules.NewBaseClient(cfg, cdc, marshaler, nil)

	bankClient := bank.NewClient(baseClient, marshaler)
	tokenClient := token.NewClient(baseClient, marshaler)
	keysClient := keys.NewClient(baseClient)

	client := &IRITAClient{
		logger:            baseClient.Logger(),
		cdc:               cdc,
		appCodec:          marshaler,
		interfaceRegistry: interfaceRegistry,
		BaseClient:        baseClient,
		Bank:              bankClient,
		Token:             tokenClient,
		Key:               keysClient,
		moduleManager:     make(map[string]types.Module),
	}

	client.RegisterModule(
		bankClient,
		tokenClient,
	)
	return *client
}

func (client *IRITAClient) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client *IRITAClient) Codec() *codec.LegacyAmino {
	return client.cdc
}

func (client *IRITAClient) AppCodec() codec.Marshaler {
	return client.appCodec
}

func (client *IRITAClient) Manager() types.BaseClient {
	return client.BaseClient
}

func (client *IRITAClient) RegisterModule(ms ...types.Module) {
	for _, m := range ms {
		_, ok := client.moduleManager[m.Name()]
		if ok {
			panic(fmt.Sprintf("%s has register", m.Name()))
		}

		m.RegisterCodec(client.cdc)
		m.RegisterInterfaceTypes(client.interfaceRegistry)
		client.moduleManager[m.Name()] = m
	}
}

func (client *IRITAClient) Module(name string) types.Module {
	return client.moduleManager[name]
}
