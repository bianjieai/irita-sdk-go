package integration_test

import (
	"github.com/bianjieai/irita-sdk-go/modules/service"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	"github.com/stretchr/testify/require"
	"time"
)

func (s IntegrationTestSuite) TestService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":"1upoint"}`
	options := `{}`

	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	definition := service.DefineServiceRequest{
		ServiceName:       s.RandStringOfLength(10),
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	result, err := s.Service.DefineService(definition, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), result.Hash)

	defi, err := s.Service.QueryServiceDefinition(definition.ServiceName)
	require.NoError(s.T(), err)
	require.Equal(s.T(), definition.ServiceName, defi.Name)
	require.Equal(s.T(), definition.Description, defi.Description)
	require.EqualValues(s.T(), definition.Tags, defi.Tags)
	require.Equal(s.T(), definition.AuthorDescription, defi.AuthorDescription)
	require.Equal(s.T(), definition.Schemas, defi.Schemas)
	require.Equal(s.T(), s.Account().Address.String(), defi.Author)

	deposit, e := sdk.ParseDecCoins("20000point")
	require.NoError(s.T(), e)
	binding := service.BindServiceRequest{
		ServiceName: definition.ServiceName,
		Deposit:     deposit,
		Pricing:     pricing,
		QoS:         10,
		Options:     options,
	}
	result, err = s.Service.BindService(binding, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), result.Hash)

	bindResp, err := s.Service.QueryServiceBinding(definition.ServiceName, s.Account().Address.String())
	require.NoError(s.T(), err)
	require.Equal(s.T(), binding.ServiceName, bindResp.ServiceName)
	require.Equal(s.T(), s.Account().Address.String(), bindResp.Provider)
	require.Equal(s.T(), binding.Pricing, bindResp.Pricing)

	input := `{"header":{},"body":{"pair":"point-usdt"}}`
	output := `{"header":{},"body":{"last":"1:100"}}`
	testResult := `{"code":200,"message":""}`

	serviceFeeCap, e := sdk.ParseDecCoins("200point")
	require.NoError(s.T(), e)

	invocation := service.InvokeServiceRequest{
		ServiceName:   definition.ServiceName,
		Providers:     []string{s.Account().Address.String()},
		Input:         input,
		ServiceFeeCap: serviceFeeCap,
		Timeout:       10,
		Repeated:      true,
		RepeatedTotal: -1,
	}

	reqCtxID, _, e := s.Service.InvokeService(invocation, baseTx)
	require.NoError(s.T(), e)

	queryRequestContextResp, e := s.Service.QueryRequestContext(reqCtxID)
	require.NoError(s.T(), e)
	require.Equal(s.T(), binding.ServiceName, queryRequestContextResp.ServiceName)

	time.Sleep(time.Second * 3)
	requests, e := s.Service.QueryServiceRequests(bindResp.ServiceName, bindResp.Provider, nil)
	require.NoError(s.T(), e)

	var reqId string
	for _, request := range requests {
		if request.RequestContextID == reqCtxID {
			reqId = request.ID
		}
	}
	require.NotEqual(s.T(), 0, len(reqId))

	request, e := s.Service.QueryServiceRequest(reqId)
	require.NoError(s.T(), e)
	require.Equal(s.T(), bindResp.ServiceName, request.ServiceName)

	responseRequest := service.InvokeServiceResponseRequest{
		RequestId: reqId,
		Output:    output,
		Result:    testResult,
	}
	resultTx, e := s.Service.InvokeServiceResponse(responseRequest, baseTx)
	require.NoError(s.T(), e)
	require.NotEmpty(s.T(), resultTx.Hash)

	response, e := s.Service.QueryServiceResponse(reqId)
	require.NoError(s.T(), e)
	require.NotEmpty(s.T(), response.Result)

	responses, e := s.Service.QueryServiceResponses(reqCtxID, 1, nil)
	require.NoError(s.T(), e)

	var exist bool
	for _, response := range responses {
		if response.RequestContextID == reqCtxID {
			exist = true
			require.NotEmpty(s.T(), response.Output)
		}
	}
	require.True(s.T(), exist)
}
