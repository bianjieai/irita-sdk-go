package nft

import (
	"context"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	"github.com/bianjieai/irita-sdk-go/types/query"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type nftClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return nftClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (nc nftClient) Name() string {
	return ModuleName
}

func (nc nftClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (nc nftClient) IssueDenom(request IssueDenomRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := nc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgIssueDenom{
		Id:     request.ID,
		Name:   request.Name,
		Schema: request.Schema,
		Sender: sender.String(),
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) MintNFT(request MintNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := nc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var recipient = sender.String()
	if len(request.Recipient) > 0 {
		if err := sdk.ValidateAccAddress(request.Recipient); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		recipient = request.Recipient
	}

	msg := &MsgMintNFT{
		Id:        request.ID,
		DenomId:   request.Denom,
		Name:      request.Name,
		URI:       request.URI,
		Data:      request.Data,
		Sender:    sender.String(),
		Recipient: recipient,
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) EditNFT(request EditNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := nc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgEditNFT{
		Id:      request.ID,
		Name:    request.Name,
		DenomId: request.Denom,
		URI:     request.URI,
		Data:    request.Data,
		Sender:  sender.String(),
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) TransferNFT(request TransferNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := nc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if err := sdk.ValidateAccAddress(request.Recipient); err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgTransferNFT{
		Id:        request.ID,
		Name:      request.Name,
		DenomId:   request.Denom,
		URI:       request.URI,
		Data:      request.Data,
		Sender:    sender.String(),
		Recipient: request.Recipient,
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) BurnNFT(request BurnNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := nc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgBurnNFT{
		Sender:  sender.String(),
		Id:      request.ID,
		DenomId: request.Denom,
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) QuerySupply(denom, creator string) (uint64, sdk.Error) {
	if len(denom) == 0 {
		return 0, sdk.Wrapf("denom is required")
	}

	if err := sdk.ValidateAccAddress(creator); err != nil {
		return 0, sdk.Wrap(err)
	}

	conn, err := nc.GenConn()

	if err != nil {
		return 0, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Supply(
		context.Background(),
		&QuerySupplyRequest{
			Owner:   creator,
			DenomId: denom,
		},
	)
	if err != nil {
		return 0, sdk.Wrap(err)
	}

	return res.Amount, nil
}

func (nc nftClient) QueryOwner(creator, denom string, pageReq *query.PageRequest) (QueryOwnerResp, sdk.Error) {
	if err := sdk.ValidateAccAddress(creator); err != nil {
		return QueryOwnerResp{}, sdk.Wrap(err)
	}

	conn, err := nc.GenConn()

	if err != nil {
		return QueryOwnerResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Owner(
		context.Background(),
		&QueryOwnerRequest{
			Owner:      creator,
			DenomId:    denom,
			Pagination: pageReq,
		},
	)
	if err != nil {
		return QueryOwnerResp{}, sdk.Wrap(err)
	}

	return res.Owner.Convert().(QueryOwnerResp), nil
}

func (nc nftClient) QueryCollection(denom string, pageReq *query.PageRequest) (QueryCollectionResp, sdk.Error) {
	if len(denom) == 0 {
		return QueryCollectionResp{}, sdk.Wrapf("denom is required")
	}

	conn, err := nc.GenConn()

	if err != nil {
		return QueryCollectionResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Collection(
		context.Background(),
		&QueryCollectionRequest{
			DenomId:    denom,
			Pagination: pageReq,
		},
	)
	if err != nil {
		return QueryCollectionResp{}, sdk.Wrap(err)
	}

	return res.Collection.Convert().(QueryCollectionResp), nil
}

func (nc nftClient) QueryDenoms(pageReq *query.PageRequest) ([]QueryDenomResp, sdk.Error) {
	conn, err := nc.GenConn()

	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Denoms(
		context.Background(),
		&QueryDenomsRequest{Pagination: pageReq},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return denoms(res.Denoms).Convert().([]QueryDenomResp), nil
}

func (nc nftClient) QueryDenom(denom string) (QueryDenomResp, sdk.Error) {
	conn, err := nc.GenConn()

	if err != nil {
		return QueryDenomResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Denom(
		context.Background(),
		&QueryDenomRequest{DenomId: denom},
	)
	if err != nil {
		return QueryDenomResp{}, sdk.Wrap(err)
	}

	return res.Denom.Convert().(QueryDenomResp), nil
}

func (nc nftClient) QueryNFT(denom, tokenID string) (QueryNFTResp, sdk.Error) {
	if len(denom) == 0 {
		return QueryNFTResp{}, sdk.Wrapf("denom is required")
	}

	if len(tokenID) == 0 {
		return QueryNFTResp{}, sdk.Wrapf("tokenID is required")
	}

	conn, err := nc.GenConn()

	if err != nil {
		return QueryNFTResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).NFT(
		context.Background(),
		&QueryNFTRequest{
			DenomId: denom,
			TokenId: tokenID,
		},
	)
	if err != nil {
		return QueryNFTResp{}, sdk.Wrap(err)
	}

	return res.NFT.Convert().(QueryNFTResp), nil
}
