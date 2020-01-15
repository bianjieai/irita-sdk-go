package tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bianjieai/irita-sdk-go/client/lcd"
	"github.com/bianjieai/irita-sdk-go/client/types"
	"github.com/bianjieai/irita-sdk-go/types/tx"
	"github.com/bianjieai/irita-sdk-go/util/constant"
)

func buildStdSignMsg(chainId, memo string, accInfo lcd.AccountInfo, msg sdk.Msg) tx.StdSignMsg {
	return tx.StdSignMsg{
		ChainID:       chainId,
		AccountNumber: accInfo.Value.AccountNumber,
		Sequence:      accInfo.Value.Sequence,
		Fee:           types.DefaultTxFee,
		Msgs:          []sdk.Msg{msg},
		Memo:          memo,
	}
}

func signAndBroadcastTx(c *client, stdSignMsg tx.StdSignMsg, commit bool) (types.BroadcastTxResult, error) {
	var (
		result types.BroadcastTxResult
	)

	// validate and sign msg
	for _, m := range stdSignMsg.Msgs {
		if err := m.ValidateBasic(); err != nil {
			return result, err
		}
	}

	txBytes, err := c.keyManager.Sign(stdSignMsg)
	if err != nil {
		return result, err
	}

	// broadcast tx
	var txBroadcastType string
	if commit {
		txBroadcastType = constant.TxBroadcastTypeCommit
	} else {
		txBroadcastType = constant.TxBroadcastTypeSync
	}

	return c.rpcClient.BroadcastTx(txBroadcastType, txBytes)
}
