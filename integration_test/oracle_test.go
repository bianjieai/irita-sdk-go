package integration_test

import (
	"github.com/bianjieai/irita-sdk-go/modules/oracle"
	"github.com/bianjieai/irita-sdk-go/modules/service"
	"github.com/stretchr/testify/require"
	//oracle "github.com/bianjieai/irita-sdk-go/modules/oracle"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

func (s IntegrationTestSuite) TestOracle() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}
	serviceName := "test-service"

	//author := val.Address
	//provider := author
	provider := []string{"xiaoming", "xiaohong"}
	//creator := author
	//creator := provider
	feedName := "test-feed"
	aggregateFunc := "avg"
	valueJsonPath := "price"
	latestHistory := 10
	description := "description"
	input := `{"header":{},"body":{}}`
	providers := provider
	timeout := 2
	//serviceDenom := "stake"
	//	serviceFeeCap := fmt.Sprintf("50%s", serviceDenom)
	threshold := 1
	frequency := 12

	createReq := oracle.CreateFeedRequest{
		FeedName:      feedName,
		AggregateFunc: aggregateFunc,
		ValueJsonPath: valueJsonPath,
		LatestHistory: uint64(latestHistory),
		Description:   description,
		Input:         input,
		Timeout:       int64(timeout),
		//ServiceFeeCap : 	serviceFeeCap,
		ResponseThreshold: uint32(threshold),
		RepeatedFrequency: uint64(frequency),
		Providers:         providers,
		ServiceName:       serviceName,
	}

	rs, err := s.Oracle.CreateFeed(createReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	rs, err = s.Oracle.StartFeed(feedName, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	rs, err = s.Oracle.PauseFeed(feedName, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	feedName = "feed-test2"

	editReq := oracle.EditFeedRequest{
		FeedName:      feedName,
		LatestHistory: uint64(11),
		Description:   description,
		Timeout:       int64(timeout),
		//ServiceFeeCap : 	,
		ResponseThreshold: uint32(threshold),
		RepeatedFrequency: uint64(frequency),
		Providers:         provider,
	}

	rs, err = s.Oracle.EditFeedRequest(editReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	state := service.PAUSED
	feedRep, err := s.Oracle.QueryFeed(feedName)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), feedRep)
	require.Equal(s.T(), state, feedRep.State)

	feedsRep, err := s.Oracle.QueryFeeds("PAUSED")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), feedsRep)

	feedValuesRep, err := s.Oracle.QueryFeedValue(feedName)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), feedValuesRep)

	//v, err := s.Node.QueryValidator(validatorID)
	//require.NoError(s.T(), err)
	//
	//vs, err := s.Node.QueryValidators(nil, 0, 0, false)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), vs)

	//updateReq := node.UpdateValidatorRequest{
	//	ID:          validatorID,
	//	Name:        "test2",
	//	Certificate: cert,
	//	Power:       10,
	//	Details:     "this is a updated test",
	//}
	//rs, err = s.Node.UpdateValidator(updateReq, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), rs.Hash)
	//
	//v, err = s.Node.QueryValidator(validatorID)
	//require.NoError(s.T(), err)
	//require.Equal(s.T(), updateReq.Name, v.Name)
	//require.Equal(s.T(), updateReq.Details, v.Details)
	//
	//rs, err = s.Node.RemoveValidator(validatorID, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), rs.Hash)
	//
	//v, err = s.Node.QueryValidator(validatorID)
	//require.Error(s.T(), err)
	//
	//grantNodeReq := node.GrantNodeRequest{
	//	Name:        "test3",
	//	Certificate: cert,
	//	Details:     "this is a grantNode test",
	//}
	//rs, err = s.Node.GrantNode(grantNodeReq, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), rs.Hash)
	//
	//noid, e := rs.Events.GetValue("grant_node", "id")
	//require.NoError(s.T(), e)
	//
	//n, err := s.Node.QueryNode(noid)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), n)
	//
	//ns, err := s.Node.QueryNodes(nil, 0, 0, false)
	//require.NoError(s.T(), err)
	//require.Equal(s.T(), len(ns), 1)
	//
	//rs, err = s.Node.RevokeNode(noid, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), rs.Hash)

}
