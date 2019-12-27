package tx

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/types"
	"gitlab.bianjie.ai/irita/irita-sdk-go/types/tx"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util/constant"
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

	// set tx fee
	fee := sdk.Coins{
		{
			Denom:  constant.TxDefaultFeeDenom,
			Amount: sdk.NewInt(constant.TxDefaultFeeAmount),
		},
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

	// validate and sign stdMsg
	stdSignMsg := tx.StdSignMsg{
		ChainID:       c.chainId,
		AccountNumber: account.Value.AccountNumber,
		Sequence:      account.Value.Sequence,
		Fee:           auth.NewStdFee(constant.TxDefaultGas, fee),
		Msgs:          []sdk.Msg{msg},
		Memo:          memo,
	}

	for _, m := range stdSignMsg.Msgs {
		if err := m.ValidateBasic(); err != nil {
			return result, err
		}
	}

	txBytes, err := c.keyManager.Sign(stdSignMsg)
	if err != nil {
		return result, err
	}

	var txBroadcastType string
	if commit {
		txBroadcastType = constant.TxBroadcastTypeCommit
	} else {
		txBroadcastType = constant.TxBroadcastTypeSync
	}

	return c.rpcClient.BroadcastTx(txBroadcastType, txBytes)
}

func buildServiceResponseMsg(reqChainID string, requestId string,
	provider sdk.AccAddress, output, errorMsg []byte) service.MsgSvcResponse {
	return service.NewMsgSvcResponse(reqChainID, requestId, provider, output, errorMsg)
}
