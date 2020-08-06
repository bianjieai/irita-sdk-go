// Package keys allows you to manage your local tendermint keystore (wallets) for iris.
//
// **NOTE:** You need to implement the [[KeyDAO]] Interface first.
//
// As a quick start:
//
// CreateRecord a new key.
//
//  client := test.NewClient()
//	name, password := "test2", "1234567890"
//
//	address, mnemonic, err := client.KeyI.Add(name, password)
//	require.NoError(client.T(), err)
//	require.NotEmpty(client.T(), address)
//	require.NotEmpty(client.T(), mnemonic)
//
package keys
