package tx

import (
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/crypto"

	codectypes "github.com/bianjieai/irita-sdk-go/codec/types"
	"github.com/bianjieai/irita-sdk-go/crypto/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/types/tx/signing"
)

// wrapper is a wrapper around the Tx proto.Message which retain the raw
// body and auth_info bytes.
type wrapper struct {
	tx *Tx

	// bodyBz represents the protobuf encoding of TxBody. This should be encoding
	// from the client using TxRaw if the tx was decoded from the wire
	bodyBz []byte

	// authInfoBz represents the protobuf encoding of TxBody. This should be encoding
	// from the client using TxRaw if the tx was decoded from the wire
	authInfoBz []byte

	// pubKeys represents the cached crypto.PubKey's that were set either from tx decoding
	// or decoded from AuthInfo when GetPubKey's was called
	pubKeys []crypto.PubKey

	pubkeyCodec types.PublicKeyCodec

	txBodyHasUnknownNonCriticals bool
}

var (
	_ ExtensionOptionsTxBuilder = &wrapper{}
	_ codectypes.IntoAny        = &wrapper{}
)

// ExtensionOptionsTxBuilder defines a Factory that can also set extensions.
type ExtensionOptionsTxBuilder interface {
	SetExtensionOptions(...*codectypes.Any)
	SetNonCriticalExtensionOptions(...*codectypes.Any)
}

func newBuilder(pubkeyCodec types.PublicKeyCodec) *wrapper {
	return &wrapper{
		tx: &Tx{
			Body: &TxBody{},
			AuthInfo: &AuthInfo{
				Fee: &Fee{},
			},
		},
		pubkeyCodec: pubkeyCodec,
	}
}

func (w *wrapper) GetMsgs() []sdk.Msg {
	return w.tx.GetMsgs()
}

func (w *wrapper) ValidateBasic() error {
	return w.tx.ValidateBasic()
}

func (w *wrapper) getBodyBytes() []byte {
	if len(w.bodyBz) == 0 {
		// if bodyBz is empty, then marshal the body. bodyBz will generally
		// be set to nil whenever SetBody is called so the result of calling
		// this method should always return the correct bytes. Note that after
		// decoding bodyBz is derived from TxRaw so that it matches what was
		// transmitted over the wire
		var err error
		w.bodyBz, err = proto.Marshal(w.tx.Body)
		if err != nil {
			panic(err)
		}
	}
	return w.bodyBz
}

func (w *wrapper) getAuthInfoBytes() []byte {
	if len(w.authInfoBz) == 0 {
		// if authInfoBz is empty, then marshal the body. authInfoBz will generally
		// be set to nil whenever SetAuthInfo is called so the result of calling
		// this method should always return the correct bytes. Note that after
		// decoding authInfoBz is derived from TxRaw so that it matches what was
		// transmitted over the wire
		var err error
		w.authInfoBz, err = proto.Marshal(w.tx.AuthInfo)
		if err != nil {
			panic(err)
		}
	}
	return w.authInfoBz
}

func (w *wrapper) GetSigners() []sdk.AccAddress {
	return w.tx.GetSigners()
}

func (w *wrapper) GetPubKeys() []crypto.PubKey {
	if w.pubKeys == nil {
		signerInfos := w.tx.AuthInfo.SignerInfos
		pubKeys := make([]crypto.PubKey, len(signerInfos))

		for i, si := range signerInfos {
			var err error
			pk := si.PublicKey
			if pk != nil {
				pubKeys[i], err = w.pubkeyCodec.Decode(si.PublicKey)
				if err != nil {
					panic(err)
				}
			}
		}

		w.pubKeys = pubKeys
	}

	return w.pubKeys
}

func (w *wrapper) GetGas() uint64 {
	return w.tx.AuthInfo.Fee.GasLimit
}

func (w *wrapper) GetFee() sdk.Coins {
	return w.tx.AuthInfo.Fee.Amount
}

func (w *wrapper) FeePayer() sdk.AccAddress {
	return w.GetSigners()[0]
}

func (w *wrapper) GetMemo() string {
	return w.tx.Body.Memo
}

func (w *wrapper) GetSignatures() [][]byte {
	return w.tx.Signatures
}

// GetTimeoutHeight returns the transaction's timeout height (if set).
func (w *wrapper) GetTimeoutHeight() uint64 {
	return w.tx.Body.TimeoutHeight
}

func (w *wrapper) SetMsgs(msgs ...sdk.Msg) error {
	anys := make([]*codectypes.Any, len(msgs))

	for i, msg := range msgs {
		var err error
		anys[i], err = codectypes.NewAnyWithValue(msg)
		if err != nil {
			return err
		}
	}

	w.tx.Body.Messages = anys

	// set bodyBz to nil because the cached bodyBz no longer matches tx.Body
	w.bodyBz = nil

	return nil
}

// SetTimeoutHeight sets the transaction's height timeout.
func (w *wrapper) SetTimeoutHeight(height uint64) {
	w.tx.Body.TimeoutHeight = height

	// set bodyBz to nil because the cached bodyBz no longer matches tx.Body
	w.bodyBz = nil
}

func (w *wrapper) SetMemo(memo string) {
	w.tx.Body.Memo = memo

	// set bodyBz to nil because the cached bodyBz no longer matches tx.Body
	w.bodyBz = nil
}

func (w *wrapper) SetGasLimit(limit uint64) {
	if w.tx.AuthInfo.Fee == nil {
		w.tx.AuthInfo.Fee = &Fee{}
	}

	w.tx.AuthInfo.Fee.GasLimit = limit

	// set authInfoBz to nil because the cached authInfoBz no longer matches tx.AuthInfo
	w.authInfoBz = nil
}

func (w *wrapper) SetFeeAmount(coins sdk.Coins) {
	if w.tx.AuthInfo.Fee == nil {
		w.tx.AuthInfo.Fee = &Fee{}
	}

	w.tx.AuthInfo.Fee.Amount = coins

	// set authInfoBz to nil because the cached authInfoBz no longer matches tx.AuthInfo
	w.authInfoBz = nil
}

func (w *wrapper) SetSignatures(signatures ...signing.SignatureV2) error {
	n := len(signatures)
	signerInfos := make([]*SignerInfo, n)
	rawSigs := make([][]byte, n)

	for i, sig := range signatures {
		var modeInfo *ModeInfo
		modeInfo, rawSigs[i] = SignatureDataToModeInfoAndSig(sig.Data)
		pk, err := w.pubkeyCodec.Encode(sig.PubKey)
		if err != nil {
			return err
		}
		signerInfos[i] = &SignerInfo{
			PublicKey: pk,
			ModeInfo:  modeInfo,
			Sequence:  sig.Sequence,
		}
	}

	w.setSignerInfos(signerInfos)
	w.setSignatures(rawSigs)

	return nil
}

func (w *wrapper) setSignerInfos(infos []*SignerInfo) {
	w.tx.AuthInfo.SignerInfos = infos
	// set authInfoBz to nil because the cached authInfoBz no longer matches tx.AuthInfo
	w.authInfoBz = nil
	// set cached pubKeys to nil because they no longer match tx.AuthInfo
	w.pubKeys = nil
}

func (w *wrapper) setSignatures(sigs [][]byte) {
	w.tx.Signatures = sigs
}

func (w *wrapper) GetTx() sdk.Tx {
	return w
}

// GetProtoTx returns the tx as a proto.Message.
func (w *wrapper) AsAny() *codectypes.Any {
	// We're sure here that w.tx is a proto.Message, so this will call
	// codectypes.NewAnyWithValue under the hood.
	return codectypes.UnsafePackAny(w.tx)
}

func (w *wrapper) GetExtensionOptions() []*codectypes.Any {
	return w.tx.Body.ExtensionOptions
}

func (w *wrapper) GetNonCriticalExtensionOptions() []*codectypes.Any {
	return w.tx.Body.NonCriticalExtensionOptions
}

func (w *wrapper) SetExtensionOptions(extOpts ...*codectypes.Any) {
	w.tx.Body.ExtensionOptions = extOpts
	w.bodyBz = nil
}

func (w *wrapper) SetNonCriticalExtensionOptions(extOpts ...*codectypes.Any) {
	w.tx.Body.NonCriticalExtensionOptions = extOpts
	w.bodyBz = nil
}
