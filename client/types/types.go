package types

import (
	"github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type (
	Coin struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}

	BroadcastTxResult struct {
		BroadcastResult ResultBroadcastTx       `json:"broadcast_result"`
		CommitResult    ResultBroadcastTxCommit `json:"commit_result"`
	}

	ResultBroadcastTx struct {
		Code uint32       `json:"code"`
		Data cmn.HexBytes `json:"data"`
		Log  string       `json:"log"`
		Hash cmn.HexBytes `json:"hash"`
	}

	ResultBroadcastTxCommit struct {
		CheckTx   abci.ResponseCheckTx   `json:"check_tx"`
		DeliverTx abci.ResponseDeliverTx `json:"deliver_tx"`
		Hash      cmn.HexBytes           `json:"hash"`
		Height    int64                  `json:"height"`
	}
)

func AccAddrFromBech32(addr string) (types.AccAddress, error) {
	return types.AccAddressFromBech32(addr)
}

func BuildCoins(txCoins types.Coins) []Coin {
	var (
		coins []Coin
	)
	if len(txCoins) == 0 {
		return coins
	}

	for _, v := range txCoins {
		coins = append(coins, Coin{
			Denom:  v.Denom,
			Amount: v.Amount.String(),
		})
	}

	return coins
}
