package modules

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/gogo/protobuf/jsonpb"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	clienttx "github.com/bianjieai/irita-sdk-go/client/tx"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	typetx "github.com/bianjieai/irita-sdk-go/types/tx"
)

// QueryTx returns the tx info
func (base baseClient) QueryTx(hash string) (sdk.ResultQueryTx, error) {
	tx, err := hex.DecodeString(hash)
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}

	res, err := base.Tx(context.Background(),tx, true)
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}

	resBlocks, err := base.getResultBlocks([]*ctypes.ResultTx{res})
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}
	return base.parseTxResult(res, resBlocks[res.Height])
}

func (base baseClient) QueryTxs(builder *sdk.EventQueryBuilder, page, size int) (sdk.ResultSearchTxs, error) {

	query := builder.Build()
	if len(query) == 0 {
		return sdk.ResultSearchTxs{}, errors.New("must declare at least one tag to search")
	}

	res, err := base.TxSearch(context.Background(),query, true, &page, &size, "asc")
	if err != nil {
		return sdk.ResultSearchTxs{}, err
	}

	resBlocks, err := base.getResultBlocks(res.Txs)
	if err != nil {
		return sdk.ResultSearchTxs{}, err
	}

	var txs []sdk.ResultQueryTx
	for i, tx := range res.Txs {
		txInfo, err := base.parseTxResult(tx, resBlocks[res.Txs[i].Height])
		if err != nil {
			return sdk.ResultSearchTxs{}, err
		}
		txs = append(txs, txInfo)
	}

	return sdk.ResultSearchTxs{
		Total: res.TotalCount,
		Txs:   txs,
	}, nil
}

func (base baseClient) QueryBlock(height int64) (sdk.BlockDetail, error) {
	block, err := base.Block(context.Background(),&height)
	if err != nil {
		return sdk.BlockDetail{}, err
	}

	blockResult, err := base.BlockResults(context.Background(),&height)
	if err != nil {
		return sdk.BlockDetail{}, err
	}

	return sdk.BlockDetail{
		BlockID:     block.BlockID,
		Block:       sdk.ParseBlock(base.encodingConfig.Amino, block.Block),
		BlockResult: sdk.ParseBlockResult(blockResult),
	}, nil
}

func (base *baseClient) buildTx(msgs []sdk.Msg, baseTx sdk.BaseTx) ([]byte, *clienttx.Factory, sdk.Error) {
	factory, err := base.prepare(baseTx)
	if err != nil {
		return nil, factory, sdk.Wrap(err)
	}

	txByte, err := factory.BuildAndSign(baseTx.From, msgs)
	if err != nil {
		return nil, factory, sdk.Wrap(err)
	}

	base.Logger().Debug("sign transaction success")
	return txByte, factory, nil
}

func (base baseClient) broadcastTx(txBytes []byte, mode sdk.BroadcastMode) (res sdk.ResultTx, err sdk.Error) {
	switch mode {
	case sdk.Commit:
		res, err = base.broadcastTxCommit(txBytes)
	case sdk.Async:
		res, err = base.broadcastTxAsync(txBytes)
	case sdk.Sync:
		res, err = base.broadcastTxSync(txBytes)
	default:
		err = sdk.Wrapf("commit mode(%s) not supported", mode)
	}
	return
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (base baseClient) broadcastTxCommit(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxCommit(context.Background(),tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if !res.CheckTx.IsOK() {
		return sdk.ResultTx{}, sdk.GetError(res.CheckTx.Codespace, res.CheckTx.Code, res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return sdk.ResultTx{}, sdk.GetError(res.DeliverTx.Codespace, res.DeliverTx.Code, res.DeliverTx.Log)
	}

	return sdk.ResultTx{
		GasWanted: res.DeliverTx.GasWanted,
		GasUsed:   res.DeliverTx.GasUsed,
		Events:    sdk.StringifyEvents(res.DeliverTx.Events),
		Hash:      res.Hash.String(),
		Height:    res.Height,
	}, nil
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node
// synchronously.
func (base baseClient) broadcastTxSync(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxSync(context.Background(),tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if res.Code != 0 {
		return sdk.ResultTx{}, sdk.GetError(sdk.RootCodespace,
			res.Code, res.Log)
	}

	return sdk.ResultTx{Hash: res.Hash.String()}, nil
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (base baseClient) broadcastTxAsync(tx []byte) (sdk.ResultTx, sdk.Error) {
	res, err := base.BroadcastTxAsync(context.Background(),tx)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	return sdk.ResultTx{Hash: res.Hash.String()}, nil
}

func (base baseClient) getResultBlocks(resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	resBlocks := make(map[int64]*ctypes.ResultBlock)
	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := base.Block(context.Background(),&resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}
	return resBlocks, nil
}

func (base baseClient) parseTxResult(res *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (sdk.ResultQueryTx, error) {
	var tx sdk.Tx
	var err error

	decode := base.encodingConfig.TxConfig.TxDecoder()
	if tx, err = decode(res.Tx); err != nil {
		return sdk.ResultQueryTx{}, err
	}

	unwrappedTx, err := typetx.Unwrap(tx)
	if err != nil {
		return sdk.ResultQueryTx{}, err
	}

	return sdk.ResultQueryTx{
		Hash:   res.Hash.String(),
		Height: res.Height,
		Tx:     *unwrappedTx,
		Result: sdk.TxResult{
			Code:      res.TxResult.Code,
			Log:       res.TxResult.Log,
			GasWanted: res.TxResult.GasWanted,
			GasUsed:   res.TxResult.GasUsed,
			Events:    sdk.StringifyEvents(res.TxResult.Events),
		},
		Timestamp: resBlock.Block.Time.Format(time.RFC3339),
	}, nil
}

func adjustGasEstimate(estimate uint64, adjustment float64) uint64 {
	return uint64(adjustment * float64(estimate))
}

func parseQueryResponse(bz []byte) (sdk.SimulationResponse, error) {
	var simRes sdk.SimulationResponse
	if err := jsonpb.Unmarshal(strings.NewReader(string(bz)), &simRes); err != nil {
		return sdk.SimulationResponse{}, err
	}
	return simRes, nil
}
