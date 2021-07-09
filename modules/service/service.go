package service

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/irisnet/service-sdk-go/codec"
	"github.com/irisnet/service-sdk-go/codec/types"
	sdk "github.com/irisnet/service-sdk-go/types"
)

type ServiceClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return ServiceClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (s ServiceClient) Name() string {
	return ModuleName
}

func (s ServiceClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

//DefineService is responsible for creating a new service definition
func (s ServiceClient) DefineService(request DefineServiceRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	author, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgDefineService{
		Name:              request.ServiceName,
		Description:       request.Description,
		Tags:              request.Tags,
		Author:            author.String(),
		AuthorDescription: request.AuthorDescription,
		Schemas:           request.Schemas,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//BindService is responsible for binding a new service definition
func (s ServiceClient) BindService(request BindServiceRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var provider = owner.String()
	if len(request.Provider) > 0 {
		if err := sdk.ValidateAccAddress(request.Provider); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		provider = request.Provider
	}

	amt, err := s.ToMinCoin(request.Deposit...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgBindService{
		ServiceName: request.ServiceName,
		Provider:    provider,
		Deposit:     amt,
		Pricing:     request.Pricing,
		QoS:         request.QoS,
		Options:     request.Options,
		Owner:       owner.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//UpdateServiceBinding updates the specified service binding
func (s ServiceClient) UpdateServiceBinding(request UpdateServiceBindingRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var provider = owner.String()
	if len(request.Provider) > 0 {
		if err := sdk.ValidateAccAddress(request.Provider); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		provider = request.Provider
	}

	amt, err := s.ToMinCoin(request.Deposit...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgUpdateServiceBinding{
		ServiceName: request.ServiceName,
		Provider:    provider,
		Deposit:     amt,
		Pricing:     request.Pricing,
		QoS:         request.QoS,
		Owner:       owner.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// DisableServiceBinding disables the specified service binding
func (s ServiceClient) DisableServiceBinding(serviceName, provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providerAddr = owner.String()
	if len(provider) > 0 {
		if err := sdk.ValidateAccAddress(provider); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		providerAddr = provider
	}

	msg := &MsgDisableServiceBinding{
		ServiceName: serviceName,
		Provider:    providerAddr,
		Owner:       owner.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// EnableServiceBinding enables the specified service binding
func (s ServiceClient) EnableServiceBinding(serviceName, provider string, deposit sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providerAddr = owner.String()
	if len(provider) > 0 {
		if err := sdk.ValidateAccAddress(provider); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		providerAddr = provider
	}

	amt, err := s.ToMinCoin(deposit...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgEnableServiceBinding{
		ServiceName: serviceName,
		Provider:    providerAddr,
		Deposit:     amt,
		Owner:       owner.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//InvokeService is responsible for invoke a new service and callback `handler`
func (s ServiceClient) InvokeService(request InvokeServiceRequest, baseTx sdk.BaseTx) (string, sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return "", sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providers []string
	for _, provider := range request.Providers {
		if err := sdk.ValidateAccAddress(provider); err != nil {
			return "", sdk.ResultTx{}, sdk.Wrap(err)
		}
		providers = append(providers, provider)
	}

	amt, err := s.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return "", sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgCallService{
		ServiceName:       request.ServiceName,
		Providers:         providers,
		Consumer:          consumer.String(),
		Input:             request.Input,
		ServiceFeeCap:     amt,
		Timeout:           request.Timeout,
		Repeated:          request.Repeated,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
	}

	//mode must be set to commit
	baseTx.Mode = sdk.Commit

	result, err := s.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return "", sdk.ResultTx{}, sdk.Wrap(err)
	}

	reqCtxID, e := result.Events.GetValue(sdk.EventTypeCreateContext, attributeKeyRequestContextID)
	if e != nil {
		return reqCtxID, result, sdk.Wrap(e)
	}

	if request.Callback == nil {
		return reqCtxID, result, nil
	}

	_, err = s.SubscribeServiceResponse(reqCtxID, request.Callback)
	return reqCtxID, result, sdk.Wrap(err)
}

func (s ServiceClient) InvokeServiceResponse(req InvokeServiceResponseRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, err
	}

	reqId := req.RequestId
	_, err = s.QueryServiceRequest(reqId)
	if err != nil {
		return sdk.ResultTx{}, err
	}

	msg := &MsgRespondService{
		RequestId: reqId,
		Provider:  provider.String(),
		Result:    req.Result,
		Output:    req.Output,
	}

	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (s ServiceClient) SubscribeServiceResponse(reqCtxID string,
	callback InvokeCallback) (subscription sdk.Subscription, err sdk.Error) {
	if len(reqCtxID) == 0 {
		return subscription, sdk.Wrapf("reqCtxID %s should not be empty", reqCtxID)
	}

	builder := sdk.NewEventQueryBuilder().AddCondition(
		sdk.NewCond(sdk.EventTypeResponseService, attributeKeyRequestContextID).EQ(sdk.EventValue(reqCtxID)),
	)

	return s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
		s.Logger().Debug(
			"consumer received response transaction sent by provider",
			"tx_hash", tx.Hash,
			"height", tx.Height,
			"reqCtxID", reqCtxID,
		)
		for _, msg := range tx.Tx.GetMsgs() {
			msg, ok := msg.(*MsgRespondService)
			if ok {
				reqCtxID2, _, _, _, err := splitRequestID(msg.RequestId)
				if err != nil {
					s.Logger().Error(
						"invalid requestID",
						"requestID", msg.RequestId,
						"errMsg", err.Error(),
					)
					continue
				}
				if reqCtxID2.String() == strings.ToUpper(reqCtxID) {
					callback(reqCtxID, msg.RequestId, msg.Result, msg.Output)
				}
			}
		}
		reqCtx, err := s.QueryRequestContext(reqCtxID)
		if err != nil || reqCtx.State == RequestContextStateToStringMap[COMPLETED] {
			_ = s.Unsubscribe(subscription)
		}
	})
}

// SplitRequestID splits the given contextID to contextID, batchCounter, requestHeight, batchRequestIndex
func splitRequestID(reqID string) (sdk.HexBytes, uint64, int64, int16, error) {
	requestID, err := hex.DecodeString(reqID)
	if err != nil {
		return nil, 0, 0, 0, errors.New("invalid request id")
	}

	if len(requestID) != requestIDLen {
		return nil, 0, 0, 0, errors.New("invalid request id")
	}

	reqCtxID := requestID[0:40]
	batchCounter := binary.BigEndian.Uint64(requestID[40:48])
	requestHeight := int64(binary.BigEndian.Uint64(requestID[48:56]))
	batchRequestIndex := int16(binary.BigEndian.Uint16(requestID[56:]))
	return reqCtxID, batchCounter, requestHeight, batchRequestIndex, nil
}

// SetWithdrawAddress sets a new withdrawal address for the specified service binding
func (s ServiceClient) SetWithdrawAddress(withdrawAddress string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if err := sdk.ValidateAccAddress(withdrawAddress); err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s invalid address", withdrawAddress)
	}
	msg := &MsgSetWithdrawAddress{
		Owner:           owner.String(),
		WithdrawAddress: withdrawAddress,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// RefundServiceDeposit refunds the deposit from the specified service binding
func (s ServiceClient) RefundServiceDeposit(serviceName, provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	if err := sdk.ValidateAccAddress(provider); err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
		Owner:       owner.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// StartRequestContext starts the specified request context
func (s ServiceClient) StartRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgStartRequestContext{
		RequestContextId: requestContextID,
		Consumer:         consumer.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// PauseRequestContext suspends the specified request context
func (s ServiceClient) PauseRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgPauseRequestContext{
		RequestContextId: requestContextID,
		Consumer:         consumer.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// KillRequestContext terminates the specified request context
func (s ServiceClient) KillRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgKillRequestContext{
		RequestContextId: requestContextID,
		Consumer:         consumer.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// UpdateRequestContext updates the specified request context
func (s ServiceClient) UpdateRequestContext(request UpdateRequestContextRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	for _, provider := range request.Providers {
		if err := sdk.ValidateAccAddress(provider); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}

	amt, err := s.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgUpdateRequestContext{
		RequestContextId:  request.RequestContextID,
		Providers:         request.Providers,
		ServiceFeeCap:     amt,
		Timeout:           request.Timeout,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
		Consumer:          consumer.String(),
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// WithdrawEarnedFees withdraws the earned fees to the specified provider
func (s ServiceClient) WithdrawEarnedFees(provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providerAddr = owner.String()
	if len(provider) > 0 {
		if err := sdk.ValidateAccAddress(provider); err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		providerAddr = provider
	}

	msg := &MsgWithdrawEarnedFees{
		Owner:    owner.String(),
		Provider: providerAddr,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//SubscribeSingleServiceRequest is responsible for registering a single service handler
func (s ServiceClient) SubscribeServiceRequest(serviceName string,
	callback RespondCallback,
	baseTx sdk.BaseTx) (subscription sdk.Subscription, err sdk.Error) {
	provider, e := s.QueryAddress(baseTx.From, baseTx.Password)
	if e != nil {
		return sdk.Subscription{}, sdk.Wrap(e)
	}

	builder := sdk.NewEventQueryBuilder().AddCondition(
		sdk.NewCond(eventTypeNewBatchRequestProvider, attributeKeyServiceName).EQ(sdk.EventValue(serviceName)),
	).AddCondition(
		sdk.NewCond(eventTypeNewBatchRequestProvider, attributeKeyProvider).EQ(sdk.EventValue(provider.String())),
	)

	return s.SubscribeNewBlock(builder, func(block sdk.EventDataNewBlock) {
		msgs := s.GenServiceResponseMsgs(block.ResultEndBlock.Events, serviceName, provider, callback)
		if msgs == nil || len(msgs) == 0 {
			s.Logger().Error("no message created",
				"serviceName", serviceName,
				"provider", provider,
			)
		}
		if _, err = s.SendBatch(msgs, baseTx); err != nil {
			s.Logger().Error("provider respond failed", "errMsg", err.Error())
		}
	})
}

// QueryDefinition return a service definition of the specified name
func (s ServiceClient) QueryServiceDefinition(serviceName string) (QueryServiceDefinitionResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryServiceDefinitionResponse{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Definition(
		context.Background(),
		&QueryDefinitionRequest{ServiceName: serviceName},
	)
	if err != nil {
		return QueryServiceDefinitionResponse{}, sdk.Wrap(err)
	}

	return resp.ServiceDefinition.Convert().(QueryServiceDefinitionResponse), nil
}

// QueryBinding return the specified service binding
func (s ServiceClient) QueryServiceBinding(serviceName string, provider string) (QueryServiceBindingResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryServiceBindingResponse{}, sdk.Wrap(err)
	}

	if err := sdk.ValidateAccAddress(provider); err != nil {
		return QueryServiceBindingResponse{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Binding(
		context.Background(),
		&QueryBindingRequest{
			ServiceName: serviceName,
			Provider:    provider,
		},
	)
	if err != nil {
		return QueryServiceBindingResponse{}, sdk.Wrap(err)
	}

	return resp.ServiceBinding.Convert().(QueryServiceBindingResponse), nil
}

// QueryBindings returns all bindings of the specified service
func (s ServiceClient) QueryServiceBindings(serviceName string) ([]QueryServiceBindingResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Bindings(
		context.Background(),
		&QueryBindingsRequest{ServiceName: serviceName},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return serviceBindings(resp.ServiceBindings).Convert().([]QueryServiceBindingResponse), nil
}

// QueryRequest returns  the active request of the specified requestID
func (s ServiceClient) QueryServiceRequest(requestID string) (QueryServiceRequestResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryServiceRequestResponse{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Request(
		context.Background(),
		&QueryRequestRequest{RequestId: requestID},
	)
	if err != nil {
		return QueryServiceRequestResponse{}, sdk.Wrap(err)
	}

	return resp.Request.Convert().(QueryServiceRequestResponse), nil
}

// QueryRequest returns all the active requests of the specified service binding
func (s ServiceClient) QueryServiceRequests(serviceName string, provider string) ([]QueryServiceRequestResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	if err := sdk.ValidateAccAddress(provider); err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Requests(
		context.Background(),
		&QueryRequestsRequest{ServiceName: serviceName, Provider: provider},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return requests(resp.Requests).Convert().([]QueryServiceRequestResponse), nil
}

// QueryRequestsByReqCtx returns all requests of the specified request context ID and batch counter
func (s ServiceClient) QueryRequestsByReqCtx(reqCtxID string, batchCounter uint64) ([]QueryServiceRequestResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).RequestsByReqCtx(
		context.Background(),
		&QueryRequestsByReqCtxRequest{
			RequestContextId: reqCtxID,
			BatchCounter:     batchCounter,
		},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return requests(resp.Requests).Convert().([]QueryServiceRequestResponse), nil
}

// QueryResponse returns a response with the speicified request ID
func (s ServiceClient) QueryServiceResponse(requestID string) (QueryServiceResponseResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryServiceResponseResponse{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Response(
		context.Background(),
		&QueryResponseRequest{RequestId: requestID},
	)
	if err != nil {
		return QueryServiceResponseResponse{}, sdk.Wrap(err)
	}

	return resp.Response.Convert().(QueryServiceResponseResponse), nil
}

// QueryResponses returns all responses of the specified request context and batch counter
func (s ServiceClient) QueryServiceResponses(reqCtxID string, batchCounter uint64) ([]QueryServiceResponseResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Responses(
		context.Background(),
		&QueryResponsesRequest{
			RequestContextId: reqCtxID,
			BatchCounter:     batchCounter,
		},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return responses(resp.Responses).Convert().([]QueryServiceResponseResponse), nil
}

// QueryRequestContext return the specified request context
func (s ServiceClient) QueryRequestContext(reqCtxID string) (QueryRequestContextResp, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryRequestContextResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).RequestContext(
		context.Background(),
		&QueryRequestContextRequest{RequestContextId: reqCtxID},
	)
	if err != nil {
		return QueryRequestContextResp{}, sdk.Wrap(err)
	}

	return resp.RequestContext.Convert().(QueryRequestContextResp), nil
}

//QueryFees return the earned fees for a provider
func (s ServiceClient) QueryFees(provider string) (sdk.Coins, sdk.Error) {
	if err := sdk.ValidateAccAddress(provider); err != nil {
		return nil, sdk.Wrap(err)
	}

	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).EarnedFees(
		context.Background(),
		&QueryEarnedFeesRequest{Provider: provider},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return res.Fees, nil
}

func (s ServiceClient) QueryParams() (QueryParamsResp, sdk.Error) {
	conn, err := s.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Params(
		context.Background(),
		&QueryParamsRequest{},
	)
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	return res.Params.Convert().(QueryParamsResp), nil
}

func (s ServiceClient) GenServiceResponseMsgs(events sdk.StringEvents, serviceName string,
	provider sdk.AccAddress,
	handler RespondCallback) (msgs []sdk.Msg) {

	var ids []string
	for _, e := range events {
		if e.Type != eventTypeNewBatchRequestProvider {
			continue
		}
		attributes := sdk.Attributes(e.Attributes)
		svcName := attributes.GetValue(attributeKeyServiceName)
		prov := attributes.GetValue(attributeKeyProvider)
		if svcName == serviceName && prov == provider.String() {
			reqIDsStr := attributes.GetValue(attributeKeyRequests)
			var idsTemp []string
			if err := json.Unmarshal([]byte(reqIDsStr), &idsTemp); err != nil {
				s.Logger().Error(
					"service request don't exist",
					attributeKeyRequestID, reqIDsStr,
					attributeKeyServiceName, serviceName,
					attributeKeyProvider, provider.String(),
					"errMsg", err.Error(),
				)
				return
			}
			ids = append(ids, idsTemp...)
		}
	}

	for _, reqID := range ids {
		request, err := s.QueryServiceRequest(reqID)
		if err != nil {
			s.Logger().Error(
				"service request don't exist",
				attributeKeyRequestID, reqID,
				attributeKeyServiceName, serviceName,
				attributeKeyProvider, provider.String(),
				"errMsg", err.Error(),
			)
			continue
		}
		//check again
		providerStr := provider.String()
		if providerStr == request.Provider && request.ServiceName == serviceName {
			output, result := handler(request.RequestContextID, reqID, request.Input)
			msgs = append(msgs, &MsgRespondService{
				RequestId: reqID,
				Provider:  providerStr,
				Output:    output,
				Result:    result,
			})
		}
	}
	return msgs
}
