package integration_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/modules/oracle"
	"github.com/bianjieai/irita-sdk-go/modules/service"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

var serviceName = generateServiceName()

func (s *IntegrationTestSuite) SetupService(ch chan<- int) {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":"1upoint"}`
	output := `{"header":{},"body":{"last":"100"}}`
	testResult := `{"code":200,"message":""}`

	coin, _ := sdk.ParseDecCoins("4point")
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Fee:      coin,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	definition := service.DefineServiceRequest{
		ServiceName:       serviceName,
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	_, err := s.Service.DefineService(definition, baseTx)
	require.NoError(s.T(), err)
	deposit, _ := sdk.ParseDecCoins("6000point")
	binding := service.BindServiceRequest{
		ServiceName: definition.ServiceName,
		Deposit:     deposit,
		Pricing:     pricing,
		QoS:         10,
		Options:     `{}`,
	}
	_, err = s.Service.BindService(binding, baseTx)
	require.NoError(s.T(), err)

	_, err = s.Service.SubscribeServiceRequest(
		definition.ServiceName,
		func(reqCtxID, reqID, input string) (string, string) {
			s.Logger().Info("Service received request", "input", input, "reqCtxID", reqCtxID, "reqID", reqID, "output", output)
			ch <- 1
			return output, testResult
		}, baseTx)

	require.NoError(s.T(), err)
}
func (s IntegrationTestSuite) TestOracle() {
	var ch = make(chan int)
	s.SetupService(ch)

	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}
	input := `{"header":{},"body":{"pair":"iris-usdt"}}`
	feedName := generateFeedName(serviceName)
	serviceFeeCap, _ := sdk.ParseCoins("1000upoint")

	sender := s.rootAccount.Address
	createReq := oracle.CreateFeedRequest{
		FeedName:          feedName,
		LatestHistory:     5,
		Description:       "fetch USDT-CNY ",
		ServiceName:       serviceName,
		Providers:         []string{sender.String()},
		Input:             input,
		Timeout:           50,
		ServiceFeeCap:     serviceFeeCap,
		RepeatedFrequency: 50,
		AggregateFunc:     "avg",
		ValueJsonPath:     "last",
		ResponseThreshold: 1,
	}

	cfrs, err := s.Oracle.CreateFeed(createReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), cfrs.Hash)

	sfrs, err := s.Oracle.StartFeed(feedName, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), sfrs.Hash)

	select {
	case <-ch:

		time.Sleep(2 * time.Second)

		feedValuesRep, err := s.Oracle.QueryFeedValue(feedName)
		require.NoError(s.T(), err)
		s.Logger().Info("Query feed value", "feedName", feedName, "result", feedValuesRep)

		editReq := oracle.EditFeedRequest{
			FeedName:          feedName,
			LatestHistory:     5,
			Description:       "fetch USDT-CNY ",
			Timeout:           3,
			ServiceFeeCap:     serviceFeeCap,
			ResponseThreshold: 1,
			RepeatedFrequency: 5,
			Providers:         []string{sender.String()},
		}

		efrs, err := s.Oracle.EditFeedRequest(editReq, baseTx)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), efrs.Hash)

		pfrs, err := s.Oracle.PauseFeed(feedName, baseTx)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), pfrs.Hash)

		feedRep, err := s.Oracle.QueryFeed(feedName)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), feedRep)

		feedsRep, err := s.Oracle.QueryFeeds("PAUSED")
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), feedsRep)
		require.Equal(s.T(), service.PAUSED, feedRep.State)
	}
}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}

func generateFeedName(serviceName string) string {
	return fmt.Sprintf("feed-%s", serviceName)
}
