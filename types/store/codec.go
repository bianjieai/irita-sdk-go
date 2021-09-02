package store

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/v2/codec"
	cryptoAmino "github.com/bianjieai/irita-sdk-go/v2/crypto/codec"
	"github.com/bianjieai/irita-sdk-go/v2/crypto/hd"
)

var cdc *codec.LegacyAmino

func init() {
	cdc = codec.NewLegacyAmino()
	cryptoAmino.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	cdc.Seal()
}

// RegisterCodec registers concrete types and interfaces on the given codec.
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Info)(nil), nil)
	cdc.RegisterConcrete(hd.BIP44Params{}, "crypto/keys/hd/BIP44Params", nil)
	cdc.RegisterConcrete(localInfo{}, "crypto/keys/localInfo", nil)
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey crypto.PubKey, err error) {
	err = cdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}
