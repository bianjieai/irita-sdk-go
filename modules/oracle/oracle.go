package oracle

import (
	"context"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type oracleClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return oracleClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (o oracleClient) Name() string {
	return ModuleName
}

func (o oracleClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (o oracleClient) CreateFeed(request CreateFeedRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgCreateFeed{
		FeedName:          request.FeedName,
		AggregateFunc:     request.AggregateFunc,
		ValueJsonPath:     request.ValueJsonPath,
		LatestHistory:     request.LatestHistory,
		Description:       request.Description,
		ServiceName:       request.ServiceName,
		Providers:         request.Providers,
		Input:             request.Input,
		Timeout:           request.Timeout,
		ServiceFeeCap:     request.ServiceFeeCap,
		RepeatedFrequency: request.RepeatedFrequency,
		ResponseThreshold: request.ResponseThreshold,
		Creator:           creator.String(),
	}

	return o.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (o oracleClient) StartFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgStartFeed{
		FeedName: feedName,
		Creator:  creator.String(),
	}

	return o.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (o oracleClient) PauseFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgPauseFeed{
		FeedName: feedName,
		Creator:  creator.String(),
	}

	return o.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (o oracleClient) EditFeedRequest(request EditFeedRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgEditFeed{
		FeedName:          request.FeedName,
		LatestHistory:     request.LatestHistory,
		Description:       request.Description,
		Providers:         request.Providers,
		Timeout:           request.Timeout,
		ServiceFeeCap:     request.ServiceFeeCap,
		RepeatedFrequency: request.RepeatedFrequency,
		ResponseThreshold: request.ResponseThreshold,
		Creator:           creator.String(),
	}

	return o.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (o oracleClient) QueryFeed(feedName string) (QueryFeedResp, sdk.Error) {
	conn, err := o.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryFeedResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Feed(
		context.Background(),
		&QueryFeedRequest{
			FeedName: feedName,
		},
	)
	if err != nil {
		return QueryFeedResp{}, sdk.Wrap(err)
	}

	return resp.Feed.Convert().(QueryFeedResp), nil
}

func (o oracleClient) QueryFeeds(state string) ([]QueryFeedResp, sdk.Error) {
	conn, err := o.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return []QueryFeedResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Feeds(
		context.Background(),
		&QueryFeedsRequest{
			State: state,
		},
	)
	if err != nil {
		return []QueryFeedResp{}, sdk.Wrap(err)
	}

	return Feeds(resp.Feeds).Convert().([]QueryFeedResp), nil
}

func (o oracleClient) QueryFeedValue(feedName string) ([]QueryFeedValueResp, sdk.Error) {
	conn, err := o.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return []QueryFeedValueResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).FeedValue(
		context.Background(),
		&QueryFeedValueRequest{
			FeedName: feedName,
		},
	)
	if err != nil {
		return []QueryFeedValueResp{}, sdk.Wrap(err)
	}

	return FeedValues(resp.FeedValues).Convert().([]QueryFeedValueResp), nil
}
