package tx

import (
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/codec"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type (
	// Generator defines an interface a client can utilize to generate an
	// application-defined concrete transaction type. The type returned must
	// implement ClientTx.
	Generator interface {
		NewTx() ClientTx
		NewFee() ClientFee
		NewSignature() ClientSignature
	}

	ClientFee interface {
		sdk.Fee
		SetGas(uint64)
		SetAmount(sdk.Coins)
	}

	ClientSignature interface {
		sdk.Signature
		SetPubKey(crypto.PubKey) error
		SetSignature([]byte)
	}

	// ClientTx defines an interface which an application-defined concrete transaction
	// type must implement. Namely, it must be able to set messages, generate
	// signatures, and provide canonical bytes to sign over. The transaction must
	// also know how to encode itself.
	ClientTx interface {
		sdk.Tx
		codec.ProtoMarshaler

		SetMsgs(...sdk.Msg) error
		GetSignatures() []sdk.Signature
		SetSignatures(...ClientSignature) error
		GetFee() sdk.Fee
		SetFee(ClientFee) error
		GetMemo() string
		SetMemo(string)

		// CanonicalSignBytes returns the canonical JSON bytes to sign over, given a
		// chain ID, along with an account and sequence number. The JSON encoding
		// ensures all field names adhere to their proto definition, default values
		// are omitted, and follows the JSON Canonical Form.
		CanonicalSignBytes(cid string, num, seq uint64) ([]byte, error)
	}
)

func (f Factory) BuildAndSignedTx(msgs ...sdk.Msg) ([]byte, error) {
	tx, err := f.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}
	return f.Sign(tx)
}

// BuildUnsignedTx builds a transaction to be signed given a set of messages. The
// transaction is initially created via the provided factory's generator. Once
// created, the fee, memo, and messages are set.
func (f Factory) BuildUnsignedTx(msgs ...sdk.Msg) (ClientTx, error) {
	if f.chainID == "" {
		return nil, fmt.Errorf("chain ID required but not specified")
	}

	fees := f.fees

	if !f.gasPrices.IsZero() {
		if !fees.IsZero() {
			return nil, errors.New("cannot provide both fees and gas prices")
		}

		glDec := sdk.NewDec(int64(f.gas))

		// Derive the fees based on the provided gas prices, where
		// fee = ceil(gasPrice * gasLimit).
		fees = make(sdk.Coins, len(f.gasPrices))

		for i, gp := range f.gasPrices {
			fee := gp.Amount.Mul(glDec)
			fees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}
	}

	clientFee := f.txGenerator.NewFee()
	clientFee.SetAmount(fees)
	clientFee.SetGas(f.gas)

	tx := f.txGenerator.NewTx()
	tx.SetMemo(f.memo)

	if err := tx.SetFee(clientFee); err != nil {
		return nil, err
	}

	if err := tx.SetSignatures(); err != nil {
		return nil, err
	}

	if err := tx.SetMsgs(msgs...); err != nil {
		return nil, err
	}

	return tx, nil
}

// BuildSimTx creates an unsigned tx with an empty single signature and returns
// the encoded transaction or an error if the unsigned transaction cannot be
// built.
func (f Factory) BuildSimTx(msgs ...sdk.Msg) ([]byte, error) {
	tx, err := f.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	// Create an empty signature literal as the ante handler will populate with a
	// sentinel pubkey.
	sig := f.txGenerator.NewSignature()

	if err := tx.SetSignatures(sig); err != nil {
		return nil, err
	}

	return tx.Marshal()
}

// Sign signs a given tx with the provided name and passphrase. If the Factory's
// Keybase is not set, a new one will be created based on the client's backend.
// The bytes signed over are canconical. The resulting signature will be set on
// the transaction. Finally, the marshaled transaction is returned. An error is
// returned upon failure.
//
// Note, It is assumed the Factory has the necessary fields set that are required
// by the CanonicalSignBytes call.
func (f Factory) Sign(tx ClientTx) ([]byte, error) {
	if f.keyManager == nil {
		return nil, errors.New("keybase must be set prior to signing a transaction")
	}

	signBytes, err := tx.CanonicalSignBytes(f.chainID, f.accountNumber, f.sequence)
	if err != nil {
		return nil, err
	}

	sigBytes, pubkey, err := f.keyManager.Sign(f.from, f.password, signBytes)
	if err != nil {
		return nil, err
	}

	sig := f.txGenerator.NewSignature()
	sig.SetSignature(sigBytes)

	if err := sig.SetPubKey(pubkey); err != nil {
		return nil, err
	}

	if err := tx.SetSignatures(sig); err != nil {
		return nil, err
	}

	return tx.Marshal()
}

// GasEstimateResponse defines a response definition for tx gas estimation.
type GasEstimateResponse struct {
	GasEstimate uint64 `json:"gas_estimate" yaml:"gas_estimate"`
}

func (gr GasEstimateResponse) String() string {
	return fmt.Sprintf("gas estimate: %d", gr.GasEstimate)
}
