package rpc

import (
	"encoding/hex"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/sdk-go/types/tx"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
)

type ResultTx struct {
	Hash     string                 `json:"hash"`
	Height   int64                  `json:"height"`
	Index    uint32                 `json:"index"`
	TxResult abci.ResponseDeliverTx `json:"tx_result"`
	Tx       auth.StdTx             `json:"std_tx"`
	Proof    types.TxProof          `json:"proof,omitempty"`
}

func (c *client) GetTx(hash string) (ResultTx, error) {
	var (
		res   ResultTx
		stdTx auth.StdTx
	)
	txBytes, err := hex.DecodeString(hash)
	if err != nil {
		return res, nil
	}
	txResult, err := c.rpc.Tx(txBytes, true)
	if err != nil {
		return res, nil
	}

	res.Hash = txResult.Hash.String()
	res.Height = txResult.Height
	res.Index = txResult.Index
	res.TxResult = txResult.TxResult
	res.Proof = txResult.Proof

	if err := tx.Cdc.UnmarshalBinaryLengthPrefixed(txResult.Tx, &stdTx); err != nil {
		return res, err
	} else {
		res.Tx = stdTx
	}

	return res, nil
}
