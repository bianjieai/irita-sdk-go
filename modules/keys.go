package modules

import (
	"fmt"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/crypto"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/types/store"
)

type keyManager struct {
	keyDAO store.KeyDAO
	algo   string
}

func (k keyManager) Sign(name, password string, data []byte) ([]byte, tmcrypto.PubKey, error) {
	info, err := k.keyDAO.Read(name, password)
	if err != nil {
		return nil, nil, fmt.Errorf("name %s not exist", name)
	}

	km, err := crypto.NewPrivateKeyManager([]byte(info.PrivKeyArmor), string(info.Algo))
	if err != nil {
		return nil, nil, fmt.Errorf("name %s not exist", name)
	}

	signByte, err := km.Sign(data)
	if err != nil {
		return nil, nil, err
	}

	return signByte, km.ExportPubKey(), nil
}

func (k keyManager) Insert(name, password string) (string, string, error) {
	if k.keyDAO.Has(name) {
		return "", "", fmt.Errorf("name %s has existed", name)
	}

	km, err := crypto.NewAlgoKeyManager(k.algo)
	if err != nil {
		return "", "", err
	}

	mnemonic, priv := km.Generate()

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       pubKey.Bytes(),
		PrivKeyArmor: string(priv.Bytes()),
		Algo:         k.algo,
	}

	if err = k.keyDAO.Write(name, password, info); err != nil {
		return "", "", err
	}

	return address, mnemonic, nil
}

func (k keyManager) Recover(name, password, mnemonic string) (string, error) {
	if k.keyDAO.Has(name) {
		return "", fmt.Errorf("name %s has existed", name)
	}

	km, err := crypto.NewMnemonicKeyManager(mnemonic, k.algo)
	if err != nil {
		return "", err
	}

	_, priv := km.Generate()

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       pubKey.Bytes(),
		PrivKeyArmor: string(priv.Bytes()),
		Algo:         k.algo,
	}

	if err = k.keyDAO.Write(name, password, info); err != nil {
		return "", err
	}

	return address, nil
}

func (k keyManager) Import(name, password, armor string) (string, error) {
	if k.keyDAO.Has(name) {
		return "", fmt.Errorf("%s has existed", name)
	}

	km := crypto.NewKeyManager()

	priv, _, err := km.ImportPrivKey(armor, password)
	if err != nil {
		return "", err
	}

	pubKey := km.ExportPubKey()
	address := types.AccAddress(pubKey.Address().Bytes()).String()

	info := store.KeyInfo{
		Name:         name,
		PubKey:       pubKey.Bytes(),
		PrivKeyArmor: string(priv.Bytes()),
		Algo:         k.algo,
	}

	if err = k.keyDAO.Write(name, password, info); err != nil {
		return "", err
	}

	return address, nil
}

func (k keyManager) Export(name, password string) (armor string, err error) {
	info, err := k.keyDAO.Read(name, password)
	if err != nil {
		return armor, fmt.Errorf("name %s not exist", name)
	}

	km, err := crypto.NewPrivateKeyManager([]byte(info.PrivKeyArmor), info.Algo)
	if err != nil {
		return "", err
	}

	return km.ExportPrivKey(password)
}

func (k keyManager) Delete(name, password string) error {
	return k.keyDAO.Delete(name, password)
}

func (k keyManager) Find(name, password string) (types.AccAddress, error) {
	info, err := k.keyDAO.Read(name, password)
	if err != nil {
		return nil, fmt.Errorf("name %s not exist", name)
	}

	pubKey, err := store.PubKeyFromBytes(info.PubKey)
	if err != nil {
		return nil, fmt.Errorf("name %s not exist", name)
	}

	return types.AccAddress(pubKey.Address().Bytes()), nil
}
