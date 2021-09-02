package types

import (
	signingtypes "github.com/bianjieai/irita-sdk-go/v2/types/tx/signing"
)

type (
	// TxBuilder defines an interface which an application-defined concrete transaction
	// type must implement. Namely, it must be able to set messages, generate
	// signatures, and provide canonical bytes to sign over. The transaction must
	// also know how to encode itself.
	TxBuilder interface {
		GetTx() Tx

		SetMsgs(msgs ...Msg) error
		SetSignatures(signatures ...signingtypes.SignatureV2) error
		SetMemo(memo string)
		SetFeeAmount(amount Coins)
		SetGasLimit(limit uint64)
		SetTimeoutHeight(height uint64)
	}

	// TxEncodingConfig defines an interface that contains transaction
	// encoders and decoders
	TxEncodingConfig interface {
		TxEncoder() TxEncoder
		TxDecoder() TxDecoder
		TxJSONEncoder() TxEncoder
		TxJSONDecoder() TxDecoder
		MarshalSignatureJSON([]signingtypes.SignatureV2) ([]byte, error)
		UnmarshalSignatureJSON([]byte) ([]signingtypes.SignatureV2, error)
	}

	// TxConfig defines an interface a client can utilize to generate an
	// application-defined concrete transaction type. The type returned must
	// implement TxBuilder.
	TxConfig interface {
		TxEncodingConfig

		NewTxBuilder() TxBuilder
		WrapTxBuilder(Tx) (TxBuilder, error)
		SignModeHandler() SignModeHandler
	}
)
