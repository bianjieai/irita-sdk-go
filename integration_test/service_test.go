package integration_test

import (
	"fmt"
	"github.com/bianjieai/irita-sdk-go/types/query"
	"github.com/stretchr/testify/require"

	sdk "github.com/bianjieai/irita-sdk-go/types"

	"github.com/bianjieai/irita-sdk-go/modules/service"
)

func (s IntegrationTestSuite) TestTemp() {
	//resp, err := s.Service.QueryServiceDefinition("servname")

	//resp, err := s.Service.QueryServiceBinding("test","iaa1ulkncs27dft8qxve9t9jpcgkka8yeuya28kz2d")

	request := query.PageRequest{
		Offset: 0,
		Limit:  20,
	}
	//resp, err := s.Service.QueryServiceBindings("test", &request)

	//done
	//resp, err := s.Service.QueryRequestContext("A4EB617E01F28910136BD5413E4409E09DFB9FD1A6A989F68EB1C3644BE548A00000000000000000")
	//done
	//resp, err := s.Service.QueryServiceRequests("test", "iaa1ulkncs27dft8qxve9t9jpcgkka8yeuya28kz2d", &request)

	//QueryServiceResponse(requestID string) (QueryServiceResponseResponse, sdk.Error)
	//todo	QueryServiceResponses QueryServiceResponse

	//done
	//resp, err := s.Service.QueryRequestsByReqCtx("A4EB617E01F28910136BD5413E4409E09DFB9FD1A6A989F68EB1C3644BE548A00000000000000000", 1, &request)
	//done
	//resp, err := s.Service.QueryServiceRequest("0347EC493A40C9FBA11872A71BDC0FCA6294DB61C418CD16A74E2B3B96ABAEE8000000000000000000000000000000010000000000000E2D0000")
	//done
	//resp, err := s.Service.QueryFees("iaa1ulkncs27dft8qxve9t9jpcgkka8yeuya28kz2d")
	//done
	resp, err := s.Service.QueryServiceResponses("F035A4B68C6329360B3CDB21DC5859F0A175A68997105A763ACBF3561B59204D0000000000000000", 1, &request)

	//resp, err := s.Service.QueryServiceResponse("6F9E8E507577FF57B43F5DF517ACBD86DF40196E0DC5CF0002D850067B7061B70000000000000000000000000000000100000000000014290000")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func (s IntegrationTestSuite) TestService() {
	//schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	//pricing := `{"price":"1stake"}`
	//options := `{}`
	//
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}
	//
	//definition := service.DefineServiceRequest{
	//	ServiceName:       s.RandStringOfLength(10),
	//	Description:       "this is a test service",
	//	Tags:              nil,
	//	AuthorDescription: "service provider",
	//	Schemas:           schemas,
	//}
	//
	//result, err := s.Service.DefineService(definition, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), result.Hash)
	//
	//defi, err := s.Service.QueryServiceDefinition(definition.ServiceName)
	//require.NoError(s.T(), err)
	//require.Equal(s.T(), definition.ServiceName, defi.Name)
	//require.Equal(s.T(), definition.Description, defi.Description)
	//require.EqualValues(s.T(), definition.Tags, defi.Tags)
	//require.Equal(s.T(), definition.AuthorDescription, defi.AuthorDescription)
	//require.Equal(s.T(), definition.Schemas, defi.Schemas)
	//require.Equal(s.T(), s.Account().Address.String(), defi.Author)
	//
	//deposit, e := sdk.ParseDecCoins("20000stake")
	//require.NoError(s.T(), e)
	//binding := service.BindServiceRequest{
	//	ServiceName: definition.ServiceName,
	//	Deposit:     deposit,
	//	Pricing:     pricing,
	//	QoS:         10,
	//	Options:     options,
	//}
	//result, err = s.Service.BindService(binding, baseTx)
	//require.NoError(s.T(), err)
	//require.NotEmpty(s.T(), result.Hash)
	//
	//bindResp, err := s.Service.QueryServiceBinding(definition.ServiceName, s.Account().Address.String())
	//require.NoError(s.T(), err)
	//require.Equal(s.T(), binding.ServiceName, bindResp.ServiceName)
	//require.Equal(s.T(), s.Account().Address.String(), bindResp.Provider)
	//require.Equal(s.T(), binding.Pricing, bindResp.Pricing)
	//
	//input := `{"header":{},"body":{"pair":"stake-usdt"}}`
	output := `{"header":{},"body":{"last":"1:100"}}`
	testResult := `{"code":200,"message":""}`
	//
	//serviceFeeCap, e := sdk.ParseDecCoins("200stake")
	//require.NoError(s.T(), e)
	//
	//invocation := service.InvokeServiceRequest{
	//	ServiceName:   definition.ServiceName,
	//	Providers:     []string{s.Account().Address.String()},
	//	Input:         input,
	//	ServiceFeeCap: serviceFeeCap,
	//	Timeout:       10,
	//	Repeated:      true,
	//	RepeatedTotal: -1,
	//}
	//
	//reqCtxID, _, e := s.Service.InvokeService(invocation, baseTx)
	//require.NoError(s.T(), e)
	//
	//queryRequestContextResp, e := s.Service.QueryRequestContext(reqCtxID)
	//require.NoError(s.T(), e)
	//require.Equal(s.T(), binding.ServiceName, queryRequestContextResp.ServiceName)

	// TODO
	//requests, e := s.Service.QueryServiceRequests(serviceName, addr, nil)
	//require.NoError(s.T(), e)
	//
	//var reqId string
	//for _, request := range requests {
	//	if request.RequestContextID == reqCtxID {
	//		reqId = request.ID
	//	}
	//}
	//require.NotEqual(s.T(), 0, len(reqId))

	reqCtxID := "09CA1790F330A6678755B1E5023F91FE1AD5E3E493C270C68591C6992DAB36520000000000000000"
	reqId := "67BBF26EE5535A3FFF877BFEE9EEA298D0339945F51FB8CDD9A149AFF6B56E950000000000000000000000000000000B0000000000000BD00000"
	//request, e := s.Service.QueryServiceRequest(reqId)
	//require.NoError(s.T(), e)
	//require.Equal(s.T(), serviceName, request.ServiceName)

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
	fmt.Println(response)

	// TODO
	responses, e := s.Service.QueryServiceResponses(reqCtxID, 0, nil)
	require.NoError(s.T(), e)

	var exist bool
	for _, response := range responses {
		if response.RequestContextID == reqCtxID {
			exist = true
			fmt.Println(response.Output)
		}
	}
	require.True(s.T(), exist)
}
