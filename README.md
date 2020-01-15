# IRITA Network Go SDK

IRITA Network Go SDK provide a within warapper around the IRITA LCD API, in addition to creating and submitting different transaction.
It includes the following core components:

- **client**: provide httpClient, LiteClient, RpcClient and TxClient for query or send transaction on IRITA
- **keys**: implement KeyManage to manage private key and accounts
- **types**: common types
- **util**: define constant and common functions

# Install

## Requirement

Go version above 1.13

## Use go mod(recommend)

Add "github.com/bianjieai/irita-sdk-go" dependency into your go.mod file.

```
require (
	github.com/bianjieai/irita-sdk-go latest
)
```

# Usage

## Key Manager

Before start using API, you should construct a Key Manager to help sign the transaction msg or verify signature. Key Manager is an Identity Manger to define who you are in the IRITA

Wo provide follow construct functions to generate Key Mange(other keyManager will coming soon):

```
NewKeyStoreKeyManager(file string, auth string) (KeyManager, error)
```

Examples:

for mnemonic:

```
func TestRecoverFromMnemonic(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("faa", "fap")
	config.Seal()

	km := keyManager{}

	mnemonic := "situate wink injury solar orange ugly behave elite roast ketchup sand elephant monitor inherit canal menu demand hockey dose clap illness hurdle elbow high"
	password := ""
	fullPath := "44'/118'/0'/0/0"

	if err := km.recoverFromMnemonic(mnemonic, password, fullPath); err != nil {
		t.Fatal(err)
	} else {
		//assert.Equal(t, "faa1s4p3m36dcw5dga5z8hteeznvd8827ulhmm857j", km.addr.String())
		t.Log(km.GetAddr().String())
	}
}
```

## Init Client

```
import (
	"github.com/bianjieai/irita-sdk-go/client"
	"github.com/bianjieai/irita-sdk-go/types"
)
var (
	baseUrl, nodeUrl string
	networkType = types.Testnet
)
km, _ := keys.NewKeyManagerFromMnemonic(mnemonic, password, fullPath)
c, _ := client.NewIRITAClient(baseUrl, nodeUrl, networkType, km)
```

Note:
- `baseUrl`: should be lcd endpoint if you want to use liteClient
- `nodeUrl`: should be irisnet node endpoint, format is `tcp://host:port`
- `networkType`: `testnet` or `mainnet`

after you init irisnetClient, it include follow clients which you can use:

- `liteClient`: lcd client for IRITA
- `rpcClient`: query IRITA info by rpc
- `txClient`: send transaction to IRITA
