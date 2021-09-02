package types

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/v2/types/tx/signing"
)

// SigVerifiableTx defines a transaction interface for all signature verification
// handlers.
type SigVerifiableTx interface {
	Tx
	GetSigners() []AccAddress
	GetPubKeys() []crypto.PubKey // If signer already has pubkey in context, this list will have nil in its place
	GetSignaturesV2() ([]signing.SignatureV2, error)
}

// Tx defines a transaction interface that supports all standard message, signature
// fee, memo, and auxiliary interfaces.
type SigTx interface {
	SigVerifiableTx

	TxWithMemo
	FeeTx
	TxWithTimeoutHeight
}
