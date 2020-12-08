package sdk

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bianjieai/irita-sdk-go/codec"
	cdctypes "github.com/bianjieai/irita-sdk-go/codec/types"
	cryptocodec "github.com/bianjieai/irita-sdk-go/crypto/codec"
	"github.com/bianjieai/irita-sdk-go/modules"
	"github.com/bianjieai/irita-sdk-go/modules/admin"
	"github.com/bianjieai/irita-sdk-go/modules/bank"
	"github.com/bianjieai/irita-sdk-go/modules/identity"
	"github.com/bianjieai/irita-sdk-go/modules/keys"
	"github.com/bianjieai/irita-sdk-go/modules/nft"
	"github.com/bianjieai/irita-sdk-go/modules/node"
	"github.com/bianjieai/irita-sdk-go/modules/params"
	"github.com/bianjieai/irita-sdk-go/modules/record"
	"github.com/bianjieai/irita-sdk-go/modules/service"
	"github.com/bianjieai/irita-sdk-go/modules/token"
	"github.com/bianjieai/irita-sdk-go/types"
	txtypes "github.com/bianjieai/irita-sdk-go/types/tx"
)

// IRITAClient define a group of api to access c network
type IRITAClient struct {
	logger         log.Logger
	moduleManager  map[string]types.Module
	encodingConfig types.EncodingConfig

	types.BaseClient
	Bank     bank.BankI
	Token    token.TokenI
	Record   record.RecordI
	NFT      nft.NFTI
	Service  service.ServiceI
	Key      keys.KeyI
	Admin    admin.AdminI
	Identity identity.IdentityI
	Node     node.ValidatorI
	Params   params.ParamsI
}

// NewIRITAClient return a instance of the  IRITAClient
func NewIRITAClient(cfg types.ClientConfig) IRITAClient {
	encodingConfig := makeEncodingConfig()
	//create a instance of baseClient
	baseClient := modules.NewBaseClient(cfg, encodingConfig, nil)

	adminClient := admin.NewClient(baseClient, encodingConfig.Marshaler)
	bankClient := bank.NewClient(baseClient, encodingConfig.Marshaler)
	idClient := identity.NewClient(baseClient, encodingConfig.Marshaler)
	tokenClient := token.NewClient(baseClient, encodingConfig.Marshaler)
	keysClient := keys.NewClient(baseClient)
	recordClient := record.NewClient(baseClient, encodingConfig.Marshaler)
	nftClient := nft.NewClient(baseClient, encodingConfig.Marshaler)
	serviceClient := service.NewClient(baseClient, encodingConfig.Marshaler)
	nodeClient := node.NewClient(baseClient, encodingConfig.Marshaler)
	paramsClient := params.NewClient(baseClient, encodingConfig.Marshaler)

	client := &IRITAClient{
		logger:         baseClient.Logger(),
		BaseClient:     baseClient,
		Bank:           bankClient,
		Token:          tokenClient,
		Key:            keysClient,
		Record:         recordClient,
		NFT:            nftClient,
		Service:        serviceClient,
		Admin:          adminClient,
		Identity:       idClient,
		Node:           nodeClient,
		Params:         paramsClient,
		moduleManager:  make(map[string]types.Module),
		encodingConfig: encodingConfig,
	}

	client.RegisterModule(
		adminClient,
		bankClient,
		idClient,
		tokenClient,
		recordClient,
		nftClient,
		serviceClient,
		nodeClient,
		paramsClient,
	)
	return *client
}

// SetLogger set the logger for irita client
func (client *IRITAClient) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

// Codec return the codec of the amnio
func (client *IRITAClient) Codec() *codec.LegacyAmino {
	return client.encodingConfig.Amino
}

// AppCodec return the codec of the protobuf
func (client *IRITAClient) AppCodec() codec.Marshaler {
	return client.encodingConfig.Marshaler
}

// Manager return the BaseClient
func (client *IRITAClient) Manager() types.BaseClient {
	return client.BaseClient
}

// RegisterModule regisger the module for irita client
func (client *IRITAClient) RegisterModule(ms ...types.Module) {
	for _, m := range ms {
		_, ok := client.moduleManager[m.Name()]
		if ok {
			panic(fmt.Sprintf("%s has register", m.Name()))
		}

		// m.RegisterCodec(client.encodingConfig.Amino)
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
		client.moduleManager[m.Name()] = m
	}
}

// Module return the subclient by the module name
func (client *IRITAClient) Module(name string) types.Module {
	return client.moduleManager[name]
}

func makeEncodingConfig() types.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := txtypes.NewTxConfig(marshaler, txtypes.DefaultSignModes)

	encodingConfig := types.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	RegisterLegacyAminoCodec(encodingConfig.Amino)
	RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// RegisterLegacyAminoCodec registers the sdk message type.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterInterface((*types.Tx)(nil), nil)
	cryptocodec.RegisterCrypto(cdc)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.v1beta1.Msg", (*types.Msg)(nil))
	txtypes.RegisterInterfaces(registry)
	cryptocodec.RegisterInterfaces(registry)
}
