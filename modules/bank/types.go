package bank

import (
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/crypto/sm2"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	maxMsgLen  = 5
	ModuleName = "bank"
)

var (
	_ sdk.Msg = MsgSend{}

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

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) MsgSend {
	return MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: amount}
}

func (msg MsgSend) Route() string {
	return ModuleName
}

func (msg MsgSend) Type() string {
	return "send"
}

func (msg MsgSend) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return errors.New("missing sender address")
	}
	if msg.ToAddress.Empty() {
		return errors.New("missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return errors.New("invalid coins")
	}
	if !msg.Amount.IsAllPositive() {
		return errors.New("invalid coins")
	}
	return nil
}

func (msg MsgSend) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiSend(in []Input, out []Output) MsgMultiSend {
	return MsgMultiSend{Inputs: in, Outputs: out}
}

func (msg MsgMultiSend) Route() string { return ModuleName }

// Implements Msg.
func (msg MsgMultiSend) Type() string { return "multisend" }

// Implements Msg.
func (msg MsgMultiSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return errors.New("invalid input coins")
	}
	if len(msg.Outputs) == 0 {
		return errors.New("invalid output coins")
	}
	// make sure all inputs and outputs are individually valid
	var totalIn, totalOut sdk.Coins
	for _, in := range msg.Inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}
		totalIn = totalIn.Add(in.Coins...)
	}
	for _, out := range msg.Outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}
		totalOut = totalOut.Add(out.Coins...)
	}
	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return errors.New("inputs and outputs don't match")
	}
	return nil
}

// Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgMultiSend) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	if len(in.Address) == 0 {
		return errors.New(fmt.Sprintf(fmt.Sprintf("account %s is invalid", in.Address.String())))
	}
	if in.Coins.Empty() {
		return errors.New("empty input coins")
	}
	if !in.Coins.IsValid() {
		return errors.New(fmt.Sprintf("invalid input coins [%s]", in.Coins))
	}
	return nil
}

// NewInput - create a transaction input, used with MsgSend
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	input := Input{
		Address: addr,
		Coins:   coins,
	}
	return input
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	if len(out.Address) == 0 {
		return errors.New(fmt.Sprintf(fmt.Sprintf("account %s is invalid", out.Address.String())))
	}
	if out.Coins.Empty() {
		return errors.New("empty input coins")
	}
	if !out.Coins.IsValid() {
		return errors.New(fmt.Sprintf("invalid input coins [%s]", out.Coins))
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgSend
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	output := Output{
		Address: addr,
		Coins:   coins,
	}
	return output
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "cosmos-sdk/MsgMultiSend", nil)

	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "cosmos-sdk/Account", nil)
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		secp256k1.PubKeyAminoName, nil)
	cdc.RegisterConcrete(sm2.PubKeySm2{},
		sm2.PubKeyAminoName, nil)
	cdc.RegisterConcrete(multisig.PubKeyMultisigThreshold{},
		multisig.PubKeyMultisigThresholdAminoRoute, nil)

	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{},
		ed25519.PrivKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{},
		secp256k1.PrivKeyAminoName, nil)
	cdc.RegisterConcrete(sm2.PrivKeySm2{},
		sm2.PrivKeyAminoName, nil)

}
