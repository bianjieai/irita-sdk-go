package tx

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/sdk-go/client/types"
	"github.com/irisnet/sdk-go/types/tx"
	"github.com/irisnet/sdk-go/util"
	"github.com/irisnet/sdk-go/util/constant"
)

func (c *client) SendToken(receiver string, coins []types.Coin, memo string, commit bool) (types.BroadcastTxResult, error) {
	var (
		result types.BroadcastTxResult
	)
	from := c.keyManager.GetAddr()

	to, err := types.AccAddrFromBech32(receiver)
	if err != nil {
		return result, err
	}

	sdkCoins, err := buildCoins(coins)
	if err != nil {
		return result, err
	}
	msg := buildBankSendMsg(from, to, sdkCoins)

	account, err := c.liteClient.QueryAccount(from.String())
	if err != nil {
		return result, err
	}

	// TODO: check balance is enough

	fee := sdk.Coins{
		{
			Denom:  constant.TxDefaultFeeDenom,
			Amount: sdk.NewInt(constant.TxDefaultFeeAmount),
		},
	}
	stdSignMsg := tx.StdSignMsg{
		ChainID:       c.chainId,
		AccountNumber: uint64(util.StrToInt64IgnoreErr(account.Value.AccountNumber)),
		Sequence:      uint64(util.StrToInt64IgnoreErr(account.Value.Sequence)),
		Fee:           auth.NewStdFee(constant.TxDefaultGas, fee...),
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

func buildCoins(icoins []types.Coin) (sdk.Coins, error) {
	var (
		coins []sdk.Coin
	)
	if len(icoins) == 0 {
		return coins, nil
	}
	for _, v := range icoins {
		amount, ok := sdk.NewIntFromString(v.Amount)
		if ok {
			coins = append(coins, sdk.Coin{
				Denom:  v.Denom,
				Amount: amount,
			})
		} else {
			return coins, fmt.Errorf("can't parse str to Int, coin is %+v", icoins)
		}
	}

	return coins, nil
}

// buildBankSendMsg builds the sending coins msg
func buildBankSendMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) bank.MsgSend {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}
