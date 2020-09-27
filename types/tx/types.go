package tx

import (
	"errors"
	"fmt"

	codectypes "github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// MaxGasWanted defines the max gas allowed.
const MaxGasWanted = uint64((1 << 63) - 1)

var _, _ codectypes.UnpackInterfacesMessage = &Tx{}, &TxBody{}
var _ sdk.Tx = &Tx{}

// GetMsgs implements the GetMsgs method on sdk.Tx.
func (t *Tx) GetMsgs() []sdk.Msg {
	if t == nil || t.Body == nil {
		return nil
	}

	anys := t.Body.Messages
	res := make([]sdk.Msg, len(anys))
	for i, any := range anys {
		msg := any.GetCachedValue().(sdk.Msg)
		res[i] = msg
	}
	return res
}

// ValidateBasic implements the ValidateBasic method on sdk.Tx.
func (t *Tx) ValidateBasic() error {
	if t == nil {
		return fmt.Errorf("bad Tx")
	}

	body := t.Body
	if body == nil {
		return fmt.Errorf("missing TxBody")
	}

	authInfo := t.AuthInfo
	if authInfo == nil {
		return fmt.Errorf("missing AuthInfo")
	}

	fee := authInfo.Fee
	if fee == nil {
		return fmt.Errorf("missing fee")
	}

	if fee.GasLimit > MaxGasWanted {
		return fmt.Errorf(
			"invalid gas supplied; %d > %d", fee.GasLimit, MaxGasWanted,
		)
	}

	if fee.Amount.IsAnyNegative() {
		return fmt.Errorf(
			"invalid fee provided: %s", fee.Amount,
		)
	}

	sigs := t.Signatures

	if len(sigs) == 0 {
		return errors.New("no signatures supplied")
	}

	if len(sigs) != len(t.GetSigners()) {
		return fmt.Errorf(
			"wrong number of signers; expected %d, got %d", t.GetSigners(), len(sigs),
		)
	}

	return nil
}

// GetSigners retrieves all the signers of a tx.
func (t *Tx) GetSigners() []sdk.AccAddress {
	var signers []sdk.AccAddress
	seen := map[string]bool{}

	for _, msg := range t.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	return signers
}

// UnpackInterfaces implements the UnpackInterfaceMessages.UnpackInterfaces method
func (t *Tx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if t.Body != nil {
		return t.Body.UnpackInterfaces(unpacker)
	}
	return nil
}

// UnpackInterfaces implements the UnpackInterfaceMessages.UnpackInterfaces method
func (m *TxBody) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, any := range m.Messages {
		var msg sdk.Msg
		if err := unpacker.UnpackAny(any, &msg); err != nil {
			return err
		}
	}
	return nil
}

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.tx.v1beta1.Tx", (*sdk.Tx)(nil))
	registry.RegisterImplementations((*sdk.Tx)(nil), &Tx{})
}
