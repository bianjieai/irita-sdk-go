package service

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	// ModuleName define module name
	ModuleName = "service"

	eventTypeNewBatchRequest         = "new_batch_request"
	eventTypeNewBatchRequestProvider = "new_batch_request_provider"
	attributeKeyRequests             = "requests"
	attributeKeyRequestID            = "request_id"
	attributeKeyRequestContextID     = "request_context_id"
	attributeKeyServiceName          = "service_name"
	attributeKeyProvider             = "provider"

	requestIDLen = 58
	contextIDLen = 40
)

var (
	_ sdk.Msg = MsgDefineService{}
	_ sdk.Msg = MsgBindService{}
	_ sdk.Msg = MsgUpdateServiceBinding{}
	_ sdk.Msg = MsgSetWithdrawAddress{}
	_ sdk.Msg = MsgDisableServiceBinding{}
	_ sdk.Msg = MsgEnableServiceBinding{}
	_ sdk.Msg = MsgRefundServiceDeposit{}
	_ sdk.Msg = MsgCallService{}
	_ sdk.Msg = MsgRespondService{}
	_ sdk.Msg = MsgPauseRequestContext{}
	_ sdk.Msg = MsgStartRequestContext{}
	_ sdk.Msg = MsgKillRequestContext{}
	_ sdk.Msg = MsgUpdateRequestContext{}
	_ sdk.Msg = MsgWithdrawEarnedFees{}

	amino = codec.New()

	// ModuleCdc references the global service module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to service and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())

	RequestContextStateToStringMap = map[RequestContextState]string{
		RUNNING:   "running",
		PAUSED:    "paused",
		COMPLETED: "completed",
	}
	StringToRequestContextStateMap = map[string]RequestContextState{
		"running":   RUNNING,
		"paused":    PAUSED,
		"completed": COMPLETED,
	}

	RequestContextBatchStateToStringMap = map[RequestContextBatchState]string{
		BATCHRUNNING:   "running",
		BATCHCOMPLETED: "completed",
	}
	StringToRequestContextBatchStateMap = map[string]RequestContextBatchState{
		"running":   BATCHRUNNING,
		"completed": BATCHCOMPLETED,
	}
)

func init() {
	registerCodec(amino)
}

func (msg MsgDefineService) Route() string { return ModuleName }

func (msg MsgDefineService) Type() string {
	return "define_service"
}

func (msg MsgDefineService) ValidateBasic() error {
	if len(msg.Author) == 0 {
		return errors.New("author missing")
	}

	if len(msg.Name) == 0 {
		return errors.New("author missing")
	}

	if len(msg.Schemas) == 0 {
		return errors.New("schemas missing")
	}

	return nil
}

func (msg MsgDefineService) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}

	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgDefineService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

func (msg MsgBindService) Type() string {
	return "bind_service"
}

func (msg MsgBindService) Route() string { return ModuleName }

func (msg MsgBindService) ValidateBasic() error {
	if len(msg.Owner) == 0 {
		return errors.New("owner missing")
	}

	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(msg.Pricing) == 0 {
		return errors.New("pricing missing")
	}
	return nil
}

func (msg MsgBindService) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgBindService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgCallService) Route() string { return ModuleName }

func (msg MsgCallService) Type() string {
	return "request_service"
}

func (msg MsgCallService) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}
	if len(msg.Providers) == 0 {
		return errors.New("providers missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(msg.Input) == 0 {
		return errors.New("input missing")
	}
	return nil
}

func (msg MsgCallService) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgCallService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

func (msg MsgRespondService) Route() string { return ModuleName }

func (msg MsgRespondService) Type() string {
	return "respond_service"
}

func (msg MsgRespondService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Result) == 0 {
		return errors.New("result missing")
	}

	if len(msg.Output) > 0 {
		if !json2.Valid([]byte(msg.Output)) {
			return errors.New("output is not valid JSON")
		}
	}

	return nil
}

func (msg MsgRespondService) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgRespondService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

func (msg MsgUpdateServiceBinding) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgUpdateServiceBinding) Type() string { return "update_service_binding" }

// GetSignBytes implements Msg.
func (msg MsgUpdateServiceBinding) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateServiceBinding) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Owner) == 0 {
		return errors.New("owner missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if !msg.Deposit.Empty() {
		return fmt.Errorf("invalid deposit: %s", msg.Deposit)
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//______________________________________________________________________

func (msg MsgSetWithdrawAddress) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgSetWithdrawAddress) Type() string { return "set_withdraw_address" }

// GetSignBytes implements Msg.
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSetWithdrawAddress) ValidateBasic() error {
	if len(msg.Owner) == 0 {
		return errors.New("owner missing")
	}

	if len(msg.WithdrawAddress) == 0 {
		return errors.New("withdrawal address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//______________________________________________________________________

func (msg MsgDisableServiceBinding) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgDisableServiceBinding) Type() string { return "disable_service_binding" }

// GetSignBytes implements Msg.
func (msg MsgDisableServiceBinding) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgDisableServiceBinding) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Owner) == 0 {
		return errors.New("owner missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgDisableServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//______________________________________________________________________

func (msg MsgEnableServiceBinding) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgEnableServiceBinding) Type() string { return "enable_service_binding" }

// GetSignBytes implements Msg.
func (msg MsgEnableServiceBinding) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgEnableServiceBinding) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Owner) == 0 {
		return errors.New("owner missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if !msg.Deposit.Empty() {
		return fmt.Errorf("invalid deposit: %s", msg.Deposit)
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgEnableServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//______________________________________________________________________

func (msg MsgRefundServiceDeposit) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgRefundServiceDeposit) Type() string { return "refund_service_deposit" }

// GetSignBytes implements Msg.
func (msg MsgRefundServiceDeposit) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRefundServiceDeposit) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Owner) == 0 {
		return errors.New("owner missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgRefundServiceDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//______________________________________________________________________

func (msg MsgPauseRequestContext) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgPauseRequestContext) Type() string { return "pause_request_context" }

// GetSignBytes implements Msg.
func (msg MsgPauseRequestContext) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgPauseRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgPauseRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

func (msg MsgStartRequestContext) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgStartRequestContext) Type() string { return "start_request_context" }

// GetSignBytes implements Msg.
func (msg MsgStartRequestContext) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgStartRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgStartRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

func (msg MsgKillRequestContext) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgKillRequestContext) Type() string { return "kill_request_context" }

// GetSignBytes implements Msg.
func (msg MsgKillRequestContext) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgKillRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgKillRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

func (msg MsgUpdateRequestContext) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgUpdateRequestContext) Type() string { return "update_request_context" }

// GetSignBytes implements Msg.
func (msg MsgUpdateRequestContext) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

func (msg MsgWithdrawEarnedFees) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgWithdrawEarnedFees) Type() string { return "withdraw_earned_fees" }

// GetSignBytes implements Msg.
func (msg MsgWithdrawEarnedFees) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawEarnedFees) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Owner) == 0 {
		return errors.New("owner missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawEarnedFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//==========================================for QueryWithResponse==========================================

func (r ServiceDefinition) Convert() interface{} {
	return QueryServiceDefinitionResponse{
		Name:              r.Name,
		Description:       r.Description,
		Tags:              r.Tags,
		Author:            r.Author,
		AuthorDescription: r.AuthorDescription,
		Schemas:           r.Schemas,
	}
}

func (b ServiceBinding) Convert() interface{} {
	return QueryServiceBindingResponse{
		ServiceName:  b.ServiceName,
		Provider:     b.Provider,
		Deposit:      b.Deposit,
		Pricing:      b.Pricing,
		QoS:          b.QoS,
		Available:    b.Available,
		DisabledTime: b.DisabledTime,
		Owner:        b.Owner,
	}

}

type serviceBindings []ServiceBinding

func (bs serviceBindings) Convert() interface{} {
	bindings := make([]QueryServiceBindingResponse, len(bs))
	for i, binding := range bs {
		bindings[i] = binding.Convert().(QueryServiceBindingResponse)
	}
	return bindings
}

func (r Request) Empty() bool {
	return len(r.ServiceName) == 0
}

func (r Request) Convert() interface{} {
	return QueryServiceRequestResponse{
		ID:                         r.ID.String(),
		ServiceName:                r.ServiceName,
		Provider:                   r.Provider,
		Consumer:                   r.Consumer,
		Input:                      r.Input,
		ServiceFee:                 r.ServiceFee,
		SuperMode:                  r.SuperMode,
		RequestHeight:              r.RequestHeight,
		ExpirationHeight:           r.ExpirationHeight,
		RequestContextID:           r.RequestContextID.String(),
		RequestContextBatchCounter: r.RequestContextBatchCounter,
	}
}

type requests []Request

func (rs requests) Convert() interface{} {
	requests := make([]QueryServiceRequestResponse, len(rs))
	for i, request := range rs {
		requests[i] = request.Convert().(QueryServiceRequestResponse)
	}
	return requests
}

func (r Response) Empty() bool {
	return len(r.Provider) == 0
}

func (r Response) Convert() interface{} {
	return QueryServiceResponseResponse{
		Provider:                   r.Provider,
		Consumer:                   r.Consumer,
		Output:                     r.Output,
		Result:                     r.Result,
		RequestContextID:           r.RequestContextID.String(),
		RequestContextBatchCounter: r.RequestContextBatchCounter,
	}
}

type responses []Response

func (rs responses) Convert() interface{} {
	responses := make([]QueryServiceResponseResponse, len(rs))
	for i, response := range rs {
		responses[i] = response.Convert().(QueryServiceResponseResponse)
	}
	return responses
}

func RequestContextStateFromString(str string) (RequestContextState, error) {
	if state, ok := StringToRequestContextStateMap[strings.ToLower(str)]; ok {
		return state, nil
	}
	return RequestContextState(0xff), fmt.Errorf("'%s' is not a valid request context state", str)
}

func (state RequestContextState) String() string {
	return RequestContextStateToStringMap[state]
}

// MarshalJSON returns the JSON representation
func (state RequestContextState) MarshalJSON() ([]byte, error) {
	return json2.Marshal(state.String())
}

func RequestContextBatchStateFromString(str string) (RequestContextBatchState, error) {
	if state, ok := StringToRequestContextBatchStateMap[strings.ToLower(str)]; ok {
		return state, nil
	}
	return RequestContextBatchState(0xff), fmt.Errorf("'%s' is not a valid request context batch state", str)
}

func (state RequestContextBatchState) String() string {
	return RequestContextBatchStateToStringMap[state]
}

// MarshalJSON returns the JSON representation
func (state RequestContextBatchState) MarshalJSON() ([]byte, error) {
	return json2.Marshal(state.String())
}

// UnmarshalJSON unmarshals raw JSON bytes into a RequestContextBatchState
func (state *RequestContextBatchState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json2.Unmarshal(data, &s); err != nil {
		return nil
	}

	bz, err := RequestContextBatchStateFromString(s)
	if err != nil {
		return err
	}

	*state = bz
	return nil
}

// UnmarshalJSON unmarshals raw JSON bytes into a RequestContextState.
func (state *RequestContextState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json2.Unmarshal(data, &s); err != nil {
		return nil
	}

	bz, err := RequestContextStateFromString(s)
	if err != nil {
		return err
	}

	*state = bz
	return nil
}

// Empty returns true if empty
func (r RequestContext) Empty() bool {
	return len(r.ServiceName) == 0
}

func (r RequestContext) Convert() interface{} {
	return QueryRequestContextResponse{
		ServiceName:        r.ServiceName,
		Providers:          r.Providers,
		Consumer:           r.Consumer,
		Input:              r.Input,
		ServiceFeeCap:      r.ServiceFeeCap,
		Timeout:            r.Timeout,
		SuperMode:          r.SuperMode,
		Repeated:           r.Repeated,
		RepeatedFrequency:  r.RepeatedFrequency,
		RepeatedTotal:      r.RepeatedTotal,
		BatchCounter:       r.BatchCounter,
		BatchRequestCount:  r.BatchRequestCount,
		BatchResponseCount: r.BatchResponseCount,
		BatchState:         r.BatchState.String(),
		State:              r.State.String(),
		ResponseThreshold:  r.ResponseThreshold,
		ModuleName:         r.ModuleName,
	}
}

func (p Params) Convert() interface{} {
	return QueryParamsResponse{
		MaxRequestTimeout:    p.MaxRequestTimeout,
		MinDepositMultiple:   p.MinDepositMultiple,
		MinDeposit:           p.MinDeposit.String(),
		ServiceFeeTax:        p.ServiceFeeTax.String(),
		SlashFraction:        p.SlashFraction.String(),
		ComplaintRetrospect:  p.ComplaintRetrospect,
		ArbitrationTimeLimit: p.ArbitrationTimeLimit,
		TxSizeLimit:          p.TxSizeLimit,
		BaseDenom:            p.BaseDenom,
	}
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDefineService{}, "irismod/service/MsgDefineService", nil)
	cdc.RegisterConcrete(MsgBindService{}, "irismod/service/MsgBindService", nil)
	cdc.RegisterConcrete(MsgUpdateServiceBinding{}, "irismod/service/MsgUpdateServiceBinding", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "irismod/service/MsgSetWithdrawAddress", nil)
	cdc.RegisterConcrete(MsgDisableServiceBinding{}, "irismod/service/MsgDisableServiceBinding", nil)
	cdc.RegisterConcrete(MsgEnableServiceBinding{}, "irismod/service/MsgEnableServiceBinding", nil)
	cdc.RegisterConcrete(MsgRefundServiceDeposit{}, "irismod/service/MsgRefundServiceDeposit", nil)
	cdc.RegisterConcrete(MsgCallService{}, "irismod/service/MsgCallService", nil)
	cdc.RegisterConcrete(MsgRespondService{}, "irismod/service/MsgRespondService", nil)
	cdc.RegisterConcrete(MsgPauseRequestContext{}, "irismod/service/MsgPauseRequestContext", nil)
	cdc.RegisterConcrete(MsgStartRequestContext{}, "irismod/service/MsgStartRequestContext", nil)
	cdc.RegisterConcrete(MsgKillRequestContext{}, "irismod/service/MsgKillRequestContext", nil)
	cdc.RegisterConcrete(MsgUpdateRequestContext{}, "irismod/service/MsgUpdateRequestContext", nil)
	cdc.RegisterConcrete(MsgWithdrawEarnedFees{}, "irismod/service/MsgWithdrawEarnedFees", nil)

	cdc.RegisterConcrete(ServiceDefinition{}, "irismod/service/ServiceDefinition", nil)
	cdc.RegisterConcrete(ServiceBinding{}, "irismod/service/ServiceBinding", nil)
	cdc.RegisterConcrete(RequestContext{}, "irismod/service/RequestContext", nil)
	cdc.RegisterConcrete(Request{}, "irismod/service/Request", nil)
	cdc.RegisterConcrete(Response{}, "irismod/service/Response", nil)

	cdc.RegisterConcrete(&Params{}, "irismod/service/Params", nil)
}
