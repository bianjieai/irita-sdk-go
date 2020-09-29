package nft

import (
	"strings"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	ModuleName = "nft"
)

var (
	_ sdk.Msg = &MsgIssueDenom{}
	_ sdk.Msg = &MsgTransferNFT{}
	_ sdk.Msg = &MsgEditNFT{}
	_ sdk.Msg = &MsgMintNFT{}
	_ sdk.Msg = &MsgBurnNFT{}
)

func (m MsgIssueDenom) Route() string {
	return ModuleName
}

func (m MsgIssueDenom) Type() string {
	return "issue_denom"
}

func (m MsgIssueDenom) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdk.Wrapf("missing sender address")
	}
	id := strings.TrimSpace(m.Id)
	if len(id) == 0 {
		return sdk.Wrapf("missing id")
	}
	return nil
}

func (m MsgIssueDenom) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(&m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (m MsgIssueDenom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

func (m MsgTransferNFT) Route() string {
	return ModuleName
}

func (m MsgTransferNFT) Type() string {
	return "transfer_nft"
}

func (m MsgTransferNFT) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdk.Wrapf("missing sender address")
	}
	if m.Recipient.Empty() {
		return sdk.Wrapf("missing recipient address")
	}

	denom := strings.TrimSpace(m.Denom)
	if len(denom) == 0 {
		return sdk.Wrapf("missing denom")
	}

	tokenID := strings.TrimSpace(m.Id)
	if len(tokenID) == 0 {
		return sdk.Wrapf("missing ID")
	}
	return nil
}

func (m MsgTransferNFT) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(&m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (m MsgTransferNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

func (m MsgEditNFT) Route() string {
	return ModuleName
}

func (m MsgEditNFT) Type() string {
	return "edit_nft"
}

func (m MsgEditNFT) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdk.Wrapf("missing sender address")
	}

	denom := strings.TrimSpace(m.Denom)
	if len(denom) == 0 {
		return sdk.Wrapf("missing denom")
	}

	tokenID := strings.TrimSpace(m.Id)
	if len(tokenID) == 0 {
		return sdk.Wrapf("missing ID")
	}
	return nil
}

func (m MsgEditNFT) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(&m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (m MsgEditNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

func (m MsgMintNFT) Route() string {
	return ModuleName
}

func (m MsgMintNFT) Type() string {
	return "mint_nft"
}

func (m MsgMintNFT) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdk.Wrapf("missing sender address")
	}

	denom := strings.TrimSpace(m.Denom)
	if len(denom) == 0 {
		return sdk.Wrapf("missing denom")
	}

	tokenID := strings.TrimSpace(m.Id)
	if len(tokenID) == 0 {
		return sdk.Wrapf("missing ID")
	}
	return nil
}

func (m MsgMintNFT) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(&m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (m MsgMintNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

func (m MsgBurnNFT) Route() string {
	return ModuleName
}

func (m MsgBurnNFT) Type() string {
	return "burn_nft"
}

func (m MsgBurnNFT) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdk.Wrapf("missing sender address")
	}

	denom := strings.TrimSpace(m.Denom)
	if len(denom) == 0 {
		return sdk.Wrapf("missing denom")
	}

	tokenID := strings.TrimSpace(m.Id)
	if len(tokenID) == 0 {
		return sdk.Wrapf("missing ID")
	}
	return nil
}

func (m MsgBurnNFT) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(&m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (m MsgBurnNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

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

func (c Collection) Convert() interface{} {
	var nfts []QueryNFTResp
	for _, nft := range c.NFTs {
		nfts = append(nfts, QueryNFTResp{
			ID:      nft.Id,
			Name:    nft.Name,
			URI:     nft.URI,
			Data:    nft.Data,
			Creator: nft.Owner.String(),
		})
	}
	return QueryCollectionResp{
		Denom: c.Denom.Convert().(QueryDenomResp),
		NFTs:  nfts,
	}
}
