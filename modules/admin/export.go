package admin

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// Client export a group api for Admin module
type Client interface {
	sdk.Module

	AddRoles(address string, roles []Role, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	RemoveRoles(address string, roles []Role, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	BlockAccount(address string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	UnblockAccount(address string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryRoles(address string) ([]Role, sdk.Error)
	QueryBlacklist(page, limit int) ([]string, sdk.Error)
}
