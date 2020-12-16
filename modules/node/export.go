package node

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

var (
	_ Client = nodeClient{}
)

// expose Record module api for user
type Client interface {
	sdk.Module

	CreateValidator(request CreateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	UpdateValidator(request UpdateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	RemoveValidator(id string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	GrantNode(request GrantNodeRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	RevokeNode(nodeId string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryValidators(key []byte, offset uint64, limit uint64, countTotal bool) ([]QueryValidatorResp, sdk.Error)
	QueryValidator(id string) (QueryValidatorResp, sdk.Error)
	QueryNodes(key []byte, offset uint64, limit uint64, countTotal bool) ([]QueryNodeResp, sdk.Error)
	QueryNode(id string) (QueryNodeResp, sdk.Error)
	QueryParams() (QueryParamsResp, sdk.Error)
}

type CreateValidatorRequest struct {
	Name        string `json:"name"`
	Certificate string `json:"certificate"`
	Power       int64  `json:"power"`
	Details     string `json:"details"`
}

type UpdateValidatorRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Certificate string `json:"certificate"`
	Power       int64  `json:"power"`
	Details     string `json:"details"`
}

type GrantNodeRequest struct {
	Name        string `json:"name"`
	Certificate string `json:"certificate"`
	Details     string `json:"details"`
}

type QueryValidatorResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Pubkey      string `json:"pubkey"`
	Certificate string `json:"certificate"`
	Power       int64  `json:"power"`
	Details     string `json:"details"`
	Jailed      bool   `json:"jailed"`
	Operator    string `json:"operator"`
}

type QueryNodeResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Certificate string `json:"certificate"`
}

// token params
type QueryParamsResp struct {
	HistoricalEntries uint32 `json:"historical_entries"`
}
