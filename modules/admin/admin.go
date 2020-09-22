package admin

import (
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

func (a adminClient) RegisterCodec(cdc *codec.Codec) {
	registerCodec(cdc)
}

func (a adminClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddRoles{},
		&MsgRemoveRoles{},
		&MsgBlockAccount{},
		&MsgUnblockAccount{},
	)
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

	msg := MsgAddRoles{
		Address:  acc,
		Roles:    roles,
		Operator: sender,
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

	msg := MsgRemoveRoles{
		Address:  acc,
		Roles:    roles,
		Operator: sender,
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

	msg := MsgBlockAccount{
		Address:  acc,
		Operator: sender,
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

	msg := MsgUnblockAccount{
		Address:  acc,
		Operator: sender,
	}
	return a.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (a adminClient) QueryRoles(address string) (roles []Role, err sdk.Error) {
	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return roles, sdk.Wrap(err)
	}

	param := struct {
		Address sdk.AccAddress
	}{
		Address: acc,
	}

	bz, e := a.Query("custom/admin/roles", param)
	if e != nil {
		return roles, sdk.Wrap(err)
	}

	if err := a.UnmarshalJSON(bz, &roles); err != nil {
		return roles, sdk.Wrap(err)
	}
	return roles, nil
}

func (a adminClient) QueryBlacklist(page, limit int) (bl []string, err sdk.Error) {
	param := struct {
		Page  int
		Limit int
	}{
		Page:  page,
		Limit: limit,
	}

	bz, e := a.Query("custom/admin/blacklist", param)
	if e != nil {
		return nil, sdk.Wrap(err)
	}

	if err := a.UnmarshalJSON(bz, &bl); err != nil {
		return nil, sdk.Wrap(err)
	}
	return bl, nil
}
