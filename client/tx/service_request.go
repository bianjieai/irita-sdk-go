package tx

import (
	"encoding/hex"
	"github.com/bianjieai/irita-sdk-go/client/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bianjieai/irita/modules/service"
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

	serviceFee, err := types.ParseCoins(r.ServiceFee)
	if err != nil {
		return result, err
	}

	msg := buildServiceRequestMsg(r.DefChainId, r.ServiceName, r.BindChainId, c.chainId,
		consumer, provider, r.MethodId, input, serviceFee, r.Profiling)
	stdSignMsg := buildStdSignMsg(c.chainId, memo, account, msg)

	return signAndBroadcastTx(c, stdSignMsg, commit)
}

func buildServiceRequestMsg(defChainID, defName, bindChainID, reqChainID string,
	consumer, provider sdk.AccAddress, methodID int16, input []byte, serviceFee sdk.Coins,
	profiling bool) service.MsgSvcRequest {
	msg := service.NewMsgSvcRequest(defChainID, defName, bindChainID, reqChainID,
		consumer, provider, methodID, input, serviceFee, profiling)
	return msg
}
