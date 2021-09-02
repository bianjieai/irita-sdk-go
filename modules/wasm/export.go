package wasm

import (
	"encoding/json"
	"errors"

	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
)

// Client define a group of interface for wasm module
type Client interface {
	sdk.Module

	Store(request StoreRequest, config sdk.BaseTx) (string, error)
	Instantiate(request InstantiateRequest, config sdk.BaseTx) (string, error)
	Execute(contractAddress string,
		abi *ContractABI,
		SentFunds sdk.Coins,
		config sdk.BaseTx) (sdk.ResultTx, error)
	Migrate(contractAddress string,
		newCodeID string,
		msgByte []byte,
		config sdk.BaseTx) (sdk.ResultTx, error)

	QueryContractInfo(address string) (*ContractInfo, error)
	QueryContract(address string, abi *ContractABI) ([]byte, error)
	ExportContractState(address string) (map[string][]byte, error)
}

//StoreRequest define a struct for Store method
type StoreRequest struct {
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode []byte
	// WASMFile can be raw or gzip file
	WASMFile string
	// Source is a valid absolute HTTPS URI to the contract's source code, optional
	Source string
	// Builder is a valid docker image name with tag, optional
	Builder string
	// InstantiatePermission access control to apply on contract creation, optional
	permission *AccessConfig
}

//InstantiateRequest define a struct for Instantiate method
type InstantiateRequest struct {
	// Admin is an optional address that can execute migrations
	Admin string
	// CodeID is the reference to the stored WASM code
	CodeID string
	// Label is optional metadata to be stored with a contract instance.
	Label string
	// InitMsg json encoded message to be passed to the contract on instantiation
	InitMsg Args
	// InitFunds coins that are transferred to the contract on instantiation
	InitFunds sdk.Coins
}

type Args map[string]interface{}

func NewArgs() Args {
	return make(map[string]interface{})
}

func (ar Args) MarshallJson() []byte {
	bz, _ := json.Marshal(ar)
	return bz
}

func (ar Args) Put(key string, value interface{}) Args {
	ar[key] = value
	return ar
}

//ContractABI define a message for executing contract
type ContractABI struct {
	Method string
	Args   Args
}

//NewContractABI return the instance of ContractABI
func NewContractABI() *ContractABI {
	return &ContractABI{
		Args: NewArgs(),
	}
}

//WithMethod set the field[method] for ContractABI
func (abi *ContractABI) WithMethod(method string) *ContractABI {
	abi.Method = method
	return abi
}

//WithArgs set the field[args] for ContractABI
func (abi *ContractABI) WithArgs(key string, value interface{}) *ContractABI {
	abi.Args[key] = value
	return abi
}

//Build marshal the ContractABI to []byte
func (abi ContractABI) Build() ([]byte, error) {
	if len(abi.Method) == 0 {
		return nil, errors.New("method is required")
	}
	data := make(map[string]interface{})
	data[abi.Method] = abi.Args
	return json.Marshal(data)
}
