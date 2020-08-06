package params

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

var (
	_ ParamsI = paramsClient{}
)

// expose params module api for user
type ParamsI interface {
	sdk.Module
	UpdateParams(request []UpdateParamRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}

type UpdateParamRequest struct {
	Module string      `json:"module"`
	Key    string      `json:"key"`
	Value  interface{} `json:"value"`
}
