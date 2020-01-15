package rpc

import (
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/types"
	itypes "github.com/bianjieai/irita-sdk-go/client/types"
	"github.com/bianjieai/irita-sdk-go/util/constant"
)

func (c *client) BroadcastTx(broadcastType string, tx types.Tx) (itypes.BroadcastTxResult, error) {
	var (
		broadcastTxResult itypes.BroadcastTxResult
	)
	switch broadcastType {
	case constant.TxBroadcastTypeSync:
		if result, err := c.rpc.BroadcastTxSync(tx); err != nil {
			return broadcastTxResult, err
		} else {
			broadcastTxResult.BroadcastResult = itypes.ResultBroadcastTx{
				Code: result.Code,
				Data: result.Data,
				Log:  result.Log,
				Hash: result.Hash,
			}
			return broadcastTxResult, nil
		}
	case constant.TxBroadcastTypeAsync:
		if res, err := c.rpc.BroadcastTxAsync(tx); err != nil {
			return broadcastTxResult, err
		} else {
			broadcastTxResult.BroadcastResult = itypes.ResultBroadcastTx{
				Code: res.Code,
				Data: res.Data,
				Log:  res.Log,
				Hash: res.Hash,
			}
			return broadcastTxResult, nil
		}
	case constant.TxBroadcastTypeCommit:
		if res, err := c.rpc.BroadcastTxCommit(tx); err != nil {
			return broadcastTxResult, err
		} else {
			broadcastTxResult.CommitResult = itypes.ResultBroadcastTxCommit{
				CheckTx:   res.CheckTx,
				DeliverTx: res.DeliverTx,
				Hash:      res.Hash,
				Height:    res.Height,
			}

			return broadcastTxResult, nil
		}
	default:
		return broadcastTxResult, errors.New("invalid broadcast type")
	}
}
