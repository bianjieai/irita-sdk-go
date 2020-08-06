package validator

import (
	"errors"
	"strings"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	ModuleName = "validator"
)

var (
	_ sdk.Msg = MsgCreateValidator{}
	_ sdk.Msg = MsgRemoveValidator{}
	_ sdk.Msg = MsgUpdateValidator{}

	amino = codec.New()

	// ModuleCdc references the global validator module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to validator and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	registerCodec(amino)
}

func (m MsgCreateValidator) Route() string {
	return ModuleName
}

func (m MsgCreateValidator) Type() string {
	return "create_validator"
}

func (m MsgCreateValidator) ValidateBasic() error {
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	if len(strings.TrimSpace(m.Name)) == 0 {
		return errors.New("validator name cannot be blank")
	}

	if len(m.Certificate) == 0 {
		return errors.New("certificate missing")
	}
	if m.Power <= 0 {
		return errors.New("power must be positive")
	}
	return nil
}

func (m MsgCreateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCreateValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

func (m MsgUpdateValidator) Route() string {
	return ModuleName
}

func (m MsgUpdateValidator) Type() string {
	return "update_validator"
}

func (m MsgUpdateValidator) ValidateBasic() error {
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	if len(m.ID) == 0 {
		return errors.New("validator id cannot be blank")
	}

	if m.Power < 0 {
		return errors.New("power can not be negative")
	}
	return nil
}

func (m MsgUpdateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUpdateValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

func (m MsgRemoveValidator) Route() string {
	return ModuleName
}

func (m MsgRemoveValidator) Type() string {
	return "remove_validator"
}

func (m MsgRemoveValidator) ValidateBasic() error {
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	if len(m.ID) == 0 {
		return errors.New("validator id cannot be blank")
	}
	return nil
}

func (m MsgRemoveValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgRemoveValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

func (v Validator) Convert() interface{} {
	return QueryValidatorResponse{
		ID:          v.ID.String(),
		Name:        v.Name,
		Pubkey:      v.Pubkey,
		Certificate: v.Certificate,
		Power:       v.Power,
		Details:     v.Description,
		Jailed:      v.Jailed,
		Operator:    v.Operator.String(),
	}
}

func (m MsgUnjailValidator) Route() string {
	return ModuleName
}

func (m MsgUnjailValidator) Type() string {
	return "unjail_validator"
}

func (m MsgUnjailValidator) ValidateBasic() error {
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	if len(m.ID) == 0 {
		return errors.New("validator name cannot be blank")
	}
	return nil
}

func (m MsgUnjailValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUnjailValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

type Validators []Validator

func (vs Validators) Convert() interface{} {
	var vrs []QueryValidatorResponse
	for _, v := range vs {
		vrs = append(vrs, QueryValidatorResponse{
			ID:          v.ID.String(),
			Name:        v.Name,
			Pubkey:      v.Pubkey,
			Certificate: v.Certificate,
			Power:       v.Power,
			Details:     v.Description,
			Jailed:      v.Jailed,
			Operator:    v.Operator.String(),
		})
	}
	return vrs
}

func (p Params) Convert() interface{} {
	return QueryParamsResponse{
		HistoricalEntries: p.HistoricalEntries,
	}
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "irita/modules/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgUpdateValidator{}, "irita/modules/MsgUpdateValidator", nil)
	cdc.RegisterConcrete(MsgRemoveValidator{}, "irita/modules/MsgRemoveValidator", nil)
	cdc.RegisterConcrete(MsgUnjailValidator{}, "irita/modules/MsgUnjailValidator", nil)
}
