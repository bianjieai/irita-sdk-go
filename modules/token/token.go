// Package token allows individuals and companies to create and issue their own tokens.
//

package token

import (
	"context"
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

func (t tokenClient) Name() string {
	return ModuleName
}

func (t tokenClient) RegisterCodec(cdc *codec.LegacyAmino) {
	RegisterLegacyAminoCodec(cdc)
}

func (t tokenClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (t tokenClient) IssueToken(req IssueTokenRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := t.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgIssueToken{
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

	msg := &MsgEditToken{
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

	msg := &MsgTransferTokenOwner{
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

	msg := &MsgMintToken{
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

	conn, err := t.GenConn()
	defer func() { _ = conn.Close() }()

	if err != nil {
		return sdk.Tokens{}, sdk.Wrap(err)
	}

	request := &QueryTokensRequest{
		Owner: ownerAddr,
	}

	res, err := NewQueryClient(conn).Tokens(context.Background(), request)
	if err != nil {
		return sdk.Tokens{}, err
	}

	tokens := make(Tokens, 0, len(res.Tokens))
	for _, eviAny := range res.Tokens {
		var evi TokenInterface
		if err = t.UnpackAny(eviAny, &evi); err != nil {
			return sdk.Tokens{}, err
		}
		tokens = append(tokens, evi.(*Token))
	}

	ts := tokens.Convert().(sdk.Tokens)
	t.SaveTokens(ts...)
	return ts, nil
}

func (t tokenClient) QueryFees(symbol string) (QueryFeesResp, error) {
	conn, err := t.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryFeesResp{}, sdk.Wrap(err)
	}

	request := &QueryFeesRequest{
		Symbol: symbol,
	}

	res, err := NewQueryClient(conn).Fees(context.Background(), request)
	if err != nil {
		return QueryFeesResp{}, err
	}

	return res.Convert().(QueryFeesResp), nil
}

func (t tokenClient) QueryParams() (QueryParamsResp, error) {
	var param Params
	if err := t.BaseClient.QueryParams(ModuleName, &param); err != nil {
		return QueryParamsResp{}, err
	}
	return param.Convert().(QueryParamsResp), nil
}
