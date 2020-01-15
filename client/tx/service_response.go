package tx

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bianjieai/irita-sdk-go/client/types"
	"gitlab.bianjie.ai/irita/irita/modules/service"
)

type ServiceResponse struct {
	ReqChainId string
	RequestId  string
	Data       string
	Provider   string
	ErrorMsg   string
}

func (c *client) PostServiceResponse(response ServiceResponse, memo string, commit bool) (types.BroadcastTxResult, error) {
	var (
		result   types.BroadcastTxResult
		provider sdk.AccAddress
	)

	// get account info
	from := c.keyManager.GetAddr()
	account, err := c.liteClient.QueryAccount(from.String())
	if err != nil {
		return result, err
	}

	// build stdMsg
	r := response
	if providerAcc, err := sdk.AccAddressFromBech32(r.Provider); err != nil {
		return result, err
	} else {
		provider = providerAcc
	}
	output, err := hex.DecodeString(r.Data)
	if err != nil {
		return result, err
	}
	errMsg, err := hex.DecodeString(r.ErrorMsg)
	if err != nil {
		return result, err
	}
	msg := buildServiceResponseMsg(r.ReqChainId, r.RequestId, provider, output, errMsg)
	stdSignMsg := buildStdSignMsg(c.chainId, memo, account, msg)

	return signAndBroadcastTx(c, stdSignMsg, commit)
}

func buildServiceResponseMsg(reqChainID string, requestId string,
	provider sdk.AccAddress, output, errorMsg []byte) service.MsgSvcResponse {
	return service.NewMsgSvcResponse(reqChainID, requestId, provider, output, errorMsg)
}
