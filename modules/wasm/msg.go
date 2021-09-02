package wasm

import (
	"encoding/json"
	"errors"

	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
)

// message types for the wasm client
const (
	RouterKey                  string = "wasm"
	TypeMsgStoreCode           string = "store-code"
	TypeMsgInstantiateContract string = "instantiate"
	TypeMsgExecuteContract     string = "execute"
	TypeMsgMigrateContract     string = "migrate"
	TypeUpdateAdmin            string = "update-contract-admin"
	TypeClearAdmin             string = "clear-contract-admin"
)

var (
	_ sdk.Msg = &MsgStoreCode{}
	_ sdk.Msg = &MsgInstantiateContract{}
	_ sdk.Msg = &MsgExecuteContract{}
	_ sdk.Msg = &MsgMigrateContract{}
	_ sdk.Msg = &MsgUpdateAdmin{}
	_ sdk.Msg = &MsgClearAdmin{}
)

// Route implement sdk.Msg
func (msg MsgStoreCode) Route() string {
	return RouterKey
}

// Type implement sdk.Msg
func (msg MsgStoreCode) Type() string {
	return TypeMsgStoreCode
}

// ValidateBasic implement sdk.Msg
func (msg MsgStoreCode) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(msg.Sender); err != nil {
		return err
	}

	if len(msg.WASMByteCode) == 0 {
		return errors.New("WASMByteCode should not be empty")
	}

	return nil
}

// GetSignBytes implement sdk.Msg
func (msg MsgStoreCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

// GetSigners implement sdk.Msg
func (msg MsgStoreCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}

// Route implement sdk.Msg
func (msg MsgInstantiateContract) Route() string {
	return RouterKey
}

// Type implement sdk.Msg
func (msg MsgInstantiateContract) Type() string {
	return TypeMsgInstantiateContract
}

// ValidateBasic implement sdk.Msg
func (msg MsgInstantiateContract) ValidateBasic() error {
	if msg.CodeID == 0 {
		return errors.New("code id is required")
	}
	if msg.Label == "" {
		return errors.New("label is required")
	}
	if len(msg.Admin) != 0 {
		if err := sdk.ValidateAccAddress(msg.Admin); err != nil {
			return err
		}
	}
	if !json.Valid(msg.InitMsg) {
		return errors.New("InitMsg is not valid json")
	}
	return sdk.ValidateAccAddress(msg.Sender)
}

// GetSignBytes implement sdk.Msg
func (msg MsgInstantiateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

// GetSigners implement sdk.Msg
func (msg MsgInstantiateContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}

// Route implement sdk.Msg
func (msg MsgExecuteContract) Route() string {
	return RouterKey
}

// Type implement sdk.Msg
func (msg MsgExecuteContract) Type() string {
	return TypeMsgExecuteContract
}

// ValidateBasic implement sdk.Msg
func (msg MsgExecuteContract) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(msg.Contract); err != nil {
		return err
	}
	if !json.Valid(msg.Msg) {
		return errors.New("InitMsg is not valid json")
	}
	return sdk.ValidateAccAddress(msg.Sender)
}

// GetSignBytes implement sdk.Msg
func (msg MsgExecuteContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

// GetSigners implement sdk.Msg
func (msg MsgExecuteContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}

// Route implement sdk.Msg
func (msg MsgMigrateContract) Route() string {
	return RouterKey
}

// Type implement sdk.Msg
func (msg MsgMigrateContract) Type() string {
	return TypeMsgMigrateContract
}

// ValidateBasic implement sdk.Msg
func (msg MsgMigrateContract) ValidateBasic() error {
	if msg.CodeID == 0 {
		return errors.New("code id is required")
	}

	if err := sdk.ValidateAccAddress(msg.Contract); err != nil {
		return err
	}

	if !json.Valid(msg.MigrateMsg) {
		return errors.New("migrate msg json")
	}
	return sdk.ValidateAccAddress(msg.Sender)
}

// GetSignBytes implement sdk.Msg
func (msg MsgMigrateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

// GetSigners implement sdk.Msg
func (msg MsgMigrateContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}

// Route implement sdk.Msg
func (msg MsgUpdateAdmin) Route() string {
	return RouterKey
}

// Type implement sdk.Msg
func (msg MsgUpdateAdmin) Type() string {
	return TypeUpdateAdmin
}

// ValidateBasic implement sdk.Msg
func (msg MsgUpdateAdmin) ValidateBasic() error {
	return sdk.ValidateAccAddress(msg.Sender)
}

// GetSignBytes implement sdk.Msg
func (msg MsgUpdateAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

// GetSigners implement sdk.Msg
func (msg MsgUpdateAdmin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}

// Route implement sdk.Msg
func (msg MsgClearAdmin) Route() string {
	return RouterKey
}

// Type implement sdk.Msg
func (msg MsgClearAdmin) Type() string {
	return TypeClearAdmin
}

// ValidateBasic implement sdk.Msg
func (msg MsgClearAdmin) ValidateBasic() error {
	return sdk.ValidateAccAddress(msg.Sender)
}

// GetSignBytes implement sdk.Msg
func (msg MsgClearAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

// GetSigners implement sdk.Msg
func (msg MsgClearAdmin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}
