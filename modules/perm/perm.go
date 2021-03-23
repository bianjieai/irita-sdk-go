package perm

import (
	"context"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type permClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return permClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (a permClient) Name() string {
	return ModuleName
}

func (a permClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (a permClient) AssignRoles(address string, roles []Role, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := a.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgAssignRoles{
		Address:  acc.String(),
		Roles:    roles,
		Operator: sender.String(),
	}
	return a.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (a permClient) UnassignRoles(address string, roles []Role, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := a.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgUnassignRoles{
		Address:  acc.String(),
		Roles:    roles,
		Operator: sender.String(),
	}
	return a.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (a permClient) BlockAccount(address string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := a.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgBlockAccount{
		Address:  acc.String(),
		Operator: sender.String(),
	}
	return a.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (a permClient) UnblockAccount(address string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := a.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgUnblockAccount{
		Address:  acc.String(),
		Operator: sender.String(),
	}
	return a.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (a permClient) QueryRoles(address string) ([]Role, sdk.Error) {
	conn, err := a.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Roles(
		context.Background(),
		&QueryRolesRequest{Address: acc.String()},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return resp.Roles, nil
}

func (a permClient) QueryBlacklist(page, limit int) ([]string, sdk.Error) {
	conn, err := a.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Blacklist(
		context.Background(),
		&QueryBlacklistRequest{},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return resp.Addresses, nil
}
