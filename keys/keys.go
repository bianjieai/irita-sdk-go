package keys

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/types"
	ctypes "github.com/irisnet/irishub/types"
	"github.com/irisnet/sdk-go/types/tx"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"io/ioutil"
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
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
		PubKey:        k.privKey.PubKey(),
		Signature:     sigBytes,
	}, nil
}

func (k *keyManager) recoveryFromKeyStore(keystoreFile string, auth string) error {
	if auth == "" {
		return fmt.Errorf("Password is missing ")
	}
	keyJson, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return err
	}
	var encryptedKey EncryptedKeyJSON
	err = json.Unmarshal(keyJson, &encryptedKey)
	if err != nil {
		return err
	}
	keyBytes, err := decryptKey(&encryptedKey, auth)
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], keyBytes[:32])
	privKey := secp256k1.PrivKeySecp256k1(keyBytesArray)
	addr := ctypes.AccAddress(privKey.PubKey().Address())
	k.addr = addr
	k.privKey = privKey
	return nil
}

func NewKeyStoreKeyManager(file string, auth string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromKeyStore(file, auth)
	return &k, err
}
