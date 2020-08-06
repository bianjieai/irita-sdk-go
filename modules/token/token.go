// Package token allows individuals and companies to create and issue their own tokens.
//
package token

import (
	"strconv"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type tokenClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) TokenI {
	return tokenClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (t tokenClient) RegisterCodec(cdc *codec.Codec) {
	registerCodec(cdc)
}

func (t tokenClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgIssueToken{},
		&MsgEditToken{},
		&MsgTransferTokenOwner{},
		&MsgMintToken{},
	)
}

func (t tokenClient) IssueToken(req IssueTokenRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := t.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgIssueToken{
		Symbol:        req.Symbol,
		Name:          req.Name,
		Scale:         req.Scale,
		MinUnit:       req.MinUnit,
		InitialSupply: req.InitialSupply,
		MaxSupply:     req.MaxSupply,
		Mintable:      req.Mintable,
		Owner:         owner,
	}

	return t.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (t tokenClient) EditToken(req EditTokenRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := t.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgEditToken{
		Symbol:    req.Symbol,
		Name:      req.Name,
		MaxSupply: req.MaxSupply,
		Mintable:  Bool(strconv.FormatBool(req.Mintable)),
		Owner:     owner,
	}

	return t.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (t tokenClient) TransferToken(to string, symbol string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := t.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	dstOwner, err := sdk.AccAddressFromBech32(to)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgTransferTokenOwner{
		SrcOwner: owner,
		DstOwner: dstOwner,
		Symbol:   symbol,
	}
	return t.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (t tokenClient) MintToken(symbol string, amount uint64, to string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := t.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	receipt := owner
	if len(to) > 0 {
		if receipt, err = sdk.AccAddressFromBech32(to); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}

	msg := MsgMintToken{
		Symbol: symbol,
		Amount: amount,
		To:     receipt,
		Owner:  owner,
	}
	return t.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (t tokenClient) QueryToken(denom string) (sdk.Token, error) {
	return t.BaseClient.QueryToken(denom)
}

func (t tokenClient) QueryTokens(owner string) (sdk.Tokens, error) {
	var ownerAddr sdk.AccAddress
	if len(owner) > 0 {
		addr, e := sdk.AccAddressFromBech32(owner)
		if e != nil {
			return nil, sdk.Wrap(e)
		}
		ownerAddr = addr
	}

	param := struct {
		Owner sdk.AccAddress
	}{
		Owner: ownerAddr,
	}

	var tokens Tokens
	bz, err := t.Query("custom/token/tokens", param)
	if err != nil {
		return sdk.Tokens{}, err
	}

	if err = t.UnmarshalJSON(bz, &tokens); err != nil {
		return sdk.Tokens{}, err
	}

	ts := tokens.Convert().(sdk.Tokens)
	t.SaveTokens(ts...)

	return ts, nil
}

func (t tokenClient) QueryFees(symbol string) (QueryFeesResponse, error) {
	param := struct {
		Symbol string
	}{
		Symbol: symbol,
	}

	var tokens tokenFees
	if err := t.QueryWithResponse("custom/token/fees", param, &tokens); err != nil {
		return QueryFeesResponse{}, err
	}

	return tokens.Convert().(QueryFeesResponse), nil
}

func (t tokenClient) QueryParams() (QueryParamsResponse, error) {
	var param Params
	if err := t.BaseClient.QueryParams(ModuleName, &param); err != nil {
		return QueryParamsResponse{}, err
	}
	return param.Convert().(QueryParamsResponse), nil
}
