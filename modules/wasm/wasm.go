package wasm

import (
	context "context"
	"io/ioutil"

	"github.com/spf13/cast"

	"github.com/bianjieai/irita-sdk-go/codec"
	codectypes "github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

var (
	_ Client = wasmClient{}
)

type wasmClient struct {
	sdk.BaseClient
}

//NewClient return the instance of wasm client
func NewClient(bc sdk.BaseClient) Client {
	return wasmClient{
		BaseClient: bc,
	}
}

//Name return the module name
func (wasm wasmClient) Name() string {
	return ModuleName
}

//RegisterCodec register the module struce to amion
func (wasm wasmClient) RegisterCodec(cdc *codec.LegacyAmino) {}

//RegisterInterfaceTypes register the module struce to InterfaceRegistry
func (wasm wasmClient) RegisterInterfaceTypes(registry codectypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

//Store upload the contract to cschain and return the codeId for user
func (wasm wasmClient) Store(request StoreRequest, config sdk.BaseTx) (string, error) {
	sender, err := wasm.QueryAddress(config.From, config.Password)
	if err != nil {
		return "", err
	}

	var byteCode []byte
	if len(request.WASMByteCode) > 0 {
		byteCode = request.WASMByteCode
	} else if len(request.WASMFile) > 0 {
		bz, err := ioutil.ReadFile(request.WASMFile)
		if err != nil {
			return "", err
		}
		byteCode = bz
	}
	msg := &MsgStoreCode{
		Sender:                sender.String(),
		WASMByteCode:          byteCode,
		Source:                request.Source,
		Builder:               request.Builder,
		InstantiatePermission: request.permission,
	}
	result, err := wasm.BuildAndSend(sdk.Msgs{msg}, config)
	if err != nil {
		return "", err
	}
	return result.Events.GetValue(sdk.EventTypeMessage, "code_id")
}

//Instantiate instantiate the contract state
func (wasm wasmClient) Instantiate(request InstantiateRequest, config sdk.BaseTx) (string, error) {
	sender, err := wasm.QueryAddress(config.From, config.Password)
	if err != nil {
		return "", err
	}

	msg := &MsgInstantiateContract{
		Sender:    sender.String(),
		Admin:     request.Admin,
		CodeID:    cast.ToUint64(request.CodeID),
		Label:     request.Label,
		InitMsg:   request.InitMsg.MarshallJson(),
		InitFunds: request.InitFunds,
	}
	result, err := wasm.BuildAndSend(sdk.Msgs{msg}, config)
	if err != nil {
		return "", err
	}
	return result.Events.GetValue(sdk.EventTypeMessage, "contract_address")
}

//Execute execute the contract method
func (wasm wasmClient) Execute(contractAddress string,
	abi *ContractABI,
	sentFunds sdk.Coins,
	config sdk.BaseTx) (sdk.ResultTx, error) {
	sender, err := wasm.QueryAddress(config.From, config.Password)
	if err != nil {
		return sdk.ResultTx{}, err
	}

	msgBytes, er := abi.Build()
	if er != nil {
		return sdk.ResultTx{}, er
	}

	msg := &MsgExecuteContract{
		Sender:    sender.String(),
		Contract:  contractAddress,
		SentFunds: sentFunds,
		Msg:       msgBytes,
	}
	return wasm.BuildAndSend(sdk.Msgs{msg}, config)
}

//Execute execute the contract method
func (wasm wasmClient) Migrate(contractAddress string,
	newCodeID string,
	msgByte []byte,
	config sdk.BaseTx) (sdk.ResultTx, error) {
	sender, err := wasm.QueryAddress(config.From, config.Password)
	if err != nil {
		return sdk.ResultTx{}, err
	}

	msg := &MsgMigrateContract{
		Sender:     sender.String(),
		Contract:   contractAddress,
		CodeID:     cast.ToUint64(newCodeID),
		MigrateMsg: msgByte,
	}
	return wasm.BuildAndSend(sdk.Msgs{msg}, config)
}

//QueryContractInfo return the contract information
func (wasm wasmClient) QueryContractInfo(address string) (*ContractInfo, error) {
	conn, err := wasm.GenConn()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	req := &QueryContractInfoRequest{
		Address: address,
	}

	res, err := NewQueryClient(conn).ContractInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res.ContractInfo, nil
}

//ExportContractState export all state data of the contract
func (wasm wasmClient) ExportContractState(address string) (map[string][]byte, error) {
	conn, err := wasm.GenConn()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	req := &QueryAllContractStateRequest{
		Address: address,
	}

	res, err := NewQueryClient(conn).AllContractState(context.Background(), req)
	if err != nil {
		return nil, err
	}

	var states = make(map[string][]byte, len(res.Models))
	for _, m := range res.Models {
		states[string(m.Key.Bytes())] = m.Value
	}

	return states, nil
}

//QueryContract execute contract's query method and return the result
func (wasm wasmClient) QueryContract(address string, abi *ContractABI) ([]byte, error) {
	conn, err := wasm.GenConn()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	msgBytes, err := abi.Build()
	if err != nil {
		return nil, err
	}

	req := &QuerySmartContractStateRequest{
		Address:   address,
		QueryData: msgBytes,
	}

	res, err := NewQueryClient(conn).SmartContractState(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}
