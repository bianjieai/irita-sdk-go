// Package modules is to warpped the API provided by each module of IRITA
//
//
package modules

import (
	"errors"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/bianjieai/irita-sdk-go/codec"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/types/tx"
	"github.com/bianjieai/irita-sdk-go/utils"
	"github.com/bianjieai/irita-sdk-go/utils/cache"
	sdklog "github.com/bianjieai/irita-sdk-go/utils/log"
)

const (
	concurrency       = 16
	cacheCapacity     = 100
	cacheExpirePeriod = 1 * time.Minute
	tryThreshold      = 3
	maxBatch          = 100
)

type baseClient struct {
	sdk.TmClient
	sdk.GRPCClient
	sdk.KeyManager
	logger         log.Logger
	cfg            *sdk.ClientConfig
	encodingConfig sdk.EncodingConfig
	l              *locker

	accountQuery
	tokenQuery
}

// NewBaseClient return the baseClient for every sub modules
func NewBaseClient(cfg sdk.ClientConfig, encodingConfig sdk.EncodingConfig, logger log.Logger) sdk.BaseClient {
	// create logger
	if logger == nil {
		logger = sdklog.NewLogger(sdklog.Config{
			Format: sdklog.FormatText,
			Level:  sdklog.InfoLevel,
		})
	}

	base := baseClient{
		TmClient:       NewRPCClient(cfg.NodeURI, encodingConfig.Amino, encodingConfig.TxConfig.TxDecoder(), logger, cfg.Timeout),
		GRPCClient:     NewGRPCClient(cfg.GRPCAddr),
		logger:         logger,
		cfg:            &cfg,
		encodingConfig: encodingConfig,
		l:              NewLocker(concurrency),
	}

	base.KeyManager = keyManager{
		keyDAO: cfg.KeyDAO,
		algo:   cfg.Algo,
	}

	c := cache.NewCache(cacheCapacity, cfg.Cached)
	base.accountQuery = accountQuery{
		Queries:    base,
		GRPCClient: base.GRPCClient,
		Logger:     base.Logger(),
		Cache:      c,
		cdc:        encodingConfig.Marshaler,
		km:         base.KeyManager,
		expiration: cacheExpirePeriod,
	}

	base.tokenQuery = tokenQuery{
		q:          base,
		GRPCClient: base.GRPCClient,
		cdc:        encodingConfig.Marshaler,
		Logger:     base.Logger(),
		Cache:      c,
	}

	return &base
}

func (base *baseClient) Logger() log.Logger {
	return base.logger
}

func (base *baseClient) SetLogger(logger log.Logger) {
	base.logger = logger
}

// Codec returns codec.
func (base *baseClient) Marshaler() codec.Marshaler {
	return base.encodingConfig.Marshaler
}

func (base *baseClient) BuildAndSend(msg []sdk.Msg, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	txByte, ctx, err := base.buildTx(msg, baseTx)
	if err != nil {
		return sdk.ResultTx{}, err
	}

	if err := base.ValidateTxSize(len(txByte), msg); err != nil {
		return sdk.ResultTx{}, err
	}
	return base.broadcastTx(txByte, ctx.Mode(), baseTx.Simulate)
}

func (base *baseClient) SendBatch(msgs sdk.Msgs, baseTx sdk.BaseTx) (rs []sdk.ResultTx, err sdk.Error) {
	if msgs == nil || len(msgs) == 0 {
		return rs, sdk.Wrapf("must have at least one message in list")
	}

	defer sdk.CatchPanic(func(errMsg string) {
		base.Logger().Error("broadcast msg failed", "errMsg", errMsg)
	})
	// validate msg
	for _, m := range msgs {
		if err := m.ValidateBasic(); err != nil {
			return rs, sdk.Wrap(err)
		}
	}
	base.Logger().Debug("validate msg success")

	// lock the account
	base.l.Lock(baseTx.From)
	defer base.l.Unlock(baseTx.From)

	batch := maxBatch
	var tryCnt = 0

resize:
	for i, ms := range utils.SubArray(batch, msgs) {
		mss := ms.(sdk.Msgs)

	retry:
		txByte, ctx, err := base.buildTx(mss, baseTx)
		if err != nil {
			return rs, err
		}

		if err := base.ValidateTxSize(len(txByte), mss); err != nil {
			base.Logger().Debug("tx is too large", "msgsLength", batch, "errMsg", err.Error())

			// filter out transactions that have been sent
			msgs = msgs[i*batch:]
			// reset the maximum number of msg in each transaction
			batch = batch / 2
			_ = base.removeCache(ctx.Address())
			goto resize
		}

		res, err := base.broadcastTx(txByte, ctx.Mode(), baseTx.Simulate)
		if err != nil {
			if sdk.Code(err.Code()) == sdk.InvalidSequence {
				base.Logger().Debug("wrong sequence,retrying ...", "address", ctx.Address(), "tryCnt", tryCnt)

				_ = base.removeCache(ctx.Address())
				if tryCnt++; tryCnt >= tryThreshold {
					return rs, err
				}
				goto retry
			}

			base.Logger().Error("broadcast transaction failed", "errMsg", err.Error())
			return rs, err
		}
		rs = append(rs, res)

		base.Logger().Info("broadcast transaction success", "txHash", res.Hash, "height", res.Height)
	}
	return rs, nil
}

func (base baseClient) QueryWithResponse(path string, data interface{}, result sdk.Response) error {
	res, err := base.Query(path, data)
	if err != nil {
		return err
	}

	if err := base.encodingConfig.Marshaler.UnmarshalJSON(res, result.(proto.Message)); err != nil {
		return err
	}

	return nil
}

func (base baseClient) Query(path string, data interface{}) ([]byte, error) {
	var bz []byte
	var err error
	if data != nil {
		bz, err = base.encodingConfig.Marshaler.MarshalJSON(data.(proto.Message))
		if err != nil {
			return nil, err
		}
	}

	opts := rpcclient.ABCIQueryOptions{
		// Height: cliCtx.Height,
		Prove: false,
	}
	result, err := base.ABCIQueryWithOptions(path, bz, opts)
	if err != nil {
		return nil, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return nil, errors.New(resp.Log)
	}

	return resp.Value, nil
}

func (base baseClient) QueryStore(key sdk.HexBytes, storeName string, height int64, prove bool) (res abci.ResponseQuery, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, "key")
	opts := rpcclient.ABCIQueryOptions{
		Prove:  prove,
		Height: height,
	}

	result, err := base.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.New(resp.Log)
	}
	return resp, nil
}

func (base *baseClient) prepare(baseTx sdk.BaseTx) (*sdk.Factory, error) {
	factory := sdk.NewFactory().
		WithChainID(base.cfg.ChainID).
		WithKeyManager(base.KeyManager).
		WithMode(base.cfg.Mode).
		WithSimulate(baseTx.Simulate).
		WithGas(base.cfg.Gas).
		WithSignModeHandler(tx.MakeSignModeHandler(tx.DefaultSignModes)).
		WithTxConfig(base.encodingConfig.TxConfig)

	addr, err := base.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	factory.WithAddress(addr.String())

	account, err := base.QueryAndRefreshAccount(addr.String())
	if err != nil {
		return nil, err
	}
	factory.WithAccountNumber(account.AccountNumber).
		WithSequence(account.Sequence).
		WithPassword(baseTx.Password)

	if !baseTx.Fee.Empty() && baseTx.Fee.IsValid() {
		fees, err := base.ToMinCoin(baseTx.Fee...)
		if err != nil {
			return nil, err
		}
		factory.WithFee(fees)
	} else {
		fees, err := base.ToMinCoin(base.cfg.Fee...)
		if err != nil {
			panic(err)
		}
		factory.WithFee(fees)
	}

	if len(baseTx.Mode) > 0 {
		factory.WithMode(baseTx.Mode)
	}

	if baseTx.Gas > 0 {
		factory.WithGas(baseTx.Gas)
	}

	if len(baseTx.Memo) > 0 {
		factory.WithMemo(baseTx.Memo)
	}
	return factory, nil
}

func (base *baseClient) ValidateTxSize(txSize int, msgs []sdk.Msg) sdk.Error {
	//var isServiceTx bool
	//for _, msg := range msgs {
	//	if msg.Route() == service.ModuleName {
	//		isServiceTx = true
	//		break
	//	}
	//}
	//if isServiceTx {
	//	var param service.Params
	//
	//	err := base.QueryParams(service.ModuleName, &param)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	if uint64(txSize) > param.TxSizeLimit {
	//		return sdk.Wrapf("tx size too large, expected: <= %d, got %d", param.TxSizeLimit, txSize)
	//	}
	//	return nil
	//
	//} else {
	//	if uint64(txSize) > base.cfg.MaxTxBytes {
	//		return sdk.Wrapf("tx size too large, expected: <= %d, got %d", base.cfg.MaxTxBytes, txSize)
	//	}
	//}
	return nil
}

type locker struct {
	shards []chan int
	size   int
}

//NewLocker implement the function of lock, can lock resources according to conditions
func NewLocker(size int) *locker {
	shards := make([]chan int, size)
	for i := 0; i < size; i++ {
		shards[i] = make(chan int, 1)
	}
	return &locker{
		shards: shards,
		size:   size,
	}
}

func (l *locker) Lock(key string) {
	ch := l.getShard(key)
	ch <- 1
}

func (l *locker) Unlock(key string) {
	ch := l.getShard(key)
	<-ch
}

func (l *locker) getShard(key string) chan int {
	index := uint(l.indexFor(key)) % uint(l.size)
	return l.shards[index]
}

func (l *locker) indexFor(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
