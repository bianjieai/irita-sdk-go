package record

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// expose Record module api for user
type RecordI interface {
	sdk.Module
	CreateRecord(request CreateRecordRequest, baseTx sdk.BaseTx) (string, sdk.Error)
	QueryRecord(recordID string) (QueryRecordResponse, sdk.Error)
}

type CreateRecordRequest struct {
	Contents []Content
}

type QueryRecordResponse struct {
	TxHash   string    `json:"tx_hash" yaml:"tx_hash"`
	Contents []Content `json:"contents" yaml:"contents"`
	Creator  string    `json:"creator" yaml:"creator"`
}
