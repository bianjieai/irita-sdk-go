package oracle

import (
	"time"

	"github.com/bianjieai/irita-sdk-go/v2/modules/service"
	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
	types "github.com/bianjieai/irita-sdk-go/v2/types"
)

var (
	_ Client = oracleClient{}
)

// expose Oracle module api for user
type Client interface {
	sdk.Module

	CreateFeed(request CreateFeedRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	StartFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	PauseFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	EditFeedRequest(request EditFeedRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryFeed(feedName string) (QueryFeedResp, sdk.Error)
	QueryFeeds(state string) ([]QueryFeedResp, sdk.Error)
	QueryFeedValue(feedName string) ([]QueryFeedValueResp, sdk.Error)
}

type CreateFeedRequest struct {
	FeedName          string
	LatestHistory     uint64
	Description       string
	ServiceName       string
	Providers         []string
	Input             string
	Timeout           int64
	ServiceFeeCap     []types.Coin
	RepeatedFrequency uint64
	AggregateFunc     string
	ValueJsonPath     string
	ResponseThreshold uint32
}

type EditFeedRequest struct {
	FeedName          string
	Description       string
	LatestHistory     uint64
	Providers         []string
	Timeout           int64
	ServiceFeeCap     []types.Coin
	RepeatedFrequency uint64
	ResponseThreshold uint32
}

type QueryFeedResp struct {
	Feed              *Feed
	ServiceName       string
	Providers         []string
	Input             string
	Timeout           int64
	ServiceFeeCap     []types.Coin
	RepeatedFrequency uint64
	ResponseThreshold uint32
	State             service.RequestContextState
}

type QueryFeedValueResp struct {
	Data      string
	Timestamp time.Time
}
