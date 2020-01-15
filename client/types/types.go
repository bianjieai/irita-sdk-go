package types

import (
	"github.com/bianjieai/irita-sdk-go/util/constant"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

var (
	DefaultTxFee = auth.NewStdFee(constant.TxDefaultGas,
		[]sdkTypes.Coin{
			{
				Denom:  constant.TxDefaultFeeDenom,
				Amount: sdkTypes.NewInt(constant.TxDefaultFeeAmount),
			},
		})
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

func AccAddrFromBech32(addr string) (sdkTypes.AccAddress, error) {
	return sdkTypes.AccAddressFromBech32(addr)
}

// ParseCoins will parse out a list of coins separated by commas.
// If nothing is provided, it returns nil Coins.
// Returned coins are sorted.
func ParseCoins(coinsStr string) (sdkTypes.Coins, error) {
	return sdkTypes.ParseCoins(coinsStr)
}

func BuildCoins(txCoins sdkTypes.Coins) []Coin {
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
