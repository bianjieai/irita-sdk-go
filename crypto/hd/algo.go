package hd

import (
	"fmt"

	"github.com/bianjieai/irita-sdk-go/crypto/keys/sm2"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/go-bip39"

	"github.com/bianjieai/irita-sdk-go/crypto/keys/secp256k1"
)

type SignatureAlgo interface {
	Name() PubKeyType
	Derive() DeriveFn
	Generate() GenerateFn
}

func NewSigningAlgoFromString(str string) (SignatureAlgo, error) {
	switch str {
	case string(Secp256k1.Name()):
		return Secp256k1, nil
	case string(Sm2.Name()):
		return Sm2, nil
	default:
		return nil, fmt.Errorf("provided algorithm `%s` is not supported", str)
	}
}

// PubKeyType defines an algorithm to derive key-pairs which can be used for cryptographic signing.
type PubKeyType string

const (
	// Secp256k1Type uses the Bitcoin secp256k1 ECDSA parameters.
	Secp256k1Type = PubKeyType("secp256k1")
	// Ed25519Type represents the Ed25519Type signature system.
	// It is currently not supported for end-user keys (wallets/ledgers).
	Ed25519Type = PubKeyType("ed25519")
	// Sr25519Type represents the Sr25519Type signature system.
	Sr25519Type = PubKeyType("sr25519")

	// Sm2Type represents the Sm2Type signature system.
	Sm2Type = PubKeyType("sm2")
)

var (
	// Secp256k1 uses the Bitcoin secp256k1 ECDSA parameters.
	Secp256k1 = secp256k1Algo{}
	Sm2       = sm2Algo{}
)

type DeriveFn func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error)
type GenerateFn func(bz []byte) crypto.PrivKey

type WalletGenerator interface {
	Derive(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error)
	Generate(bz []byte) crypto.PrivKey
}

type secp256k1Algo struct{}

func (s secp256k1Algo) Name() PubKeyType {
	return Secp256k1Type
}

// Derive derives and returns the secp256k1 private key for the given seed and HD path.
func (s secp256k1Algo) Derive() DeriveFn {
	return func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error) {
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		masterPriv, ch := ComputeMastersFromSeed(seed)
		if len(hdPath) == 0 {
			return masterPriv[:], nil
		}
		derivedKey, err := DerivePrivateKeyForPath(masterPriv, ch, hdPath)
		return derivedKey[:], err
	}
}

// Generate generates a secp256k1 private key from the given bytes.
func (s secp256k1Algo) Generate() GenerateFn {
	return func(bz []byte) crypto.PrivKey {
		var bzArr [32]byte
		copy(bzArr[:], bz)
		return &secp256k1.PrivKey{Key: bzArr[:]}
	}
}

type sm2Algo struct{}

func (s sm2Algo) Name() PubKeyType {
	return Sm2Type
}

// Derive derives and returns the secp256k1 private key for the given seed and HD path.
func (s sm2Algo) Derive() DeriveFn {
	return func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error) {
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		masterPriv, ch := ComputeMastersFromSeed(seed)
		if len(hdPath) == 0 {
			return masterPriv[:], nil
		}
		derivedKey, err := DerivePrivateKeyForPath(masterPriv, ch, hdPath)
		return derivedKey[:], err
	}
}

// Generate generates a sm2 private key from the given bytes.
func (s sm2Algo) Generate() GenerateFn {
	return func(bz []byte) crypto.PrivKey {
		var bzArr [sm2.PrivKeySize]byte
		copy(bzArr[:], bz)
		return &sm2.PrivKey{Key: bzArr[:]}
	}
}
