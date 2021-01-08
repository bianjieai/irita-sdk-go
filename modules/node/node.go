package node

import (
	"context"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	query "github.com/bianjieai/irita-sdk-go/types/query"
)

type nodeClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return nodeClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (n nodeClient) Name() string {
	return ModuleName
}

func (n nodeClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (n nodeClient) CreateValidator(request CreateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := n.QueryAddress(baseTx.From, baseTx.Password)
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

	return n.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (n nodeClient) UpdateValidator(request UpdateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := n.QueryAddress(baseTx.From, baseTx.Password)
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

	return n.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (n nodeClient) RemoveValidator(id string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := n.QueryAddress(baseTx.From, baseTx.Password)
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

	return n.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (n nodeClient) GrantNode(request GrantNodeRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := n.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgGrantNode{
		Name:        request.Name,
		Certificate: request.Certificate,
		Operator:    creator.String(),
	}

	return n.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (n nodeClient) RevokeNode(nodeId string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := n.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	vID, er := sdk.HexBytesFrom(nodeId)
	if er != nil {
		return sdk.ResultTx{}, sdk.Wrap(er)
	}

	msg := &MsgRevokeNode{
		Id:       vID.String(),
		Operator: creator.String(),
	}

	return n.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (n nodeClient) QueryValidators(pageReq *query.PageRequest) ([]QueryValidatorResp, sdk.Error) {
	conn, err := n.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Validators(
		context.Background(),
		&QueryValidatorsRequest{
			Pagination: pageReq,
		},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return validators(resp.Validators).Convert().([]QueryValidatorResp), nil
}

func (n nodeClient) QueryValidator(id string) (QueryValidatorResp, sdk.Error) {
	conn, err := n.GenConn()
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

func (n nodeClient) QueryNodes(pageReq *query.PageRequest) ([]QueryNodeResp, sdk.Error) {
	conn, err := n.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Nodes(
		context.Background(),
		&QueryNodesRequest{
			Pagination: pageReq,
		},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return nodes(resp.Nodes).Convert().([]QueryNodeResp), nil
}

func (n nodeClient) QueryNode(id string) (QueryNodeResp, sdk.Error) {
	conn, err := n.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryNodeResp{}, sdk.Wrap(err)
	}

	vID, err := sdk.HexBytesFrom(id)
	if err != nil {
		return QueryNodeResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Node(
		context.Background(),
		&QueryNodeRequest{
			Id: vID.String(),
		},
	)
	if err != nil {
		return QueryNodeResp{}, sdk.Wrap(err)
	}

	return resp.Node.Convert().(QueryNodeResp), nil
}

func (n nodeClient) QueryParams() (QueryParamsResp, sdk.Error) {
	conn, err := n.GenConn()
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
