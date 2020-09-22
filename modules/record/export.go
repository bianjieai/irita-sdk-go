package record

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// expose Record module api for user
type RecordI interface {
	sdk.Module
	CreateRecord(request CreateRecordRequest, baseTx sdk.BaseTx) (string, sdk.Error)
	QueryRecord(request QueryRecordRequest) (QueryRecordResponse, sdk.Error)
}

type CreateRecordRequest struct {
	Contents []Content
}

type QueryRecordRequest struct {
	RecordID string `json:"record_id"`
	Prove    bool   `json:"prove"`
	Height   int64  `json:"height"`
}

type QueryRecordResponse struct {
	Record Data           `json:"record"`
	Proof  sdk.ProofValue `json:"proof"`
	Height int64          `json:"height"`
}
type Data struct {
	TxHash   string    `json:"tx_hash" yaml:"tx_hash"`
	Contents []Content `json:"contents" yaml:"contents"`
	Creator  string    `json:"creator" yaml:"creator"`
}
