package tx

import (
	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/unknownproto"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// DefaultTxDecoder returns a default protobuf TxDecoder using the provided Marshaler and PublicKeyCodec
func DefaultTxDecoder(cdc *codec.ProtoCodec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var raw TxRaw

		// reject all unknown proto fields in the root TxRaw
		if err := unknownproto.RejectUnknownFieldsStrict(txBytes, &raw); err != nil {
			return nil, err
		}

		if err := cdc.UnmarshalBinaryBare(txBytes, &raw); err != nil {
			return nil, err
		}

		var body TxBody

		// allow non-critical unknown fields in TxBody
		txBodyHasUnknownNonCriticals, err := unknownproto.RejectUnknownFields(raw.BodyBytes, &body, true)
		if err != nil {
			return nil, err
		}

		if err = cdc.UnmarshalBinaryBare(raw.BodyBytes, &body); err != nil {
			return nil, err
		}

		var authInfo AuthInfo

		// reject all unknown proto fields in AuthInfo
		if err = unknownproto.RejectUnknownFieldsStrict(raw.AuthInfoBytes, &authInfo); err != nil {
			return nil, err
		}

		if err = cdc.UnmarshalBinaryBare(raw.AuthInfoBytes, &authInfo); err != nil {
			return nil, err
		}

		theTx := &Tx{
			Body:       &body,
			AuthInfo:   &authInfo,
			Signatures: raw.Signatures,
		}

		return &wrapper{
			tx:                           theTx,
			bodyBz:                       raw.BodyBytes,
			authInfoBz:                   raw.AuthInfoBytes,
			txBodyHasUnknownNonCriticals: txBodyHasUnknownNonCriticals,
		}, nil
	}
}

// DefaultTxDecoder returns a default protobuf JSON TxDecoder using the provided Marshaler and PublicKeyCodec
func DefaultJSONTxDecoder(cdc *codec.ProtoCodec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var theTx Tx
		if err := cdc.UnmarshalJSON(txBytes, &theTx);err != nil {
			return nil, err
		}

		return &wrapper{
			tx: &theTx,
		}, nil
	}
}
