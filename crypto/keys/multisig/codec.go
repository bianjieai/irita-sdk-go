package multisig

import (
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/sr25519"

	"github.com/bianjieai/irita-sdk-go/v2/crypto/keys/ed25519"
	"github.com/bianjieai/irita-sdk-go/v2/crypto/keys/secp256k1"
	cryptotypes "github.com/bianjieai/irita-sdk-go/v2/crypto/types"
)

// TODO: Figure out API for others to either add their own pubkey types, or
// to make verify / marshal accept a Cdc.
const (
	PubKeyAminoRoute = "tendermint/PubKeyMultisigThreshold"
)

var AminoCdc = amino.NewCodec()

func init() {
	AminoCdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	AminoCdc.RegisterInterface((*cryptotypes.PubKey)(nil), nil)
	AminoCdc.RegisterConcrete(ed25519.PubKey{},
		ed25519.PubKeyName, nil)
	AminoCdc.RegisterConcrete(sr25519.PubKey{},
		sr25519.PubKeyName, nil)
	AminoCdc.RegisterConcrete(&secp256k1.PubKey{},
		secp256k1.PubKeyName, nil)
	AminoCdc.RegisterConcrete(&LegacyAminoPubKey{},
		PubKeyAminoRoute, nil)
}
