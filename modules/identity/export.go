package identity

import (
	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
)

type Client interface {
	sdk.Module

	CreateIdentity(request CreateIdentityRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	UpdateIdentity(request UpdateIdentityRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryIdentity(id string) (QueryIdentityResp, sdk.Error)
}

type PubkeyInfo struct {
	PubKey     string          `json:"pub_key"`
	PubKeyAlgo PubKeyAlgorithm `json:"pub_key_algo"`
}

type CreateIdentityRequest struct {
	ID          string      `json:"id"`
	PubkeyInfo  *PubkeyInfo `json:"pubkey_info"`
	Certificate string      `json:"certificate"`
	Credentials *string     `json:"credentials"`
}

type UpdateIdentityRequest struct {
	ID          string      `json:"id"`
	PubkeyInfo  *PubkeyInfo `json:"pubkey_info"`
	Certificate string      `json:"certificate"`
	Credentials *string     `json:"credentials"`
}

type QueryIdentityResp struct {
	ID           string       `json:"id"`
	PubkeyInfos  []PubkeyInfo `json:"pubkey_infos"`
	Certificates []string     `json:"certificates"`
	Credentials  string       `json:"credentials"`
	Owner        string       `json:"owner"`
}
