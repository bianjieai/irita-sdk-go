package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type serviceClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) ServiceI {
	return serviceClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (s serviceClient) Name() string {
	return ModuleName
}

func (s serviceClient) RegisterCodec(cdc *codec.LegacyAmino) {
	registerCodec(cdc)
}

func (s serviceClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgDefineService{},
		&MsgBindService{},
		&MsgUpdateServiceBinding{},
		&MsgSetWithdrawAddress{},
		&MsgDisableServiceBinding{},
		&MsgEnableServiceBinding{},
		&MsgRefundServiceDeposit{},
		&MsgCallService{},
		&MsgRespondService{},
		&MsgPauseRequestContext{},
		&MsgStartRequestContext{},
		&MsgKillRequestContext{},
		&MsgUpdateRequestContext{},
		&MsgWithdrawEarnedFees{},
	)
}

//DefineService is responsible for creating a new service definition
func (s serviceClient) DefineService(request DefineServiceRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	author, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgDefineService{
		Name:              request.ServiceName,
		Description:       request.Description,
		Tags:              request.Tags,
		Author:            author,
		AuthorDescription: request.AuthorDescription,
		Schemas:           request.Schemas,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//BindService is responsible for binding a new service definition
func (s serviceClient) BindService(request BindServiceRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var provider = owner
	if len(request.Provider) > 0 {
		provider, err = sdk.AccAddressFromBech32(request.Provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
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
		Owner:       owner,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//UpdateServiceBinding updates the specified service binding
func (s serviceClient) UpdateServiceBinding(request UpdateServiceBindingRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var provider = owner
	if len(request.Provider) > 0 {
		provider, err = sdk.AccAddressFromBech32(request.Provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
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
		Owner:       owner,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// DisableServiceBinding disables the specified service binding
func (s serviceClient) DisableServiceBinding(serviceName, provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providerAddr = owner
	if len(provider) > 0 {
		providerAddr, err = sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}

	msg := &MsgDisableServiceBinding{
		ServiceName: serviceName,
		Provider:    providerAddr,
		Owner:       owner,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// EnableServiceBinding enables the specified service binding
func (s serviceClient) EnableServiceBinding(serviceName, provider string, deposit sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providerAddr = owner
	if len(provider) > 0 {
		providerAddr, err = sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}

	amt, err := s.ToMinCoin(deposit...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgEnableServiceBinding{
		ServiceName: serviceName,
		Provider:    providerAddr,
		Deposit:     amt,
		Owner:       owner,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//InvokeService is responsible for invoke a new service and callback `handler`
func (s serviceClient) InvokeService(request InvokeServiceRequest, baseTx sdk.BaseTx) (string, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		p, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return "", sdk.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	amt, err := s.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	msg := &MsgCallService{
		ServiceName:       request.ServiceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             request.Input,
		ServiceFeeCap:     amt,
		Timeout:           request.Timeout,
		SuperMode:         request.SuperMode,
		Repeated:          request.Repeated,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
	}

	//mode must be set to commit
	baseTx.Mode = sdk.Commit

	result, err := s.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	reqCtxID, e := result.Events.GetValue(sdk.EventTypeMessage, attributeKeyRequestContextID)
	if e != nil {
		return reqCtxID, sdk.Wrap(e)
	}

	if request.Callback == nil {
		return reqCtxID, nil
	}

	_, err = s.SubscribeServiceResponse(reqCtxID, request.Callback)
	return reqCtxID, sdk.Wrap(err)
}

func (s serviceClient) SubscribeServiceResponse(reqCtxID string,
	callback InvokeCallback) (subscription sdk.Subscription, err sdk.Error) {
	if len(reqCtxID) == 0 {
		return subscription, sdk.Wrapf("reqCtxID %s should not be empty", reqCtxID)
	}

	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.NewCond(sdk.EventTypeMessage, attributeKeyRequestContextID).
			EQ(sdk.EventValue(reqCtxID)))

	return s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
		s.Logger().Debug("consumer received response transaction sent by provider",
			"tx_hash", tx.Hash,
			"height", tx.Height,
			"reqCtxID", reqCtxID,
		)
		for _, msg := range tx.Tx.GetMsgs() {
			msg, ok := msg.(*MsgRespondService)
			if ok {
				reqCtxID2, _, _, _, err := splitRequestID(msg.RequestId.String())
				if err != nil {
					s.Logger().Error("invalid requestID",
						"requestID", msg.RequestId.String(),
						"errMsg", err.Error(),
					)
					continue
				}
				if reqCtxID2.String() == strings.ToUpper(reqCtxID) {
					callback(reqCtxID, msg.RequestId.String(), msg.Output)
				}
			}
		}
		reqCtx, err := s.QueryRequestContext(reqCtxID)
		if err != nil || reqCtx.State == RequestContextStateToStringMap[COMPLETED] {
			_ = s.Unsubscribe(subscription)
		}
	})
}

// SetWithdrawAddress sets a new withdrawal address for the specified service binding
func (s serviceClient) SetWithdrawAddress(withdrawAddress string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	withdrawAddr, err := sdk.AccAddressFromBech32(withdrawAddress)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s invalid address", withdrawAddress)
	}
	msg := &MsgSetWithdrawAddress{
		Owner:           owner,
		WithdrawAddress: withdrawAddr,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// RefundServiceDeposit refunds the deposit from the specified service binding
func (s serviceClient) RefundServiceDeposit(serviceName, provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providerAddr sdk.AccAddress
	if len(provider) > 0 {
		providerAddr, err = sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}

	msg := &MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    providerAddr,
		Owner:       owner,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// StartRequestContext starts the specified request context
func (s serviceClient) StartRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgStartRequestContext{
		RequestContextId: sdk.MustHexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// PauseRequestContext suspends the specified request context
func (s serviceClient) PauseRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgPauseRequestContext{
		RequestContextId: sdk.MustHexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// KillRequestContext terminates the specified request context
func (s serviceClient) KillRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}
	msg := &MsgKillRequestContext{
		RequestContextId: sdk.MustHexBytesFrom(requestContextID),
		Consumer:         consumer,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// UpdateRequestContext updates the specified request context
func (s serviceClient) UpdateRequestContext(request UpdateRequestContextRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		p, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		providers = append(providers, p)
	}

	amt, err := s.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgUpdateRequestContext{
		RequestContextId:  sdk.MustHexBytesFrom(request.RequestContextID),
		Providers:         providers,
		ServiceFeeCap:     amt,
		Timeout:           request.Timeout,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
		Consumer:          consumer,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// WithdrawEarnedFees withdraws the earned fees to the specified provider
func (s serviceClient) WithdrawEarnedFees(provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	owner, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providerAddr = owner
	if len(provider) > 0 {
		providerAddr, err = sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
	}

	msg := &MsgWithdrawEarnedFees{
		Owner:    owner,
		Provider: providerAddr,
	}
	return s.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//SubscribeSingleServiceRequest is responsible for registering a single service handler
func (s serviceClient) SubscribeServiceRequest(serviceName string,
	callback RespondCallback,
	baseTx sdk.BaseTx) (subscription sdk.Subscription, err sdk.Error) {
	provider, e := s.QueryAddress(baseTx.From, baseTx.Password)
	if e != nil {
		return sdk.Subscription{}, sdk.Wrap(e)
	}

	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.NewCond(eventTypeNewBatchRequestProvider, attributeKeyProvider).
			EQ(sdk.EventValue(provider.String()))). //TODO
		AddCondition(sdk.NewCond(eventTypeNewBatchRequest, attributeKeyServiceName).
			EQ(sdk.EventValue(serviceName)),
		)
	return s.SubscribeNewBlock(builder, func(block sdk.EventDataNewBlock) {
		msgs := s.GenServiceResponseMsgs(block.ResultEndBlock.Events, serviceName, provider, callback)
		if _, err = s.SendBatch(msgs, baseTx); err != nil {
			s.Logger().Error("provider respond failed",
				"errMsg", err.Error(),
			)
		}
	})
}

// QueryDefinition return a service definition of the specified name
func (s serviceClient) QueryServiceDefinition(serviceName string) (QueryServiceDefinitionResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryServiceDefinitionResponse{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Definition(context.Background(), &QueryDefinitionRequest{
		ServiceName: serviceName,
	})
	if err != nil {
		return QueryServiceDefinitionResponse{}, sdk.Wrap(err)
	}

	return resp.ServiceDefinition.Convert().(QueryServiceDefinitionResponse), nil
}

// QueryBinding return the specified service binding
func (s serviceClient) QueryServiceBinding(serviceName string, provider sdk.AccAddress) (QueryServiceBindingResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryServiceBindingResponse{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Binding(context.Background(), &QueryBindingRequest{
		ServiceName: serviceName,
		Provider:    provider,
	})
	if err != nil {
		return QueryServiceBindingResponse{}, sdk.Wrap(err)
	}

	return resp.ServiceBinding.Convert().(QueryServiceBindingResponse), nil
}

// QueryBindings returns all bindings of the specified service
func (s serviceClient) QueryServiceBindings(serviceName string) ([]QueryServiceBindingResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Bindings(context.Background(), &QueryBindingsRequest{
		ServiceName: serviceName,
	})
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return serviceBindings(resp.ServiceBindings).Convert().([]QueryServiceBindingResponse), nil
}

// QueryRequest returns  the active request of the specified requestID
func (s serviceClient) QueryServiceRequest(requestID string) (QueryServiceRequestResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryServiceRequestResponse{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Request(context.Background(), &QueryRequestRequest{
		RequestId: sdk.MustHexBytesFrom(requestID),
	})
	if err != nil {
		return QueryServiceRequestResponse{}, sdk.Wrap(err)
	}

	return resp.Request.Convert().(QueryServiceRequestResponse), nil
}

// QueryRequest returns all the active requests of the specified service binding
func (s serviceClient) QueryServiceRequests(serviceName string, provider sdk.AccAddress) ([]QueryServiceRequestResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Requests(context.Background(), &QueryRequestsRequest{
		ServiceName: serviceName,
		Provider:    provider,
	})
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return requests(resp.Requests).Convert().([]QueryServiceRequestResponse), nil
}

// QueryRequestsByReqCtx returns all requests of the specified request context ID and batch counter
func (s serviceClient) QueryRequestsByReqCtx(reqCtxID string, batchCounter uint64) ([]QueryServiceRequestResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).RequestsByReqCtx(context.Background(), &QueryRequestsByReqCtxRequest{
		RequestContextId: sdk.MustHexBytesFrom(reqCtxID),
		BatchCounter:     batchCounter,
	})
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return requests(resp.Requests).Convert().([]QueryServiceRequestResponse), nil
}

// QueryResponse returns a response with the speicified request ID
func (s serviceClient) QueryServiceResponse(requestID string) (QueryServiceResponseResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryServiceResponseResponse{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Response(context.Background(), &QueryResponseRequest{
		RequestId: sdk.MustHexBytesFrom(requestID),
	})

	if err != nil {
		return QueryServiceResponseResponse{}, sdk.Wrap(err)
	}

	return resp.Response.Convert().(QueryServiceResponseResponse), nil
}

// QueryResponses returns all responses of the specified request context and batch counter
func (s serviceClient) QueryServiceResponses(reqCtxID string, batchCounter uint64) ([]QueryServiceResponseResponse, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Responses(context.Background(), &QueryResponsesRequest{
		RequestContextId: sdk.MustHexBytesFrom(reqCtxID),
		BatchCounter:     batchCounter,
	})
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	return responses(resp.Responses).Convert().([]QueryServiceResponseResponse), nil
}

// QueryRequestContext return the specified request context
func (s serviceClient) QueryRequestContext(reqCtxID string) (QueryRequestContextResp, sdk.Error) {
	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return QueryRequestContextResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).RequestContext(context.Background(), &QueryRequestContextRequest{
		RequestContextId: sdk.MustHexBytesFrom(reqCtxID),
	})
	if err != nil {
		return QueryRequestContextResp{}, sdk.Wrap(err)
	}

	return resp.RequestContext.Convert().(QueryRequestContextResp), nil
}

//QueryFees return the earned fees for a provider
func (s serviceClient) QueryFees(provider string) (sdk.Coins, sdk.Error) {
	address, addressErr := sdk.AccAddressFromBech32(provider)
	if addressErr != nil {
		return nil, sdk.Wrap(addressErr)
	}

	conn, err := s.GenConn()
	defer conn.Close()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).EarnedFees(context.Background(), &QueryEarnedFeesRequest{
		Provider: address,
	})
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return res.Fees, nil
}

func (s serviceClient) QueryParams() (QueryParamsResp, sdk.Error) {
	var param Params
	if err := s.BaseClient.QueryParams(ModuleName, &param); err != nil {
		return QueryParamsResp{}, err
	}
	return param.Convert().(QueryParamsResp), nil
}

func (s serviceClient) GenServiceResponseMsgs(events sdk.StringEvents, serviceName string,
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
				s.Logger().Error("service request don't exist",
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
			s.Logger().Error("service request don't exist",
				attributeKeyRequestID, reqID,
				attributeKeyServiceName, serviceName,
				attributeKeyProvider, provider.String(),
				"errMsg", err.Error(),
			)
			continue
		}
		//check again
		if provider.Equals(request.Provider) && request.ServiceName == serviceName {
			output, result := handler(request.RequestContextID, reqID, request.Input)
			msgs = append(msgs, &MsgRespondService{
				RequestId: sdk.MustHexBytesFrom(reqID),
				Provider:  provider,
				Output:    output,
				Result:    result,
			})
		}
	}
	return msgs
}
