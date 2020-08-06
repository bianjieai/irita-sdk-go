package validator

import (
	"fmt"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type validatorClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) ValidatorI {
	return validatorClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (v validatorClient) RegisterCodec(cdc *codec.Codec) {
	registerCodec(cdc)
}

func (v validatorClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateValidator{},
		&MsgUpdateValidator{},
		&MsgRemoveValidator{},
		&MsgUnjailValidator{},
	)
}

func (v validatorClient) CreateValidator(request CreateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := v.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgCreateValidator{
		Name:        request.Name,
		Certificate: request.Certificate,
		Description: request.Details,
		Power:       request.Power,
		Operator:    creator,
	}

	return v.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (v validatorClient) UpdateValidator(request UpdateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := v.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	vID, er := sdk.HexBytesFrom(request.ID)
	if er != nil {
		return sdk.ResultTx{}, sdk.Wrap(er)
	}

	msg := MsgUpdateValidator{
		ID:          vID,
		Name:        request.Name,
		Certificate: request.Certificate,
		Description: request.Details,
		Power:       request.Power,
		Operator:    creator,
	}

	return v.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (v validatorClient) RemoveValidator(id string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := v.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	vID, er := sdk.HexBytesFrom(id)
	if er != nil {
		return sdk.ResultTx{}, sdk.Wrap(er)
	}
	msg := MsgRemoveValidator{
		ID:       vID,
		Operator: creator,
	}

	return v.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (v validatorClient) Unjail(id string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := v.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	vID, er := sdk.HexBytesFrom(id)
	if er != nil {
		return sdk.ResultTx{}, sdk.Wrap(er)
	}

	msg := MsgUnjailValidator{
		ID:       vID,
		Operator: creator,
	}

	return v.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (v validatorClient) QueryValidators(page, limit int, jailed ...bool) ([]QueryValidatorResponse, sdk.Error) {
	var filter string
	if len(jailed) > 0 {
		filter = fmt.Sprintf("%v", jailed[0])
	}

	params := struct {
		Page   int
		Limit  int
		Jailed string
	}{
		Page:   page,
		Limit:  limit,
		Jailed: filter,
	}
	var vs Validators

	if err := v.QueryWithResponse("custom/validator/validators", params, &vs); err != nil {
		return nil, sdk.Wrap(err)
	}
	return vs.Convert().([]QueryValidatorResponse), nil
}

func (v validatorClient) QueryValidator(id string) (QueryValidatorResponse, sdk.Error) {
	params := struct {
		ID string
	}{
		ID: id,
	}

	var validator Validator
	if err := v.QueryWithResponse("custom/validator/validator", params, &validator); err != nil {
		return QueryValidatorResponse{}, sdk.Wrap(err)
	}
	return validator.Convert().(QueryValidatorResponse), nil
}

func (v validatorClient) QueryParams() (QueryParamsResponse, sdk.Error) {
	var p Params
	if err := v.QueryWithResponse("custom/validator/parameters", nil, &p); err != nil {
		return QueryParamsResponse{}, sdk.Wrap(err)
	}
	return p.Convert().(QueryParamsResponse), nil
}
