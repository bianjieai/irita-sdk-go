package types

import (
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/types/tx/signing"
)

type (
	// Msg defines the interface a transaction message must fulfill.
	Msg interface {
		proto.Message

		// Return the message type.
		// Must be alphanumeric or empty.
		Route() string

		// Returns a human-readable string for the message, intended for utilization
		// within tags
		Type() string

		// ValidateBasic does a simple validation check that
		// doesn't require access to any other information.
		ValidateBasic() error

		// Get the canonical byte representation of the Msg.
		GetSignBytes() []byte

		// Signers returns the addrs of signers that must sign.
		// CONTRACT: All signatures must be present to be valid.
		// CONTRACT: Returns addrs in some deterministic order.
		GetSigners() []AccAddress
	}

	// Fee defines an interface for an application application-defined concrete
	// transaction type to be able to set and return the transaction fee.
	Fee interface {
		GetGas() uint64
		GetAmount() Coins
	}

	// Signature defines an interface for an application application-defined
	// concrete transaction type to be able to set and return transaction signatures.
	Signature interface {
		GetPubKey() crypto.PubKey
		GetSignature() []byte
	}

	// Tx defines the interface a transaction must fulfill.
	Tx interface {
		// Gets the all the transaction's messages.
		GetMsgs() []Msg

		// ValidateBasic does a simple and lightweight validation check that doesn't
		// require access to any other information.
		ValidateBasic() error
	}

	// FeeTx defines the interface to be implemented by Tx to use the FeeDecorators
	FeeTx interface {
		Tx
		GetGas() uint64
		GetFee() Coins
		FeePayer() AccAddress
	}

	// Tx must have GetMemo() method to use ValidateMemoDecorator
	TxWithMemo interface {
		Tx
		GetMemo() string
	}

	// TxWithTimeoutHeight extends the Tx interface by allowing a transaction to
	// set a height timeout.
	TxWithTimeoutHeight interface {
		Tx

		GetTimeoutHeight() uint64
	}

	// Factory defines an interface which an application-defined concrete transaction
	// type must implement. Namely, it must be able to set messages, generate
	// signatures, and provide canonical bytes to sign over. The transaction must
	// also know how to encode itself.
	TxBuilder1 interface {
		GetTx() SigTx

		SetMsgs(msgs ...Msg) error
		SetSignatures(signatures ...signing.SignatureV2) error
		SetMemo(memo string)
		SetFeeAmount(amount Coins)
		SetGasLimit(limit uint64)
		SetTimeoutHeight(height uint64)
	}
)

// TxDecoder unmarshals transaction bytes
type TxDecoder func(txBytes []byte) (Tx, error)

// TxEncoder marshals transaction to bytes
type TxEncoder func(tx Tx) ([]byte, error)
