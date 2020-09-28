package auth

import (
	"encoding/json"
	"errors"

	"github.com/tendermint/tendermint/crypto"

	cryptocodec "github.com/bianjieai/irita-sdk-go/crypto/codec"

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

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error
}

var _ Account = (*BaseAccount)(nil)

// GetAddress Implements sdk.Account.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

// SetAddress Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey() (pk crypto.PubKey) {
	if len(acc.PubKey) == 0 {
		return nil
	}

	pk, _ = cryptocodec.PubKeyFromBytes(acc.PubKey)
	return pk
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	if pubKey == nil {
		acc.PubKey = nil
	} else {
		acc.PubKey = pubKey.Bytes()
	}
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
	account := sdk.BaseAccount{
		Address:       acc.Address,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}

	var pkStr string
	if acc.PubKey == nil {
		return account
	}

	var pk crypto.PubKey
	pk, _ = cryptocodec.PubKeyFromBytes(acc.PubKey)

	pkStr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pk)
	if err != nil {
		panic(err)
	}

	account.PubKey = pkStr
	return account
}
