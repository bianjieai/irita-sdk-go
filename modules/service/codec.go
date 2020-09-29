package service

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
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*ServiceI)(nil), nil)

	cdc.RegisterConcrete(&ServiceDefinition{}, "irismod/service/ServiceDefinition", nil)
	cdc.RegisterConcrete(&ServiceBinding{}, "irismod/service/ServiceBinding", nil)
	cdc.RegisterConcrete(&RequestContext{}, "irismod/service/RequestContext", nil)
	cdc.RegisterConcrete(&Request{}, "irismod/service/Request", nil)
	cdc.RegisterConcrete(&Response{}, "irismod/service/Response", nil)

	cdc.RegisterConcrete(&MsgDefineService{}, "irismod/service/MsgDefineService", nil)
	cdc.RegisterConcrete(&MsgBindService{}, "irismod/service/MsgBindService", nil)
	cdc.RegisterConcrete(&MsgUpdateServiceBinding{}, "irismod/service/MsgUpdateServiceBinding", nil)
	cdc.RegisterConcrete(&MsgSetWithdrawAddress{}, "irismod/service/MsgSetWithdrawAddress", nil)
	cdc.RegisterConcrete(&MsgDisableServiceBinding{}, "irismod/service/MsgDisableServiceBinding", nil)
	cdc.RegisterConcrete(&MsgEnableServiceBinding{}, "irismod/service/MsgEnableServiceBinding", nil)
	cdc.RegisterConcrete(&MsgRefundServiceDeposit{}, "irismod/service/MsgRefundServiceDeposit", nil)
	cdc.RegisterConcrete(&MsgCallService{}, "irismod/service/MsgCallService", nil)
	cdc.RegisterConcrete(&MsgRespondService{}, "irismod/service/MsgRespondService", nil)
	cdc.RegisterConcrete(&MsgPauseRequestContext{}, "irismod/service/MsgPauseRequestContext", nil)
	cdc.RegisterConcrete(&MsgStartRequestContext{}, "irismod/service/MsgStartRequestContext", nil)
	cdc.RegisterConcrete(&MsgKillRequestContext{}, "irismod/service/MsgKillRequestContext", nil)
	cdc.RegisterConcrete(&MsgUpdateRequestContext{}, "irismod/service/MsgUpdateRequestContext", nil)
	cdc.RegisterConcrete(&MsgWithdrawEarnedFees{}, "irismod/service/MsgWithdrawEarnedFees", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgDefineService{},
		&MsgBindService{},
		&MsgUpdateServiceBinding{},
		&MsgSetWithdrawAddress{},
		&MsgDisableServiceBinding{},
		&MsgEnableServiceBinding{},
		&MsgRefundServiceDeposit{},
		&MsgCallService{},
		&MsgRespondService{},
		&MsgPauseRequestContext{},
		&MsgStartRequestContext{},
		&MsgKillRequestContext{},
		&MsgUpdateRequestContext{},
		&MsgWithdrawEarnedFees{},
	)
}
