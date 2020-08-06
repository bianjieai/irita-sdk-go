package validator

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
	Unjail(id string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryValidators(page, limit int, jailed ...bool) ([]QueryValidatorResponse, sdk.Error)
	QueryValidator(id string) (QueryValidatorResponse, sdk.Error)
	QueryParams() (QueryParamsResponse, sdk.Error)
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

type QueryValidatorResponse struct {
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
type QueryParamsResponse struct {
	HistoricalEntries uint32 `json:"historical_entries"`
}
