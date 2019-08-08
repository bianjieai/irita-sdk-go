# IRIS Network Go SDK

IRIS Network Go SDK provide a within warapper around the IRISnet LCD API, in addition to creating and submitting different transaction.
It includes the following core components:

- **client**: provide httpClient, LiteClient, RpcClient and TxClient for query or send transaction on IRISnet
- **keys**: implement KeyManage to manage private key and accounts
- **types**: common types
- **util**: define constant and common functions

# Install

## Requirement

Go version above 1.11

## Use go mod(recommend)

Add "github.com/irisnet/sdk-go" dependency into your go.mod file.

```go
require (
	github.com/irisnet/sdk-go latest
)
```

# Usage

## Key Manager

Before start using API, you should construct a Key Manager to help sign the transaction msg or verify signature. Key Manager is an Identity Manger to define who you are in the IRISnet

Wo provide follow construct functions to generate Key Mange(other keyManager will coming soon):

```go
NewKeyStoreKeyManager(file string, auth string) (KeyManager, error)
```

Examples:

for keyStore:

```go
func TestNewKeyStoreKeyManager(t *testing.T) {
	file := "./ks_1234567890.json"
	if km, err := NewKeyStoreKeyManager(file, "1234567890"); err != nil {
		t.Fatal(err)
	} else {
		msg := []byte("hello world")
		signature, err := km.GetPrivKey().Sign(msg)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(km.GetAddr().String())


		assert.Equal(t, km.GetPrivKey().PubKey().VerifyBytes(msg, signature), true)
	}
}
```

## Init Client

```go
import (
	"github.com/irisnet/sdk-go/client"
	"github.com/irisnet/sdk-go/types"
)
var (
	baseUrl, nodeUrl string
	networkType = types.Testnet
)
km, _ := keys.NewKeyStoreKeyManager("../keys/ks_1234567890.json", "1234567890")
c, _ := client.NewIRISnetClient(baseUrl, nodeUrl, networkType, km)
```

Note:
- `baseUrl`: should be lcd endpoint if you want to use liteClient
- `nodeUrl`: should be irisnet node endpoint, format is `tcp://host:port`
- `networkType`: `testnet` or `mainnet`

after you init irisnetClient, it include follow clients which you can use:

- `liteClient`: lcd client for IRISnet
- `rpcClient`: query IRISnet info by rpc
- `txClient`: send transaction on IRISnet

