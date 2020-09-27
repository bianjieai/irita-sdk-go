package nft

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

func (o Owner) Convert() interface{} {
	var idcs []IDC
	for _, idc := range o.IDCollections {
		idcs = append(idcs, IDC{
			Denom:    idc.Denom,
			TokenIDs: idc.Ids,
		})
	}
	return QueryOwnerResp{
		Address: o.Address.String(),
		IDCs:    idcs,
	}
}

type nft interface {
	GetID() string
	GetOwner() sdk.AccAddress
	GetURI() string
	GetName() string
	GetData() string
	Convert() interface{}
}

// GetID returns the ID of the token
func (this BaseNFT) GetID() string { return this.Id }

// GetOwner returns the account address that owns the NFTI
func (this BaseNFT) GetOwner() sdk.AccAddress { return this.Owner }

// GetURI returns the path to optional extra properties
func (this BaseNFT) GetURI() string { return this.URI }

// GetName returns the path to optional extra properties
func (this BaseNFT) GetName() string { return this.Name }

// GetData returns the metadata of nft
func (this BaseNFT) GetData() string { return this.Data }

func (this BaseNFT) Convert() interface{} {
	return QueryNFTResp{
		ID:      this.Id,
		Name:    this.Name,
		URI:     this.URI,
		Data:    this.Data,
		Creator: this.Owner.String(),
	}
}

type NFTs []BaseNFT

func (this Denom) Convert() interface{} {
	return QueryDenomResp{
		ID:      this.Id,
		Name:    this.Name,
		Schema:  this.Schema,
		Creator: this.Creator.String(),
	}
}

type denoms []Denom

func (this denoms) Convert() interface{} {
	var denoms []QueryDenomResp
	for _, denom := range this {
		denoms = append(denoms, denom.Convert().(QueryDenomResp))
	}
	return denoms
}

// QueryCollectionResp of non fungible tokens
type collection struct {
	Denom Denom `json:"denom" yaml:"denom"` // name of the collection; not exported to clients
	NFTs  NFTs  `json:"nfts" yaml:"nfts"`   // NFTs that belong to a collection
}

func (c Collection) Convert() interface{} {
	var nfts []QueryNFTResp
	for _, nft := range c.NFTs {
		nfts = append(nfts, QueryNFTResp{
			ID:      nft.GetID(),
			Name:    nft.GetName(),
			URI:     nft.GetURI(),
			Data:    nft.GetData(),
			Creator: nft.GetOwner().String(),
		})
	}
	return QueryCollectionResp{
		Denom: c.Denom.Convert().(QueryDenomResp),
		NFTs:  nfts,
	}
}
