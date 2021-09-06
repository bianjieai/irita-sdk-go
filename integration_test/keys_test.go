package integration_test

import (
	"fmt"
	sdk "github.com/bianjieai/irita-sdk-go"
	"github.com/bianjieai/irita-sdk-go/crypto"
	"github.com/bianjieai/irita-sdk-go/crypto/codec"
	"github.com/bianjieai/irita-sdk-go/modules"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/types/store"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/sm2"
)

func (s IntegrationTestSuite) TestKeys() {
	name, password := s.RandStringOfLength(20), s.RandStringOfLength(8)

	address, mnemonic, err := s.Key.Add(name, password)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), address)
	require.NotEmpty(s.T(), mnemonic)

	address1, err := s.Key.Show(name, password)
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, address1)

	privKeyArmor, err := s.Key.Export(name, password)
	require.NoError(s.T(), err)

	err = s.Key.Delete(name, password)
	require.NoError(s.T(), err)

	address2, err := s.Key.Import(name, password, privKeyArmor)
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, address2)

	err = s.Key.Delete(name, password)
	require.NoError(s.T(), err)

	address3, err := s.Key.Recover(name, password, mnemonic)
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, address3)
}

func (s IntegrationTestSuite) TestImportAccount() {
	mnemonic := "enlist front supreme key favorite stem wrestle client trip orange burst abandon someone beyond umbrella crop unknown bind ivory february grace elephant detect board"
	keyDao := store.NewMemory(nil)
	km, err := crypto.NewMnemonicKeyManager(mnemonic, sm2.KeyType)
	require.NoError(s.T(), err)

	address := types.AccAddress(km.ExportPubKey().Address()).String()
	fmt.Println(address)

	_, priv := km.Generate()
	ki := store.KeyInfo{
		Name:         "user",
		PubKey:       codec.MarshalPubkey(km.ExportPubKey()),
		PrivKeyArmor: string(codec.MarshalPrivKey(priv)),
		Algo:         sm2.KeyType,
	}
	err = keyDao.Write("user", "password", ki)
	require.NoError(s.T(), err)
	options := []types.Option{
		types.KeyDAOOption(keyDao),
		types.TimeoutOption(10),
	}

	cfg, err := types.NewClientConfig(nodeURI, grpcAddr, chainID, options...)
	require.NoError(s.T(), err)
	client := sdk.NewIRITAClient(cfg)
	require.Equal(s.T(), address, "iaa1nylsmagjt8jck7uy4pg9lhgac0002gsjquvs2z")
	acc, err := client.Bank.QueryAccount(address)
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, acc.Address.String())

}

func (s IntegrationTestSuite) TestRecoverAccount() {
	mnemonic := "enlist front supreme key favorite stem wrestle client trip orange burst abandon someone beyond umbrella crop unknown bind ivory february grace elephant detect board"
	address, err := s.Key.Recover("user", "passwd", mnemonic)
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, "iaa1nylsmagjt8jck7uy4pg9lhgac0002gsjquvs2z")

}

func (s IntegrationTestSuite) TestInitKeyManager() {
	keyDao := store.NewMemory(nil)
	keyManger := modules.NewKeyManager(keyDao, "sm2")
	address, err := keyManger.Recover("user", "password", "enlist front supreme key favorite stem wrestle client trip orange burst abandon someone beyond umbrella crop unknown bind ivory february grace elephant detect board")
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, "iaa1nylsmagjt8jck7uy4pg9lhgac0002gsjquvs2z")
}
