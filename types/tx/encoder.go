package tx

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/bianjieai/irita-sdk-go/codec"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// DefaultTxEncoder returns a default protobuf TxEncoder using the provided Marshaler
func DefaultTxEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		txWrapper, ok := tx.(*wrapper)
		if !ok {
			return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, tx)
		}

		raw := &TxRaw{
			BodyBytes:     txWrapper.getBodyBytes(),
			AuthInfoBytes: txWrapper.getAuthInfoBytes(),
			Signatures:    txWrapper.tx.Signatures,
		}

		return proto.Marshal(raw)
	}
}

// DefaultJSONTxEncoder returns a default protobuf JSON TxEncoder using the provided Marshaler.
func DefaultJSONTxEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		txWrapper, ok := tx.(*wrapper)
		if ok {
			return codec.ProtoMarshalJSON(txWrapper.tx)
		}

		protoTx, ok := tx.(*Tx)
		if ok {
			return codec.ProtoMarshalJSON(protoTx)
		}

		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, tx)

	}
}
