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

func (r recordClient) RegisterCodec(cdc *codec.Codec) {
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

	msg := MsgCreateRecord{
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

func (r recordClient) QueryRecord(recordID string) (QueryRecordResponse, sdk.Error) {
	rID, err := hex.DecodeString(recordID)
	if err != nil {
		return QueryRecordResponse{}, sdk.Wrapf("invalid record id, must be hex encoded string,but got %s", recordID)
	}

	param := struct {
		RecordID []byte `json:"record_id"`
	}{
		RecordID: rID,
	}

	var record Record

	if err := r.QueryWithResponse("custom/record/record", param, &record); err != nil {
		return QueryRecordResponse{}, sdk.Wrap(err)
	}

	return record.Convert().(QueryRecordResponse), nil
}
