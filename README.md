# IRITA Go SDK

IRITA GO SDK makes a simple package of API provided by IRITA Chain, which provides great convenience for users to quickly develop applications based on IRITA.

## install

### Requirement

- Go version above 1.14.0

### Use Go Mod

```text
require (
    github.com/bianjieai/irita-sdk-go latest
)
```

## Usage

### Init Client

The initialization SDK code is as follows:

```go
options := []types.Option{
    types.KeyDAOOption(store.NewMemory(nil)),
    types.TimeoutOption(10),
}
cfg, err := types.NewClientConfig(nodeURI, chainID, options...)
if err != nil {
    panic(err)
}
client := sdk.NewIRITAClient(cfg)
```

The `ClientConfig` component mainly contains the parameters used in the SDK, the specific meaning is shown in the table below

| Iterm      | Type          | Description                                                                                           |
| ---------- | ------------- | ----------------------------------------------------------------------------------------------------- |
| NodeURI    | string        | The RPC address of the irita node connected to the SDK, for example: localhost: 26657                 |
| ChainID    | string        | ChainID of irita, for example: `irita`                                                                |
| Gas        | uint64        | The maximum gas to be paid for the transaction, for example: `20000`                                  |
| Fee        | DecCoins      | Transaction fees to be paid for transactions                                                          |
| KeyDAO     | KeyDAO        | Private key management interface, If the user does not provide it, the default `LevelDB` will be used |
| Mode       | enum          | Transaction broadcast mode, value: `Sync`,`Async`, `Commit`                                           |
| Algo       | enum          | Private key generation algorithm(`sm2`,`secp256k1`)                                                   |
| Timeout    | time.Duration | Transaction timeout, for example: `5s`                                                                |
| Level      | string        | Log output level, for example: `info`                                                                 |
| MaxTxBytes | uint64        | The maximum number of transaction bytes supported by the connected node, default: `1073741824`(5M)    |

If you want to use `SDK` to send a transfer transaction, the example is as follows:

```go
coins, err := types.ParseDecCoins("100point")
to := "caa1rgnu8grzt6mwnjg7jss7w0sfyjn67g4em9njf5"
baseTx := types.BaseTx{
    From:     "username",
    Gas:      20000,
    Memo:     "test",
    Mode:     types.Commit,
    Password: "password",
}

result, err := client.Bank.Send(to, coins, baseTx)
```

**Note**: If you use the relevant API for sending transactions, you should implement the `KeyDAO` interface.

### KeyDAO

 The interface definition is as follows:

```go
// KeyInfo saves the basic information of the key
type KeyInfo struct {
    Name         string `json:"name"`
    PubKey       []byte `json:"pubkey"`
    PrivKeyArmor string `json:"priv_key_armor"`
    Algo         string `json:"algo"`
}

type KeyDAO interface {
    // Write will use user password to encrypt data and save to file, the file name is user name
    Write(name, password string, store KeyInfo) error

    // Read will read encrypted data from file and decrypt with user password
    Read(name, password string) (KeyInfo, error)

    // Delete will delete user data and use user password to verify permissions
    Delete(name, password string) error

    // Has returns whether the specified user name exists
    Has(name string) bool
}
```

There are three different ways to implement the keyDAO interface in the SDK:

- Based on levelDB(`LevelDBDAO`)
- Based on local file system(`FileDAO`)
- Based on memory(`MemoryDAO`)

Located under package `types/store`
