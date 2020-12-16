// Package token allows individuals and companies to create and issue their own tokens.
//

package token

import (
	"context"
	"github.com/bianjieai/irita-sdk-go/types/query"
	"strconv"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type tokenClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return tokenClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (t tokenClient) Name() string {
	return ModuleName
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
		Owner:         owner.String(),
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
		Owner:     owner.String(),
	}

	return t.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (t tokenClient) TransferToken(to string, symbol string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := t.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if err := sdk.ValidateAccAddress(to); err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgTransferTokenOwner{
		SrcOwner: owner.String(),
		DstOwner: to,
		Symbol:   symbol,
	}
	return t.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (t tokenClient) MintToken(symbol string, amount uint64, to string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := t.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	receipt := owner.String()
	if len(to) != 0 {
		if err := sdk.ValidateAccAddress(to); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		} else {
			receipt = to
		}
	}

	msg := &MsgMintToken{
		Symbol: symbol,
		Amount: amount,
		To:     receipt,
		Owner:  owner.String(),
	}
	return t.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (t tokenClient) QueryToken(denom string) (sdk.Token, error) {
	return t.BaseClient.QueryToken(denom)
}

func (t tokenClient) QueryTokens(owner string, pageReq *query.PageRequest) (sdk.Tokens, error) {
	var ownerAddr string
	if len(owner) > 0 {
		if err := sdk.ValidateAccAddress(owner); err != nil {
			return nil, sdk.Wrap(err)
		}
		ownerAddr = owner
	}

	conn, err := t.GenConn()
	defer func() { _ = conn.Close() }()

	if err != nil {
		return sdk.Tokens{}, sdk.Wrap(err)
	}

	request := &QueryTokensRequest{
		Owner:      ownerAddr,
		Pagination: pageReq,
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
	defer func() { _ = conn.Close() }()
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
	conn, err := t.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Params(
		context.Background(),
		&QueryParamsRequest{},
	)
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	return res.Params.Convert().(QueryParamsResp), nil
}
