package record

import (
	"encoding/hex"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type recordClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) RecordI {
	return recordClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (r recordClient) Name() string {
	return ModuleName
}

func (r recordClient) RegisterCodec(cdc *codec.LegacyAmino) {
	registerCodec(cdc)
}

func (r recordClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateRecord{},
	)
}

func (r recordClient) CreateRecord(request CreateRecordRequest, baseTx sdk.BaseTx) (string, sdk.Error) {
	creator, err := r.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	msg := &MsgCreateRecord{
		Contents: request.Contents,
		Creator:  creator,
	}

	res, err := r.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return "", err
	}

	recordID, er := res.Events.GetValue(eventTypeCreateRecord, attributeKeyRecordID)
	if er != nil {
		return "", sdk.Wrap(er)
	}

	return recordID, nil
}

func (r recordClient) QueryRecord(request QueryRecordReq) (QueryRecordResp, sdk.Error) {
	rID, err := hex.DecodeString(request.RecordID)
	if err != nil {
		return QueryRecordResp{}, sdk.Wrapf("invalid record id, must be hex encoded string,but got %s", request.RecordID)
	}

	recordKey := GetRecordKey(rID)

	res, err := r.QueryStore(recordKey, ModuleName, request.Height, request.Prove)
	if err != nil {
		return QueryRecordResp{}, sdk.Wrap(err)
	}

	var record Record
	if err := r.Marshaler.UnmarshalBinaryBare(res.Value, &record); err != nil {
		return QueryRecordResp{}, sdk.Wrap(err)
	}

	result := record.Convert().(QueryRecordResp)

	result.Proof = sdk.ProofValue{
		Proof: r.MustMarshalJSON(res.ProofOps),
		Path:  []string{ModuleName, string(recordKey)},
		Value: res.Value,
	}
	result.Height = res.Height
	return result, nil
}
