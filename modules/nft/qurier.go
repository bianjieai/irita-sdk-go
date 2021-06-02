package nft

import sdk "github.com/bianjieai/irita-sdk-go/types"

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
func (this BaseNFT) GetOwner() sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(this.Owner)
	if err != nil {
		return sdk.AccAddress{}
	}

	return owner
}

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
		Creator: this.Owner,
	}
}
