package modules

import (
	"context"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	rpc "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bianjieai/irita-sdk-go/codec"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/utils/uuid"
)

type rpcClient struct {
	rpc.Client
	log.Logger
	cdc    *codec.Codec
}

func NewRPCClient(remote string,
	cdc *codec.Codec,
	logger log.Logger,
	timeout uint) sdk.TmClient {
	client, err := rpchttp.NewWithTimeout(remote, "/websocket", timeout)
	if err != nil {
		panic(err)
	}

	_ = client.Start()
	return rpcClient{
		Client: client,
		Logger: logger,
		cdc:    cdc,
	}
}

//=============================================================================
//SubscribeNewBlock implement WSClient interface
func (r rpcClient) SubscribeNewBlock(builder *sdk.EventQueryBuilder,
	handler sdk.EventNewBlockHandler) (sdk.Subscription, sdk.Error) {
	if builder == nil {
		builder = sdk.NewEventQueryBuilder()
	}

	builder.AddCondition(sdk.Cond(sdk.TypeKey).EQ(tmtypes.EventNewBlock))
	query := builder.Build()

	return r.SubscribeAny(query, func(data sdk.EventData) {
		handler(data.(sdk.EventDataNewBlock))
	})
}

//SubscribeTx implement WSClient interface
func (r rpcClient) SubscribeTx(builder *sdk.EventQueryBuilder, handler sdk.EventTxHandler) (sdk.Subscription, sdk.Error) {
	if builder == nil {
		builder = sdk.NewEventQueryBuilder()
	}
	query := builder.AddCondition(sdk.Cond(sdk.TypeKey).EQ(sdk.TxValue)).Build()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		handler(data.(sdk.EventDataTx))
	})
}

func (r rpcClient) SubscribeNewBlockHeader(handler sdk.EventNewBlockHeaderHandler) (sdk.Subscription, sdk.Error) {
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		handler(data.(sdk.EventDataNewBlockHeader))
	})
}

func (r rpcClient) SubscribeValidatorSetUpdates(handler sdk.EventValidatorSetUpdatesHandler) (sdk.Subscription, sdk.Error) {
	query := tmtypes.QueryForEvent(tmtypes.EventValidatorSetUpdates).String()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		handler(data.(sdk.EventDataValidatorSetUpdates))
	})
}

func (r rpcClient) Resubscribe(subscription sdk.Subscription, handler sdk.EventHandler) (err sdk.Error) {
	_, err = r.SubscribeAny(subscription.Query, handler)
	return
}

func (r rpcClient) Unsubscribe(subscription sdk.Subscription) sdk.Error {
	r.Info("end to subscribe event","query", subscription.Query,"subscriber", subscription.ID)
	err := r.Client.Unsubscribe(subscription.Ctx, subscription.ID, subscription.Query)
	if err != nil {
		r.Error("unsubscribe failed","query", subscription.Query,"subscriber", subscription.ID,"errMsg",err.Error())
		return sdk.Wrap(err)
	}
	return nil
}

func (r rpcClient) SubscribeAny(query string, handler sdk.EventHandler) (subscription sdk.Subscription, err sdk.Error) {
	ctx := context.Background()
	subscriber := getSubscriber()
	ch, e := r.Subscribe(ctx, subscriber, query, 0)
	if e != nil {
		return subscription, sdk.Wrap(e)
	}

	r.Info("subscribe event","query", subscription.Query,"subscriber", subscription.ID)

	subscription = sdk.Subscription{
		Ctx:   ctx,
		Query: query,
		ID:    subscriber,
	}

	go func() {
		for {
			data := <-ch
			go func() {
				defer sdk.CatchPanic(func(errMsg string) {
					r.Error("unsubscribe failed","query", subscription.Query,"subscriber", subscription.ID,"errMsg",err.Error())
				})

				switch data := data.Data.(type) {
				case tmtypes.EventDataTx:
					handler(r.parseTx(data))
					return
				case tmtypes.EventDataNewBlock:
					handler(r.parseNewBlock(data))
					return
				case tmtypes.EventDataNewBlockHeader:
					handler(r.parseNewBlockHeader(data))
					return
				case tmtypes.EventDataValidatorSetUpdates:
					handler(r.parseValidatorSetUpdates(data))
					return
				default:
					handler(data)
				}
			}()
		}
	}()
	return
}

func (r rpcClient) parseTx(data sdk.EventData) sdk.EventDataTx {
	tx := data.(tmtypes.EventDataTx)
	var stdTx sdk.StdTx
	if err := r.cdc.UnmarshalBinaryBare(tx.Tx, &stdTx); err != nil {
		return sdk.EventDataTx{}
	}
	hash := sdk.HexBytes(tx.Tx.Hash()).String()
	result := sdk.TxResult{
		Code:      tx.Result.Code,
		Log:       tx.Result.Log,
		GasWanted: tx.Result.GasWanted,
		GasUsed:   tx.Result.GasUsed,
		Events:    sdk.ParseEvents(tx.Result.Events),
	}
	return sdk.EventDataTx{
		Hash:   hash,
		Height: tx.Height,
		Index:  tx.Index,
		Tx:     stdTx,
		Result: result,
	}
}

func (r rpcClient) parseNewBlock(data sdk.EventData) sdk.EventDataNewBlock {
	block := data.(tmtypes.EventDataNewBlock)
	return sdk.EventDataNewBlock{
		Block: sdk.ParseBlock(r.cdc, block.Block),
		ResultBeginBlock: sdk.ResultBeginBlock{
			Events: sdk.ParseEvents(block.ResultBeginBlock.Events),
		},
		ResultEndBlock: sdk.ResultEndBlock{
			Events:           sdk.ParseEvents(block.ResultEndBlock.Events),
			ValidatorUpdates: sdk.ParseValidatorUpdate(block.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseNewBlockHeader(data sdk.EventData) sdk.EventDataNewBlockHeader {
	blockHeader := data.(tmtypes.EventDataNewBlockHeader)
	return sdk.EventDataNewBlockHeader{
		Header: blockHeader.Header,
		ResultBeginBlock: sdk.ResultBeginBlock{
			Events: sdk.ParseEvents(blockHeader.ResultBeginBlock.Events),
		},
		ResultEndBlock: sdk.ResultEndBlock{
			Events:           sdk.ParseEvents(blockHeader.ResultEndBlock.Events),
			ValidatorUpdates: sdk.ParseValidatorUpdate(blockHeader.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseValidatorSetUpdates(data sdk.EventData) sdk.EventDataValidatorSetUpdates {
	validatorSet := data.(tmtypes.EventDataValidatorSetUpdates)
	return sdk.EventDataValidatorSetUpdates{
		ValidatorUpdates: sdk.ParseValidators(r.cdc, validatorSet.ValidatorUpdates),
	}
}

func getSubscriber() string {
	subscriber := "irita-sdk-go"
	id, err := uuid.NewV1()
	if err == nil {
		subscriber = fmt.Sprintf("%s-%s", subscriber, id.String())
	}
	return subscriber
}
