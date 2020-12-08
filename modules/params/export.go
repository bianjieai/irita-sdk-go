package params

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

var (
	_ Client = paramsClient{}
)

// expose params module api for user
type Client interface {
	sdk.Module

	UpdateParams(request []UpdateParamRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}

type UpdateParamRequest struct {
	Module string `json:"module"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
