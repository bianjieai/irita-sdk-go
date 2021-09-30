package slashing

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	RouterKey              string = "slashing"
	TypeMsgUnjailValidator string = "unjail-validator"
)

var (
	_ sdk.Msg = &MsgUnjailValidator{}
)

// Route implement sdk.Msg
func (m MsgUnjailValidator) Route() string {
	return RouterKey
}

// Type implement sdk.Msg
func (m MsgUnjailValidator) Type() string {
	return TypeMsgUnjailValidator
}

func (m *MsgUnjailValidator) ValidateBasic() error {
	return sdk.ValidateAccAddress(m.Operator)
}

func (m *MsgUnjailValidator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgUnjailValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
}
