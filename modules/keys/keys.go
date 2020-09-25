package keys

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type keysClient struct {
	sdk.KeyManager
}

func NewClient(keyManager sdk.KeyManager) KeyI {
	return keysClient{keyManager}
}

func (k keysClient) Add(name, password string) (string, string, sdk.Error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, sdk.Wrap(err)
}

func (k keysClient) Recover(name, password, mnemonic string) (string, sdk.Error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic)
	return address, sdk.Wrap(err)
}

func (k keysClient) Import(name, password, privKeyArmor string) (string, sdk.Error) {
	address, err := k.KeyManager.Import(name, password, privKeyArmor)
	return address, sdk.Wrap(err)
}

func (k keysClient) Export(name, password string) (string, sdk.Error) {
	keystore, err := k.KeyManager.Export(name, password)
	return keystore, sdk.Wrap(err)
}

func (k keysClient) Delete(name, password string) sdk.Error {
	err := k.KeyManager.Delete(name, password)
	return sdk.Wrap(err)
}

func (k keysClient) Show(name, password string) (string, sdk.Error) {
	_, address, err := k.KeyManager.Find(name, password)
	if err != nil {
		return "", sdk.Wrap(err)
	}
	return address.String(), nil
}
