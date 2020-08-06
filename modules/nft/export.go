package nft

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// expose NFT module api for user
type NFTI interface {
	sdk.Module
	IssueDenom(request IssueDenomRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	MintNFT(request MintNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	EditNFT(request EditNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	TransferNFT(request TransferNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	BurnNFT(request BurnNFTRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QuerySupply(denomID, creator string) (uint64, sdk.Error)
	QueryOwner(creator, denomID string) (QueryOwnerResponse, sdk.Error)
	QueryCollection(denomID string) (QueryCollectionResponse, sdk.Error)
	QueryDenom(denomID string) (QueryDenomResponse, sdk.Error)
	QueryDenoms() ([]QueryDenomResponse, sdk.Error)
	QueryNFT(denomID, tokenID string) (QueryNFTResponse, sdk.Error)
}

type IssueDenomRequest struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type MintNFTRequest struct {
	Denom     string `json:"denom"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	URI       string `json:"uri"`
	Data      string `json:"data"`
	Recipient string `json:"recipient"`
}

type EditNFTRequest struct {
	Denom string `json:"denom"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	URI   string `json:"uri"`
	Data  string `json:"data"`
}

type TransferNFTRequest struct {
	Denom     string `json:"denom"`
	ID        string `json:"id"`
	URI       string `json:"uri"`
	Data      string `json:"data"`
	Name      string `json:"name"`
	Recipient string `json:"recipient"`
}

type BurnNFTRequest struct {
	Denom string `json:"denom"`
	ID    string `json:"id"`
}

// IDC defines a set of nft ids that belong to a specific
type IDC struct {
	Denom    string   `json:"denom" yaml:"denom"`
	TokenIDs []string `json:"token_ids" yaml:"token_ids"`
}

type QueryOwnerResponse struct {
	Address string `json:"address" yaml:"address"`
	IDCs    []IDC  `json:"idcs" yaml:"idcs"`
}

// BaseNFT non fungible token definition
type QueryNFTResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URI     string `json:"uri"`
	Data    string `json:"data"`
	Creator string `json:"creator"`
}

type QueryDenomResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Schema  string `json:"schema"`
	Creator string `json:"creator"`
}

type QueryCollectionResponse struct {
	Denom QueryDenomResponse `json:"denom" yaml:"denom"`
	NFTs  []QueryNFTResponse `json:"nfts" yaml:"nfts"`
}
