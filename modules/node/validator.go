package node

import (
	"context"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	query "github.com/bianjieai/irita-sdk-go/types/query"
)

type validatorClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return validatorClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (v validatorClient) Name() string {
	return ModuleName
}

func (v validatorClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (v validatorClient) CreateValidator(request CreateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := v.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgCreateValidator{
		Name:        request.Name,
		Certificate: request.Certificate,
		Description: request.Details,
		Power:       request.Power,
		Operator:    creator.String(),
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

	msg := &MsgUpdateValidator{
		Id:          vID.String(),
		Name:        request.Name,
		Certificate: request.Certificate,
		Description: request.Details,
		Power:       request.Power,
		Operator:    creator.String(),
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
	msg := &MsgRemoveValidator{
		Id:       vID.String(),
		Operator: creator.String(),
	}

	return v.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (v validatorClient) QueryValidators(key []byte, offset uint64, limit uint64, countTotal bool) ([]QueryValidatorResp, sdk.Error) {
	conn, err := v.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Validators(
		context.Background(),
		&QueryValidatorsRequest{
			Pagination: &query.PageRequest{
				Key:        key,
				Offset:     offset,
				Limit:      limit,
				CountTotal: countTotal,
			},
		},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return validators(resp.Validators).Convert().([]QueryValidatorResp), nil
}

func (v validatorClient) QueryValidator(id string) (QueryValidatorResp, sdk.Error) {
	conn, err := v.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryValidatorResp{}, sdk.Wrap(err)
	}

	vID, err := sdk.HexBytesFrom(id)
	if err != nil {
		return QueryValidatorResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Validator(
		context.Background(),
		&QueryValidatorRequest{
			Id: vID.String(),
		},
	)
	if err != nil {
		return QueryValidatorResp{}, sdk.Wrap(err)
	}

	return resp.Validator.Convert().(QueryValidatorResp), nil
}

func (v validatorClient) QueryParams() (QueryParamsResp, sdk.Error) {
	conn, err := v.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Params(
		context.Background(),
		&QueryParamsRequest{},
	)
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	return resp.Params.Convert().(QueryParamsResp), nil
}

//func (v validatorClient) GrantValidator(id string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
//	creator, err := v.QueryAddress(baseTx.From, baseTx.Password)
//	if err != nil {
//		return sdk.ResultTx{}, sdk.Wrap(err)
//	}
//
//	vID, er := sdk.HexBytesFrom(id)
//	if er != nil {
//		return sdk.ResultTx{}, sdk.Wrap(er)
//	}
//	msg := &MsgRemoveValidator{
//		Id:       vID.String(),
//		Operator: creator.String(),
//	}
//
//	return v.BuildAndSend([]sdk.Msg{msg}, baseTx)
//}