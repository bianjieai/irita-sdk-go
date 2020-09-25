package tx

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/unknownproto"
	cryptotypes "github.com/bianjieai/irita-sdk-go/crypto/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// DefaultTxDecoder returns a default protobuf TxDecoder using the provided Marshaler and PublicKeyCodec
func DefaultTxDecoder(cdc *codec.ProtoCodec, keyCodec cryptotypes.PublicKeyCodec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var raw TxRaw

		// reject all unknown proto fields in the root TxRaw
		err := unknownproto.RejectUnknownFieldsStrict(txBytes, &raw)
		if err != nil {
			return nil, err
		}

		err = cdc.UnmarshalBinaryBare(txBytes, &raw)
		if err != nil {
			return nil, err
		}

		var body TxBody

		// allow non-critical unknown fields in TxBody
		txBodyHasUnknownNonCriticals, err := unknownproto.RejectUnknownFields(raw.BodyBytes, &body, true)
		if err != nil {
			return nil, err
		}

		err = cdc.UnmarshalBinaryBare(raw.BodyBytes, &body)
		if err != nil {
			return nil, err
		}

		var authInfo AuthInfo

		// reject all unknown proto fields in AuthInfo
		err = unknownproto.RejectUnknownFieldsStrict(raw.AuthInfoBytes, &authInfo)
		if err != nil {
			return nil, err
		}

		err = cdc.UnmarshalBinaryBare(raw.AuthInfoBytes, &authInfo)
		if err != nil {
			return nil, err
		}

		theTx := &Tx{
			Body:       &body,
			AuthInfo:   &authInfo,
			Signatures: raw.Signatures,
		}

		pks, err := extractPubKeys(theTx, keyCodec)
		if err != nil {
			return nil, err
		}

		return &wrapper{
			tx:                           theTx,
			bodyBz:                       raw.BodyBytes,
			authInfoBz:                   raw.AuthInfoBytes,
			pubKeys:                      pks,
			pubkeyCodec:                  keyCodec,
			txBodyHasUnknownNonCriticals: txBodyHasUnknownNonCriticals,
		}, nil
	}
}

// DefaultTxDecoder returns a default protobuf JSON TxDecoder using the provided Marshaler and PublicKeyCodec
func DefaultJSONTxDecoder(cdc *codec.ProtoCodec, keyCodec cryptotypes.PublicKeyCodec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var theTx Tx
		err := cdc.UnmarshalJSON(txBytes, &theTx)
		if err != nil {
			return nil, err
		}

		pks, err := extractPubKeys(&theTx, keyCodec)
		if err != nil {
			return nil, err
		}

		return &wrapper{
			tx:          &theTx,
			pubKeys:     pks,
			pubkeyCodec: keyCodec,
		}, nil
	}
}

func extractPubKeys(tx *Tx, keyCodec cryptotypes.PublicKeyCodec) ([]crypto.PubKey, error) {
	if tx.AuthInfo == nil {
		return []crypto.PubKey{}, nil
	}

	signerInfos := tx.AuthInfo.SignerInfos
	pks := make([]crypto.PubKey, len(signerInfos))
	for i, si := range signerInfos {
		pk, err := keyCodec.Decode(si.PublicKey)
		if err != nil {
			return nil, err
		}
		pks[i] = pk
	}
	return pks, nil
}
