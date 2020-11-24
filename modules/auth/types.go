package auth

import (
	"encoding/json"
	"errors"

	"github.com/tendermint/tendermint/crypto"

	codectypes "github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// Account is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements Account.
type Account interface {
	GetAddress() sdk.AccAddress
	SetAddress(sdk.AccAddress) error // errors if already set.

	GetPubKey(unpacker codectypes.AnyUnpacker) (crypto.PubKey,error)
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error
}

var _ Account = (*BaseAccount)(nil)

// GetAddress Implements sdk.Account.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(acc.Address)
	return addr
}

// SetAddress Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}

	acc.Address = addr.String()
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey(unpacker codectypes.AnyUnpacker) (pk crypto.PubKey,err error) {
	if acc.PubKey == nil {
		return nil,nil
	}

	var pubKey crypto.PubKey
	if err = unpacker.UnpackAny(acc.PubKey, &pubKey);err != nil {
		return nil,err
	}
	return pubKey,nil
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	any, err := codectypes.PackAny(pubKey)
	if err != nil {
		return err
	}
	acc.PubKey = any
	return nil
}

// GetAccountNumber Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

func (acc BaseAccount) String() string {
	out, _ := json.Marshal(acc)
	return string(out)
}

// Convert return a sdk.BaseAccount
func (acc *BaseAccount) Convert() interface{} {
	return sdk.BaseAccount{
		Address:       sdk.MustAccAddressFromBech32(acc.Address),
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}
}
