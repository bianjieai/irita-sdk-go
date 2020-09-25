package tx

import (
	"fmt"

	"github.com/bianjieai/irita-sdk-go/crypto/types"
	"github.com/bianjieai/irita-sdk-go/types/tx/signing"
)

// SignatureDataToModeInfoAndSig converts a SignatureData to a ModeInfo and raw bytes signature
func SignatureDataToModeInfoAndSig(data signing.SignatureData) (*ModeInfo, []byte) {
	if data == nil {
		return nil, nil
	}

	switch data := data.(type) {
	case *signing.SingleSignatureData:
		return &ModeInfo{
			Sum: &ModeInfo_Single_{
				Single: &ModeInfo_Single{Mode: data.SignMode},
			},
		}, data.Signature
	case *signing.MultiSignatureData:
		n := len(data.Signatures)
		modeInfos := make([]*ModeInfo, n)
		sigs := make([][]byte, n)

		for i, d := range data.Signatures {
			modeInfos[i], sigs[i] = SignatureDataToModeInfoAndSig(d)
		}

		multisig := types.MultiSignature{
			Signatures: sigs,
		}
		sig, err := multisig.Marshal()
		if err != nil {
			panic(err)
		}

		return &ModeInfo{
			Sum: &ModeInfo_Multi_{
				Multi: &ModeInfo_Multi{
					Bitarray:  data.BitArray,
					ModeInfos: modeInfos,
				},
			},
		}, sig
	default:
		panic(fmt.Sprintf("unexpected signature data type %T", data))
	}
}

// ModeInfoAndSigToSignatureData converts a ModeInfo and raw bytes signature to a SignatureData or returns
// an error
func ModeInfoAndSigToSignatureData(modeInfo *ModeInfo, sig []byte) (signing.SignatureData, error) {
	switch modeInfo := modeInfo.Sum.(type) {
	case *ModeInfo_Single_:
		return &signing.SingleSignatureData{
			SignMode:  modeInfo.Single.Mode,
			Signature: sig,
		}, nil

	case *ModeInfo_Multi_:
		multi := modeInfo.Multi

		sigs, err := decodeMultisignatures(sig)
		if err != nil {
			return nil, err
		}

		sigv2s := make([]signing.SignatureData, len(sigs))
		for i, mi := range multi.ModeInfos {
			sigv2s[i], err = ModeInfoAndSigToSignatureData(mi, sigs[i])
			if err != nil {
				return nil, err
			}
		}

		return &signing.MultiSignatureData{
			BitArray:   multi.Bitarray,
			Signatures: sigv2s,
		}, nil

	default:
		panic(fmt.Errorf("unexpected ModeInfo data type %T", modeInfo))
	}
}

// decodeMultisignatures safely decodes the the raw bytes as a MultiSignature protobuf message
func decodeMultisignatures(bz []byte) ([][]byte, error) {
	multisig := types.MultiSignature{}
	err := multisig.Unmarshal(bz)
	if err != nil {
		return nil, err
	}
	// NOTE: it is import to reject multi-signatures that contain unrecognized fields because this is an exploitable
	// malleability in the protobuf message. Basically an attacker could bloat a MultiSignature message with unknown
	// fields, thus bloating the transaction and causing it to fail.
	if len(multisig.XXX_unrecognized) > 0 {
		return nil, fmt.Errorf("rejecting unrecognized fields found in MultiSignature")
	}
	return multisig.Signatures, nil
}
