package keys

import (
	"github.com/bianjieai/irita-sdk-go/types/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type KeyManager interface {
	Sign(msg tx.StdSignMsg) ([]byte, error)
	GetPrivKey() crypto.PrivKey
	GetAddr() types.AccAddress
}

type keyManager struct {
	privKey  crypto.PrivKey
	addr     types.AccAddress
	mnemonic string
}

func (k *keyManager) Sign(msg tx.StdSignMsg) ([]byte, error) {
	sig, err := k.makeSignature(msg)
	if err != nil {
		return nil, err
	}

	newTx := auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo)
	bz, err := tx.Cdc.MarshalBinaryLengthPrefixed(newTx)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func (k *keyManager) GetPrivKey() crypto.PrivKey {
	return k.privKey
}

func (k *keyManager) GetAddr() types.AccAddress {
	return k.addr
}

func (k *keyManager) makeSignature(msg tx.StdSignMsg) (sig auth.StdSignature, err error) {
	if err != nil {
		return
	}
	sigBytes, err := k.privKey.Sign(msg.Bytes())
	if err != nil {
		return
	}
	return auth.StdSignature{
		PubKey:    k.privKey.PubKey(),
		Signature: sigBytes,
	}, nil
}

func (k *keyManager) recoverFromMnemonic(mnemonic, password, fullPath string) error {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return err
	}

	masterPriv, chainCode := hd.ComputeMastersFromSeed(seed)
	privateKey, err := hd.DerivePrivateKeyForPath(masterPriv, chainCode, fullPath)

	if err != nil {
		return err
	}

	k.privKey = secp256k1.PrivKeySecp256k1(privateKey)
	k.addr = types.AccAddress(k.privKey.PubKey().Address())
	return nil
}

func NewKeyManagerFromMnemonic(mnemonic, password, fullPath string) (KeyManager, error) {
	km := keyManager{}
	return &km, km.recoverFromMnemonic(mnemonic, password, fullPath)
}
