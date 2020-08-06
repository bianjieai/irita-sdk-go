package nft

import (
	"encoding/binary"

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

func (nc nftClient) RegisterCodec(cdc *codec.Codec) {
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

	msg := MsgIssueDenom{
		ID:     request.ID,
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

	msg := MsgMintNFT{
		ID:        request.ID,
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

	msg := MsgEditNFT{
		ID:     request.ID,
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

	msg := MsgTransferNFT{
		ID:        request.ID,
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

	msg := MsgBurnNFT{
		Sender: sender,
		ID:     request.ID,
		Denom:  request.Denom,
	}
	return nc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (nc nftClient) QuerySupply(denom, creator string) (uint64, sdk.Error) {
	if len(denom) == 0 {
		return 0, sdk.Wrapf("denom is required")
	}

	address, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return 0, sdk.Wrap(err)
	}

	param := struct {
		Denom string
		Owner sdk.AccAddress
	}{
		Denom: denom,
		Owner: address,
	}

	bz, er := nc.Query("custom/nft/supply", param)
	if er != nil {
		return 0, sdk.Wrap(err)
	}
	supply := binary.LittleEndian.Uint64(bz)
	return supply, nil
}

func (nc nftClient) QueryOwner(creator, denom string) (QueryOwnerResponse, sdk.Error) {
	if len(denom) == 0 {
		return QueryOwnerResponse{}, sdk.Wrapf("denom is required")
	}

	address, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return QueryOwnerResponse{}, sdk.Wrap(err)
	}

	param := struct {
		Denom string
		Owner sdk.AccAddress
	}{
		Denom: denom,
		Owner: address,
	}

	var owner owner
	if err := nc.QueryWithResponse("custom/nft/owner", param, &owner); err != nil {
		return QueryOwnerResponse{}, sdk.Wrap(err)
	}
	return owner.Convert().(QueryOwnerResponse), nil
}

func (nc nftClient) QueryCollection(denom string) (QueryCollectionResponse, sdk.Error) {
	if len(denom) == 0 {
		return QueryCollectionResponse{}, sdk.Wrapf("denom is required")
	}

	param := struct {
		Denom string
	}{
		Denom: denom,
	}

	var collection collection
	if err := nc.QueryWithResponse("custom/nft/collection", param, &collection); err != nil {
		return QueryCollectionResponse{}, sdk.Wrap(err)
	}
	return collection.Convert().(QueryCollectionResponse), nil
}

func (nc nftClient) QueryDenoms() ([]QueryDenomResponse, sdk.Error) {
	bz, err := nc.Query("custom/nft/denoms", nil)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	var denoms denoms
	if err := nc.UnmarshalJSON(bz, &denoms); err != nil {
		return nil, sdk.Wrap(err)
	}
	return denoms.Convert().([]QueryDenomResponse), nil
}

func (nc nftClient) QueryDenom(name string) (QueryDenomResponse, sdk.Error) {
	param := struct {
		Denom string
	}{
		Denom: name,
	}

	bz, err := nc.Query("custom/nft/denom", param)
	if err != nil {
		return QueryDenomResponse{}, sdk.Wrap(err)
	}

	var denom Denom
	if err := nc.UnmarshalJSON(bz, &denom); err != nil {
		return QueryDenomResponse{}, sdk.Wrap(err)
	}
	return denom.Convert().(QueryDenomResponse), nil
}

func (nc nftClient) QueryNFT(denom, tokenID string) (QueryNFTResponse, sdk.Error) {
	if len(denom) == 0 {
		return QueryNFTResponse{}, sdk.Wrapf("denom is required")
	}

	if len(tokenID) == 0 {
		return QueryNFTResponse{}, sdk.Wrapf("tokenID is required")
	}

	param := struct {
		Denom   string
		TokenID string
	}{
		Denom:   denom,
		TokenID: tokenID,
	}

	var nft BaseNFT
	if err := nc.QueryWithResponse("custom/nft/nft", param, &nft); err != nil {
		return QueryNFTResponse{}, sdk.Wrap(err)
	}
	return nft.Convert().(QueryNFTResponse), nil
}
