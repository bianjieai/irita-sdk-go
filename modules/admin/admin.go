package admin

import (
	"context"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type adminClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) AdminI {
	return adminClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (a adminClient) Name() string {
	return ModuleName
}

func (a adminClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (a adminClient) AddRoles(address string, roles []Role, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := a.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgAddRoles{
		Address:  acc.String(),
		Roles:    roles,
		Operator: sender.String(),
	}
	return a.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (a adminClient) RemoveRoles(address string, roles []Role, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := a.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgRemoveRoles{
		Address:  acc.String(),
		Roles:    roles,
		Operator: sender.String(),
	}
	return a.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (a adminClient) BlockAccount(address string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
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

func (a adminClient) UnblockAccount(address string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
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

func (a adminClient) QueryRoles(address string) ([]Role, sdk.Error) {
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

func (a adminClient) QueryBlacklist(page, limit int) ([]string, sdk.Error) {
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
