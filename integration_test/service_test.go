package integration_test

import (
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/modules/service"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

func (s IntegrationTestSuite) TestService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":"1point"}`

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
	require.Equal(s.T(), s.Account().Address, defi.Author)

	deposit, e := sdk.ParseDecCoins("20000point")
	require.NoError(s.T(), e)
	binding := service.BindServiceRequest{
		ServiceName: definition.ServiceName,
		Deposit:     deposit,
		Pricing:     pricing,
		QoS:         1,
	}
	result, err = s.Service.BindService(binding, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), result.Hash)

	bindResp, err := s.Service.QueryServiceBinding(definition.ServiceName, s.Account().Address)
	require.NoError(s.T(), err)
	require.Equal(s.T(), binding.ServiceName, bindResp.ServiceName)
	require.Equal(s.T(), s.Account().Address, bindResp.Provider)
	require.Equal(s.T(), binding.Pricing, bindResp.Pricing)

	input := `{"pair":"point-usdt"}`
	output := `{"last":"1:100"}`
	testResult := `{"code":200,"message":""}`

	var sub1 sdk.Subscription
	callback := func(reqCtxID, reqID, input string) (string, string) {
		_, err := s.Service.QueryServiceRequest(reqID)
		require.NoError(s.T(), err)
		return output, testResult
	}
	sub1, err = s.Service.SubscribeServiceRequest(definition.ServiceName, callback, baseTx)
	require.NoError(s.T(), err)

	serviceFeeCap, e := sdk.ParseDecCoins("200point")
	require.NoError(s.T(), e)

	invocation := service.InvokeServiceRequest{
		ServiceName:   definition.ServiceName,
		Providers:     []string{s.Account().Address.String()},
		Input:         input,
		ServiceFeeCap: serviceFeeCap,
		Timeout:       3,
		SuperMode:     false,
		Repeated:      true,
		RepeatedTotal: -1,
	}

	var requestContextID string
	var sub2 sdk.Subscription
	var exit = make(chan int)

	requestContextID, err = s.Service.InvokeService(invocation, baseTx)
	require.NoError(s.T(), err)

	sub2, err = s.Service.SubscribeServiceResponse(requestContextID, func(reqCtxID, reqID, responses string) {
		require.Equal(s.T(), reqCtxID, requestContextID)
		require.Equal(s.T(), output, responses)
		request, err := s.Service.QueryServiceRequest(reqID)
		require.NoError(s.T(), err)
		require.Equal(s.T(), reqCtxID, request.RequestContextID)
		require.Equal(s.T(), reqID, request.ID)
		require.Equal(s.T(), input, request.Input)

		exit <- 1
	})
	require.NoError(s.T(), err)

	for {
		select {
		case <-exit:
			require.NoError(s.T(), s.Unsubscribe(sub1))
			require.NoError(s.T(), s.Unsubscribe(sub2))
			goto loop
		case <-time.After(2 * time.Minute):
			require.Panics(s.T(), func() {}, "test service timeout")
		}
	}

loop:
	_, err = s.Service.PauseRequestContext(requestContextID, baseTx)
	require.NoError(s.T(), err)

	_, err = s.Service.StartRequestContext(requestContextID, baseTx)
	require.NoError(s.T(), err)

	request, err := s.Service.QueryRequestContext(requestContextID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), request.ServiceName, invocation.ServiceName)
	require.Equal(s.T(), request.Input, invocation.Input)

	addr, _, err := s.Key.Add(s.RandStringOfLength(30), "1234567890")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), addr)

	_, err = s.Service.SetWithdrawAddress(addr, baseTx)
	require.NoError(s.T(), err)

	fee, err := s.Service.QueryFees(s.Account().Address.String())
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), fee)

	//acc := s.GetRandAccount()

	//TODO
	//rs, err := s.ServiceI.WithdrawEarnedFees(acc.Address.String(), baseTx)
	//require.NoError(s.T(), err)
	//
	//withdrawFee, er := rs.Events.GetValue("transfer", "amount")
	//require.NoError(s.T(), er)
	//require.Equal(s.T(), fee.String(), withdrawFee)
}
