package node

import (
	"errors"
	"strings"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	ModuleName = "node"
)

var (
	_ sdk.Msg = &MsgCreateValidator{}
	_ sdk.Msg = &MsgRemoveValidator{}
	_ sdk.Msg = &MsgUpdateValidator{}
)

func (m MsgCreateValidator) Route() string {
	return ModuleName
}

func (m MsgCreateValidator) Type() string {
	return "create_validator"
}

func (m MsgCreateValidator) ValidateBasic() error {
	if len(m.Operator) == 0 {
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
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCreateValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
}

func (m MsgUpdateValidator) Route() string {
	return ModuleName
}

func (m MsgUpdateValidator) Type() string {
	return "update_validator"
}

func (m MsgUpdateValidator) ValidateBasic() error {
	if len(m.Operator) == 0 {
		return errors.New("operator missing")
	}
	if len(m.Id) == 0 {
		return errors.New("validator id cannot be blank")
	}

	if m.Power < 0 {
		return errors.New("power can not be negative")
	}
	return nil
}

func (m MsgUpdateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUpdateValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
}

func (m MsgRemoveValidator) Route() string {
	return ModuleName
}

func (m MsgRemoveValidator) Type() string {
	return "remove_validator"
}

func (m MsgRemoveValidator) ValidateBasic() error {
	if len(m.Operator) == 0 {
		return errors.New("operator missing")
	}
	if len(m.Id) == 0 {
		return errors.New("validator id cannot be blank")
	}
	return nil
}

func (m MsgRemoveValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgRemoveValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
}

func (v Validator) Convert() interface{} {
	return QueryValidatorResp{
		ID:          v.Id,
		Name:        v.Name,
		Pubkey:      v.Pubkey,
		Certificate: v.Certificate,
		Power:       v.Power,
		Details:     v.Description,
		Jailed:      v.Jailed,
		Operator:    v.Operator,
	}
}

type validators []Validator

func (vs validators) Convert() interface{} {
	var vrs []QueryValidatorResp
	for _, v := range vs {
		vrs = append(vrs, QueryValidatorResp{
			ID:          v.Id,
			Name:        v.Name,
			Pubkey:      v.Pubkey,
			Certificate: v.Certificate,
			Power:       v.Power,
			Details:     v.Description,
			Jailed:      v.Jailed,
			Operator:    v.Operator,
		})
	}
	return vrs
}

func (p Params) Convert() interface{} {
	return QueryParamsResp(p)
}
