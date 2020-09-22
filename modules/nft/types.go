package nft

import (
	"strings"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	ModuleName = "nft"
)

var (
	_ sdk.Msg = MsgIssueDenom{}
	_ sdk.Msg = MsgTransferNFT{}
	_ sdk.Msg = MsgEditNFT{}
	_ sdk.Msg = MsgMintNFT{}
	_ sdk.Msg = MsgBurnNFT{}

	amino = codec.New()

	// ModuleCdc references the global x/gov module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/gov and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	registerCodec(amino)
}

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
	id := strings.TrimSpace(m.ID)
	if len(id) == 0 {
		return sdk.Wrapf("missing id")
	}
	return nil
}

func (m MsgIssueDenom) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(m)
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

	tokenID := strings.TrimSpace(m.ID)
	if len(tokenID) == 0 {
		return sdk.Wrapf("missing ID")
	}
	return nil
}

func (m MsgTransferNFT) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(m)
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

	tokenID := strings.TrimSpace(m.ID)
	if len(tokenID) == 0 {
		return sdk.Wrapf("missing ID")
	}
	return nil
}

func (m MsgEditNFT) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(m)
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

	tokenID := strings.TrimSpace(m.ID)
	if len(tokenID) == 0 {
		return sdk.Wrapf("missing ID")
	}
	return nil
}

func (m MsgMintNFT) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(m)
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

	tokenID := strings.TrimSpace(m.ID)
	if len(tokenID) == 0 {
		return sdk.Wrapf("missing ID")
	}
	return nil
}

func (m MsgBurnNFT) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (m MsgBurnNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueDenom{}, "irismod/nft/MsgIssueDenom", nil)
	cdc.RegisterConcrete(MsgMintNFT{}, "irismod/nft/MsgMintNFT", nil)
	cdc.RegisterConcrete(MsgTransferNFT{}, "irismod/nft/MsgTransferNFT", nil)
	cdc.RegisterConcrete(MsgEditNFT{}, "irismod/nft/MsgEditNFT", nil)
	cdc.RegisterConcrete(MsgBurnNFT{}, "irismod/nft/MsgBurnNFT", nil)

	cdc.RegisterInterface((*nft)(nil), nil)
	cdc.RegisterConcrete(BaseNFT{}, "irismod/nft/BaseNFT", nil)
}
