package tx

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/types"
	"gitlab.bianjie.ai/irita/irita-sdk-go/types/tx"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util/constant"
	"gitlab.bianjie.ai/irita/irita/modules/service"
)

type ServiceRequest struct {
	ServiceName string
	BindChainId string
	DefChainId  string
	MethodId    int16
	Provider    string
	Consumer    string
	ServiceFee  string
	Data        string
	Profiling   bool
}

func (c *client) PostServiceRequest(request ServiceRequest, memo string, commit bool) (
	types.BroadcastTxResult, error) {
	var (
		result             types.BroadcastTxResult
		consumer, provider sdk.AccAddress
		input              []byte
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
	r := request
	if consumerAcc, err := sdk.AccAddressFromBech32(r.Consumer); err != nil {
		return result, err
	} else {
		consumer = consumerAcc
	}
	if providerAcc, err := sdk.AccAddressFromBech32(r.Provider); err != nil {
		return result, err
	} else {
		provider = providerAcc
	}
	if bytes, err := hex.DecodeString(r.Data); err != nil {
		return result, err
	} else {
		input = bytes
	}
	serviceFee := []sdk.Coin{
		{
			Denom:  constant.TxDefaultFeeDenom,
			Amount: sdk.NewInt(1000000000000000000),
		},
	}
	msg := buildServiceRequestMsg(r.DefChainId, r.ServiceName, r.BindChainId, c.chainId,
		consumer, provider, r.MethodId, input, serviceFee, r.Profiling)

	// validate and sign stdMsg
	stdSignMsg := tx.StdSignMsg{
		ChainID:       c.chainId,
		AccountNumber: uint64(util.StrToInt64IgnoreErr(account.Value.AccountNumber)),
		Sequence:      uint64(util.StrToInt64IgnoreErr(account.Value.Sequence)),
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

func buildServiceRequestMsg(defChainID, defName, bindChainID, reqChainID string,
	consumer, provider sdk.AccAddress, methodID int16, input []byte, serviceFee sdk.Coins,
	profiling bool) service.MsgSvcRequest {
	msg := service.NewMsgSvcRequest(defChainID, defName, bindChainID, reqChainID,
		consumer, provider, methodID, input, serviceFee, profiling)
	return msg
}
