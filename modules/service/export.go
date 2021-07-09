package service

import (
	"time"

	sdk "github.com/irisnet/service-sdk-go/types"
)

// Tx defines a set of transaction interfaces in the service module
type Tx interface {
	DefineService(request DefineServiceRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	BindService(request BindServiceRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	InvokeService(request InvokeServiceRequest, baseTx sdk.BaseTx) (string, sdk.ResultTx, sdk.Error)
	InvokeServiceResponse(request InvokeServiceResponseRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	SetWithdrawAddress(withdrawAddress string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	UpdateServiceBinding(request UpdateServiceBindingRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	DisableServiceBinding(serviceName, provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	EnableServiceBinding(serviceName, provider string, deposit sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	RefundServiceDeposit(serviceName, provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	PauseRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	StartRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	KillRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	UpdateRequestContext(request UpdateRequestContextRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	WithdrawEarnedFees(provider string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	SubscribeServiceRequest(serviceName string, callback RespondCallback, baseTx sdk.BaseTx) (sdk.Subscription, sdk.Error)
	SubscribeServiceResponse(reqCtxID string, callback InvokeCallback) (sdk.Subscription, sdk.Error)
}

// Query defines a set of query interfaces in the service module
type Query interface {
	QueryServiceDefinition(serviceName string) (QueryServiceDefinitionResponse, sdk.Error)
	QueryServiceBinding(serviceName string, provider string) (QueryServiceBindingResponse, sdk.Error)
	QueryServiceBindings(serviceName string) ([]QueryServiceBindingResponse, sdk.Error)
	QueryServiceRequest(requestID string) (QueryServiceRequestResponse, sdk.Error)
	QueryServiceRequests(serviceName string, provider string) ([]QueryServiceRequestResponse, sdk.Error)
	QueryRequestsByReqCtx(requestContextID string, batchCounter uint64) ([]QueryServiceRequestResponse, sdk.Error)
	QueryServiceResponse(requestID string) (QueryServiceResponseResponse, sdk.Error)
	QueryServiceResponses(requestContextID string, batchCounter uint64) ([]QueryServiceResponseResponse, sdk.Error)
	QueryRequestContext(requestContextID string) (QueryRequestContextResp, sdk.Error)
	QueryFees(provider string) (sdk.Coins, sdk.Error)
	QueryParams() (QueryParamsResp, sdk.Error)
}

// Client defines a set of interfaces in the service module
type Client interface {
	sdk.Module
	Tx
	Query
}

// InvokeCallback defines the callback function for service calls
type InvokeCallback func(reqCtxID, reqID, result, responses string)

// RespondCallback defines the callback function of the service response
type RespondCallback func(reqCtxID, reqID, input string) (output string, result string)

// Registry defines a set of service invocation interfaces
type Registry map[string]RespondCallback

// Request defines a request which contains the detailed request data
type QueryServiceRequestResponse struct {
	ID                         string    `json:"id"`
	ServiceName                string    `json:"service_name"`
	Provider                   string    `json:"provider"`
	Consumer                   string    `json:"consumer"`
	Input                      string    `json:"input"`
	ServiceFee                 sdk.Coins `json:"service_fee"`
	RequestHeight              int64     `json:"request_height"`
	ExpirationHeight           int64     `json:"expiration_height"`
	RequestContextID           string    `json:"request_context_id"`
	RequestContextBatchCounter uint64    `json:"request_context_batch_counter"`
}

// Response defines a response
type QueryServiceResponseResponse struct {
	Provider                   string `json:"provider"`
	Consumer                   string `json:"consumer"`
	Output                     string `json:"output"`
	Result                     string `json:"error"`
	RequestContextID           string `json:"request_context_id"`
	RequestContextBatchCounter uint64 `json:"request_context_batch_counter"`
}

// DefineServiceRequest defines the request parameters of the service definition
type DefineServiceRequest struct {
	ServiceName       string   `json:"service_name"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
	AuthorDescription string   `json:"author_description"`
	Schemas           string   `json:"schemas"`
}

// QueryServiceDefinitionResponse represents a service definition
type QueryServiceDefinitionResponse struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
	Author            string   `json:"author"`
	AuthorDescription string   `json:"author_description"`
	Schemas           string   `json:"schemas"`
}

// BindServiceRequest defines the request parameters of the service binding
type BindServiceRequest struct {
	ServiceName string       `json:"service_name"`
	Deposit     sdk.DecCoins `json:"deposit"`
	Pricing     string       `json:"pricing"`
	QoS         uint64       `json:"qos"`
	Options     string       `json:"options"`
	Provider    string       `json:"provider"`
}

// UpdateServiceBindingRequest defines a message to update a service binding
type UpdateServiceBindingRequest struct {
	ServiceName string       `json:"service_name"`
	Deposit     sdk.DecCoins `json:"deposit"`
	Pricing     string       `json:"pricing"`
	QoS         uint64       `json:"qos"`
	Provider    string       `json:"provider"`
}

// QueryServiceBindingResponse defines a struct for service binding
type QueryServiceBindingResponse struct {
	ServiceName  string    `json:"service_name"`
	Provider     string    `json:"provider"`
	Deposit      sdk.Coins `json:"deposit"`
	Pricing      string    `json:"pricing"`
	QoS          uint64    `json:"qos"`
	Options      string    `json:"options"`
	Available    bool      `json:"available"`
	DisabledTime time.Time `json:"disabled_time"`
	Owner        string    `json:"owner"`
}

// InvokeServiceRequest defines the request parameters of the service call
type InvokeServiceRequest struct {
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Input             string       `json:"input"`
	ServiceFeeCap     sdk.DecCoins `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	Repeated          bool         `json:"repeated"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	Callback          InvokeCallback
}

// InvokeServiceResponseRequest defines the request parameters of the service response
type InvokeServiceResponseRequest struct {
	RequestId string `json:"request_id"`
	Output    string `json:"output"`
	Result    string `json:"result"`
}

// UpdateRequestContextRequest defines a message to update a request context
type UpdateRequestContextRequest struct {
	RequestContextID  string       `json:"request_context_id"`
	Providers         []string     `json:"providers"`
	ServiceFeeCap     sdk.DecCoins `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
}

// QueryRequestContextResp defines a context which holds request-related data
type QueryRequestContextResp struct {
	ServiceName        string    `json:"service_name"`
	Providers          []string  `json:"providers"`
	Consumer           string    `json:"consumer"`
	Input              string    `json:"input"`
	ServiceFeeCap      sdk.Coins `json:"service_fee_cap"`
	Timeout            int64     `json:"timeout"`
	Repeated           bool      `json:"repeated"`
	RepeatedFrequency  uint64    `json:"repeated_frequency"`
	RepeatedTotal      int64     `json:"repeated_total"`
	BatchCounter       uint64    `json:"batch_counter"`
	BatchRequestCount  uint32    `json:"batch_request_count"`
	BatchResponseCount uint32    `json:"batch_response_count"`
	BatchState         string    `json:"batch_state"`
	State              string    `json:"state"`
	ResponseThreshold  uint32    `json:"response_threshold"`
	ModuleName         string    `json:"module_name"`
}

type QueryParamsResp struct {
	MaxRequestTimeout    int64         `json:"max_request_timeout"`
	MinDepositMultiple   int64         `json:"min_deposit_multiple"`
	MinDeposit           string        `json:"min_deposit"`
	ServiceFeeTax        string        `json:"service_fee_tax"`
	SlashFraction        string        `json:"slash_fraction"`
	ComplaintRetrospect  time.Duration `json:"complaint_retrospect"`
	ArbitrationTimeLimit time.Duration `json:"arbitration_time_limit"`
	TxSizeLimit          uint64        `json:"tx_size_limit"`
	BaseDenom            string        `json:"base_denom"`
}
