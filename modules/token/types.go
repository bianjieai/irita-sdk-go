package token

import (
	json2 "encoding/json"
	"errors"
	"strconv"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	ModuleName = "token"
)

var (
	_ sdk.Msg = &MsgIssueToken{}
	_ sdk.Msg = &MsgEditToken{}
	_ sdk.Msg = &MsgMintToken{}
	_ sdk.Msg = &MsgBurnToken{}
	_ sdk.Msg = &MsgTransferTokenOwner{}
)

func (msg MsgIssueToken) Route() string { return ModuleName }

// Implements Msg.
func (msg MsgIssueToken) Type() string { return "issue_token" }

// Implements Msg.
func (msg MsgIssueToken) ValidateBasic() error {
	if len(msg.Owner) == 0 {
		return errors.New("owner must be not empty")
	}

	if err := sdk.ValidateAccAddress(msg.Owner); err != nil {
		return sdk.Wrap(err)
	}

	if len(msg.Symbol) == 0 {
		return errors.New("symbol must be not empty")
	}

	if len(msg.Name) == 0 {
		return errors.New("name must be not empty")
	}

	if len(msg.MinUnit) == 0 {
		return errors.New("minUnit must be not empty")
	}

	return nil
}

// Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Owner)}
}

// GetSignBytes implements Msg
func (msg MsgTransferTokenOwner) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferTokenOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.SrcOwner)}
}

func (msg MsgTransferTokenOwner) ValidateBasic() error {
	if len(msg.SrcOwner) == 0 {
		return errors.New("srcOwner must be not empty")
	}

	if err := sdk.ValidateAccAddress(msg.SrcOwner); err != nil {
		return sdk.Wrap(err)
	}

	if len(msg.DstOwner) == 0 {
		return errors.New("dstOwner must be not empty")
	}

	if err := sdk.ValidateAccAddress(msg.DstOwner); err != nil {
		return sdk.Wrap(err)
	}

	if len(msg.Symbol) == 0 {
		return errors.New("symbol must be not empty")
	}

	return nil
}

func (msg MsgTransferTokenOwner) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgTransferTokenOwner) Type() string { return "transfer_token_owner" }

func (msg MsgEditToken) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgEditToken) Type() string { return "edit_token" }

// ValidateBasic implements Msg
func (msg MsgEditToken) ValidateBasic() error {
	if len(msg.Owner) == 0 {
		return errors.New("owner must be not empty")
	}

	if err := sdk.ValidateAccAddress(msg.Owner); err != nil {
		return sdk.Wrap(err)
	}

	if len(msg.Symbol) == 0 {
		return errors.New("symbol must be not empty")
	}
	return nil
}

// GetSignBytes implements Msg
func (msg MsgEditToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgEditToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Owner)}
}

func (msg MsgMintToken) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgMintToken) Type() string { return "mint_token" }

// GetSignBytes implements Msg
func (msg MsgMintToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Owner)}
}

// ValidateBasic implements Msg
func (msg MsgMintToken) ValidateBasic() error {
	if len(msg.Owner) == 0 {
		return errors.New("owner must be not empty")
	}

	if err := sdk.ValidateAccAddress(msg.Owner); err != nil {
		return sdk.Wrap(err)
	}

	if len(msg.Symbol) == 0 {
		return errors.New("symbol must be not empty")
	}
	return nil
}

func (msg MsgBurnToken) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgBurnToken) Type() string { return "burn_token" }

// GetSignBytes implements Msg
func (msg MsgBurnToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgBurnToken) ValidateBasic() error {
	if len(msg.Sender) == 0 {
		return errors.New("sender must be not empty")
	}

	if err := sdk.ValidateAccAddress(msg.Sender); err != nil {
		return sdk.Wrap(err)
	}

	if len(msg.Symbol) == 0 {
		return errors.New("symbol must be not empty")
	}

	if msg.Amount == 0 {
		return errors.New("invalid token amount ")
	}
	return nil
}

type Bool string

func (b Bool) ToBool() bool {
	v := string(b)
	if len(v) == 0 {
		return false
	}
	result, _ := strconv.ParseBool(v)
	return result
}

func (b Bool) String() string {
	return string(b)
}

// Marshal needed for protobuf compatibility
func (b Bool) Marshal() ([]byte, error) {
	return []byte(b), nil
}

// Unmarshal needed for protobuf compatibility
func (b *Bool) Unmarshal(data []byte) error {
	*b = Bool(data[:])
	return nil
}

// Marshals to JSON using string
func (b Bool) MarshalJSON() ([]byte, error) {
	return json2.Marshal(b.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (b *Bool) UnmarshalJSON(data []byte) error {
	var s string
	err := json2.Unmarshal(data, &s)
	if err != nil {
		return nil
	}
	*b = Bool(s)
	return nil
}

// GetSymbol implements exported.TokenI
func (t Token) GetSymbol() string {
	return t.Symbol
}

// GetName implements exported.TokenI
func (t Token) GetName() string {
	return t.Name
}

// GetScale implements exported.TokenI
func (t Token) GetScale() uint32 {
	return t.Scale
}

// GetMinUnit implements exported.TokenI
func (t Token) GetMinUnit() string {
	return t.MinUnit
}

// GetInitialSupply implements exported.TokenI
func (t Token) GetInitialSupply() uint64 {
	return t.InitialSupply
}

// GetMaxSupply implements exported.TokenI
func (t Token) GetMaxSupply() uint64 {
	return t.MaxSupply
}

// GetMintable implements exported.TokenI
func (t Token) GetMintable() bool {
	return t.Mintable
}

// GetOwner implements exported.TokenI
func (t Token) GetOwner() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(t.Owner)
}

func (t Token) Convert() interface{} {
	return sdk.Token{
		Symbol:        t.Symbol,
		Name:          t.Name,
		Scale:         t.Scale,
		MinUnit:       t.MinUnit,
		InitialSupply: t.InitialSupply,
		MaxSupply:     t.MaxSupply,
		Mintable:      t.Mintable,
		Owner:         t.Owner,
	}
}

type Tokens []TokenInterface

func (ts Tokens) Convert() interface{} {
	var tokens sdk.Tokens
	for _, t := range ts {
		tokens = append(tokens, sdk.Token{
			Symbol:        t.GetSymbol(),
			Name:          t.GetName(),
			Scale:         t.GetScale(),
			MinUnit:       t.GetMinUnit(),
			InitialSupply: t.GetInitialSupply(),
			MaxSupply:     t.GetMaxSupply(),
			Mintable:      t.GetMintable(),
			Owner:         t.GetOwner().String(),
		})
	}
	return tokens
}

type TokenInterface interface {
	GetSymbol() string
	GetName() string
	GetScale() uint32
	GetMinUnit() string
	GetInitialSupply() uint64
	GetMaxSupply() uint64
	GetMintable() bool
	GetOwner() sdk.AccAddress
}

func (p Params) Convert() interface{} {
	return QueryParamsResp{
		TokenTaxRate:      p.TokenTaxRate.String(),
		IssueTokenBaseFee: p.IssueTokenBaseFee.String(),
		MintTokenFeeRatio: p.MintTokenFeeRatio.String(),
	}
}

func (t QueryFeesResponse) Convert() interface{} {
	return QueryFeesResp(t)
}
