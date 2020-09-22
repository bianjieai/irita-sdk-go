package sdk

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bianjieai/irita-sdk-go/codec"
	cdctypes "github.com/bianjieai/irita-sdk-go/codec/types"
	"github.com/bianjieai/irita-sdk-go/modules"
	"github.com/bianjieai/irita-sdk-go/modules/admin"
	"github.com/bianjieai/irita-sdk-go/modules/bank"
	"github.com/bianjieai/irita-sdk-go/modules/identity"
	"github.com/bianjieai/irita-sdk-go/modules/keys"
	"github.com/bianjieai/irita-sdk-go/modules/nft"
	"github.com/bianjieai/irita-sdk-go/modules/params"
	"github.com/bianjieai/irita-sdk-go/modules/record"
	"github.com/bianjieai/irita-sdk-go/modules/service"
	"github.com/bianjieai/irita-sdk-go/modules/token"
	"github.com/bianjieai/irita-sdk-go/modules/validator"
	"github.com/bianjieai/irita-sdk-go/std"
	"github.com/bianjieai/irita-sdk-go/types"
)

type IRITAClient struct {
	logger            log.Logger
	cdc               *codec.Codec
	appCodec          *std.Codec
	interfaceRegistry cdctypes.InterfaceRegistry
	moduleManager     map[string]types.Module

	types.BaseClient
	Token     token.TokenI
	Record    record.RecordI
	Validator validator.ValidatorI
	Identity  identity.IdentityI
	NFT       nft.NFTI
	Admin     admin.AdminI
	Params    params.ParamsI
	Bank      bank.BankI
	Service   service.ServiceI
	Key       keys.KeyI
}

func NewIRITAClient(cfg types.ClientConfig) IRITAClient {
	//create cdc for encoding and decoding
	cdc := types.NewCodec()
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	appCodec := std.NewAppCodec(cdc, interfaceRegistry)

	//create a instance of baseClient
	baseClient := modules.NewBaseClient(cfg, appCodec,nil)

	bankClient := bank.NewClient(baseClient, appCodec)
	tokenClient := token.NewClient(baseClient, appCodec)
	recordClient := record.NewClient(baseClient, appCodec)
	nftClient := nft.NewClient(baseClient, appCodec)
	serviceClient := service.NewClient(baseClient, appCodec)
	keysClient := keys.NewClient(baseClient)
	adminClient := admin.NewClient(baseClient, appCodec)
	paramsClient := params.NewClient(baseClient, appCodec)
	validatorClient := validator.NewClient(baseClient, appCodec)
	identityClient := identity.NewClient(baseClient, appCodec)

	client := &IRITAClient{
		logger:            baseClient.Logger(),
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		BaseClient:        baseClient,
		Bank:              bankClient,
		Token:             tokenClient,
		Key:               keysClient,
		Record:            recordClient,
		NFT:               nftClient,
		Service:           serviceClient,
		Admin:             adminClient,
		Params:            paramsClient,
		Validator:         validatorClient,
		Identity:          identityClient,
		moduleManager:     make(map[string]types.Module),
	}

	client.RegisterModule(
		bankClient,
		tokenClient,
		recordClient,
		nftClient,
		serviceClient,
		adminClient,
		paramsClient,
		validatorClient,
		identityClient,
	)
	return *client
}

func (client *IRITAClient) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client *IRITAClient) Codec() *codec.Codec {
	return client.cdc
}

func (client *IRITAClient) AppCodec() *std.Codec {
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
