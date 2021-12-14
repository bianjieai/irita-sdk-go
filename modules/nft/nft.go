package nft

import (
	"encoding/binary"
	"fmt"
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

	owner, e := sdk.AccAddressFromBech32(creator)
	if e != nil {
		return 0, sdk.Wrap(e)
	}

	param := struct {
		Denom string
		Owner sdk.AccAddress
	}{
		Denom: denom,
		Owner: owner,
	}

	bz, err := nc.Query(fmt.Sprintf(nftPath, "supply"), param)
	if err != nil {
		return 0, sdk.Wrap(err)
	}

	supply := binary.LittleEndian.Uint64(bz)
	return supply, nil
}

func (nc nftClient) QueryOwner(creator, denom string, _ *query.PageRequest) (QueryOwnerResp, sdk.Error) {
	owner, e := sdk.AccAddressFromBech32(creator)
	if e != nil {
		return QueryOwnerResp{}, sdk.Wrap(e)
	}

	param := struct {
		Denom string
		Owner sdk.AccAddress
	}{
		Denom: denom,
		Owner: owner,
	}

	var res QueryOwnerResp
	if err := nc.QueryWithResponse(fmt.Sprintf(nftPath, "owner"), param, &res); err != nil {
		return QueryOwnerResp{}, sdk.Wrap(err)
	}

	return res.Convert().(QueryOwnerResp), nil
}

func (nc nftClient) QueryCollection(denom string, pageReq *query.PageRequest) (QueryCollectionResp, sdk.Error) {
	if len(denom) == 0 {
		return QueryCollectionResp{}, sdk.Wrapf("denom is required")
	}

	param := struct {
		Denom string
	}{
		Denom: denom,
	}

	var res QueryCollectionResp
	if err := nc.QueryWithResponse(fmt.Sprintf(nftPath, "collection"), param, &res); err != nil {
		return QueryCollectionResp{}, sdk.Wrap(err)
	}

	return res.Convert().(QueryCollectionResp), nil
}

func (nc nftClient) QueryDenoms(_ *query.PageRequest) ([]QueryDenomResp, sdk.Error) {
	var res QueryDenomResps
	if err := nc.QueryWithResponse(fmt.Sprintf(nftPath, "denoms"), nil, &res); err != nil {
		return nil, sdk.Wrap(err)
	}

	return res.Convert().(QueryDenomResps), nil
}

func (nc nftClient) QueryDenom(denom string) (QueryDenomResp, sdk.Error) {
	if len(denom) == 0 {
		return QueryDenomResp{}, sdk.Wrapf("denom is required")
	}

	param := struct {
		ID string
	}{
		ID: denom,
	}

	var res QueryDenomResp
	if err := nc.QueryWithResponse(fmt.Sprintf(nftPath, "denom"), param, &res); err != nil {
		return QueryDenomResp{}, sdk.Wrap(err)
	}

	return res.Convert().(QueryDenomResp), nil
}

func (nc nftClient) QueryNFT(denom, tokenID string) (QueryNFTResp, sdk.Error) {
	if len(denom) == 0 {
		return QueryNFTResp{}, sdk.Wrapf("denom is required")
	}

	if len(tokenID) == 0 {
		return QueryNFTResp{}, sdk.Wrapf("tokenID is required")
	}

	param := struct {
		Denom   string
		TokenID string
	}{
		Denom:   denom,
		TokenID: tokenID,
	}

	var res BaseNFT
	if err := nc.QueryWithResponse(fmt.Sprintf(nftPath, "nft"), param, &res); err != nil {
		return QueryNFTResp{}, sdk.Wrap(err)
	}

	return res.Convert().(QueryNFTResp), nil
}
