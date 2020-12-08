package node

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

var (
	_ ValidatorI = validatorClient{}
)

// expose Record module api for user
type ValidatorI interface {
	sdk.Module

	CreateValidator(request CreateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	UpdateValidator(request UpdateValidatorRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	RemoveValidator(id string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryValidators(key []byte, offset uint64, limit uint64, countTotal bool) ([]QueryValidatorResp, sdk.Error)
	QueryValidator(id string) (QueryValidatorResp, sdk.Error)
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

// token params
type QueryParamsResp struct {
	HistoricalEntries uint32 `json:"historical_entries"`
}
