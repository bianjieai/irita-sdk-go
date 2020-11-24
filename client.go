package sdk

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bianjieai/irita-sdk-go/codec"
	cdctypes "github.com/bianjieai/irita-sdk-go/codec/types"
	cryptocodec "github.com/bianjieai/irita-sdk-go/crypto/codec"
	"github.com/bianjieai/irita-sdk-go/modules"
	"github.com/bianjieai/irita-sdk-go/modules/bank"
	"github.com/bianjieai/irita-sdk-go/modules/keys"
	"github.com/bianjieai/irita-sdk-go/modules/nft"
	"github.com/bianjieai/irita-sdk-go/modules/record"
	"github.com/bianjieai/irita-sdk-go/modules/service"
	"github.com/bianjieai/irita-sdk-go/modules/token"
	"github.com/bianjieai/irita-sdk-go/types"
	txtypes "github.com/bianjieai/irita-sdk-go/types/tx"
)

type IRITAClient struct {
	logger         log.Logger
	moduleManager  map[string]types.Module
	encodingConfig types.EncodingConfig

	types.BaseClient
	Bank    bank.BankI
	Token   token.TokenI
	Record  record.RecordI
	NFT     nft.NFTI
	Service service.ServiceI
	Key     keys.KeyI
}

func NewIRITAClient(cfg types.ClientConfig) IRITAClient {
	encodingConfig := makeEncodingConfig()
	//create a instance of baseClient
	baseClient := modules.NewBaseClient(cfg, encodingConfig, nil)

	bankClient := bank.NewClient(baseClient, encodingConfig.Marshaler)
	tokenClient := token.NewClient(baseClient, encodingConfig.Marshaler)
	keysClient := keys.NewClient(baseClient)
	recordClient := record.NewClient(baseClient, encodingConfig.Marshaler)
	nftClient := nft.NewClient(baseClient, encodingConfig.Marshaler)
	serviceClient := service.NewClient(baseClient, encodingConfig.Marshaler)

	client := &IRITAClient{
		logger:         baseClient.Logger(),
		BaseClient:     baseClient,
		Bank:           bankClient,
		Token:          tokenClient,
		Key:            keysClient,
		Record:         recordClient,
		NFT:            nftClient,
		Service:        serviceClient,
		moduleManager:  make(map[string]types.Module),
		encodingConfig: encodingConfig,
	}

	client.RegisterModule(
		bankClient,
		tokenClient,
		recordClient,
		nftClient,
		serviceClient,
	)
	return *client
}

func (client *IRITAClient) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client *IRITAClient) Codec() *codec.LegacyAmino {
	return client.encodingConfig.Amino
}

func (client *IRITAClient) AppCodec() codec.Marshaler {
	return client.encodingConfig.Marshaler
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

		// m.RegisterCodec(client.encodingConfig.Amino)
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
		client.moduleManager[m.Name()] = m
	}
}

func (client *IRITAClient) Module(name string) types.Module {
	return client.moduleManager[name]
}

func makeEncodingConfig() types.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := txtypes.NewTxConfig(marshaler,txtypes.DefaultSignModes)

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
