package tx

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/types"
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
	msg := buildBankMultiSendMsg(from, to, sdkCoins)

	account, err := c.liteClient.QueryAccount(from.String())
	if err != nil {
		return result, err
	}
	stdSignMsg := buildStdSignMsg(c.chainId, memo, account, msg)

	return signAndBroadcastTx(c, stdSignMsg, commit)
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

// buildBankMultiSendMsg builds the sending coins msg
func buildBankMultiSendMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) bank.MsgSend {
	msg := bank.NewMsgSend(from, to, coins)
	return msg
}
