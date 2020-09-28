package nft

import (
	"context"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type nftClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) NFTI {
	return nftClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (nc nftClient) Name() string {
	return ModuleName
}

func (nc nftClient) RegisterCodec(cdc *codec.LegacyAmino) {
	registerCodec(cdc)
}

func (nc nftClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgIssueDenom{},
		&MsgMintNFT{},
		&MsgEditNFT{},
		&MsgTransferNFT{},
		&MsgBurnNFT{},
	)
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
		Sender: sender,
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) MintNFT(request MintNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := nc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var recipient = sender
	if len(request.Recipient) > 0 {
		recipient, err = sdk.AccAddressFromBech32(request.Recipient)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}

	msg := &MsgMintNFT{
		Id:        request.ID,
		Denom:     request.Denom,
		URI:       request.URI,
		Data:      request.Data,
		Sender:    sender,
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
		Id:     request.ID,
		Name:   request.Name,
		Denom:  request.Denom,
		URI:    request.URI,
		Data:   request.Data,
		Sender: sender,
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) TransferNFT(request TransferNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := nc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	recipient, err := sdk.AccAddressFromBech32(request.Recipient)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgTransferNFT{
		Id:        request.ID,
		Name:      request.Name,
		Denom:     request.Denom,
		URI:       request.URI,
		Data:      request.Data,
		Sender:    sender,
		Recipient: recipient,
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) BurnNFT(request BurnNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := nc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgBurnNFT{
		Sender: sender,
		Id:     request.ID,
		Denom:  request.Denom,
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) QuerySupply(denom, creator string) (uint64, sdk.Error) {
	if len(denom) == 0 {
		return 0, sdk.Wrapf("denom is required")
	}

	address, err1 := sdk.AccAddressFromBech32(creator)
	if err1 != nil {
		return 0, sdk.Wrap(err1)
	}

	conn, err := nc.GenConn()
	defer conn.Close()
	if err != nil {
		return 0, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Supply(context.Background(), &QuerySupplyRequest{
		Owner: address,
		Denom: denom,
	})
	if err != nil {
		return 0, sdk.Wrap(err)
	}

	return res.Amount, nil
}

func (nc nftClient) QueryOwner(creator, denom string) (QueryOwnerResp, sdk.Error) {
	if len(denom) == 0 {
		return QueryOwnerResp{}, sdk.Wrapf("denom is required")
	}

	address, err1 := sdk.AccAddressFromBech32(creator)
	if err1 != nil {
		return QueryOwnerResp{}, sdk.Wrap(err1)
	}

	conn, err := nc.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryOwnerResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Owner(context.Background(), &QueryOwnerRequest{
		Owner: address,
		Denom: denom,
	})
	if err != nil {
		return QueryOwnerResp{}, sdk.Wrap(err)
	}

	return res.Owner.Convert().(QueryOwnerResp), nil
}

func (nc nftClient) QueryCollection(denom string) (QueryCollectionResp, sdk.Error) {
	if len(denom) == 0 {
		return QueryCollectionResp{}, sdk.Wrapf("denom is required")
	}

	conn, err := nc.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryCollectionResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Collection(context.Background(), &QueryCollectionRequest{
		Denom: denom,
	})
	if err != nil {
		return QueryCollectionResp{}, sdk.Wrap(err)
	}

	return res.Collection.Convert().(QueryCollectionResp), nil
}

func (nc nftClient) QueryDenoms() ([]QueryDenomResp, sdk.Error) {
	conn, err := nc.GenConn()
	defer conn.Close()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Denoms(context.Background(), &QueryDenomsRequest{})
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	var denoms denoms
	denoms = res.Denoms

	return denoms.Convert().([]QueryDenomResp), nil
}

func (nc nftClient) QueryDenom(denom string) (QueryDenomResp, sdk.Error) {
	conn, err := nc.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryDenomResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Denom(context.Background(), &QueryDenomRequest{
		Denom: denom,
	})
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
	defer conn.Close()
	if err != nil {
		return QueryNFTResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).NFT(context.Background(), &QueryNFTRequest{
		Denom: denom,
		Id:    tokenID,
	})
	if err != nil {
		return QueryNFTResp{}, sdk.Wrap(err)
	}

	return res.NFT.Convert().(QueryNFTResp), nil
}
