package types

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/codec"
	cdctypes "github.com/bianjieai/irita-sdk-go/codec/types"
)

//The purpose of this interface is to convert the irita system type to the user receiving type
// and standardize the user interface
type Response interface {
	Convert() interface{}
}

type SplitAble interface {
	Len() int
	Sub(begin, end int) SplitAble
}

type Module interface {
	Name() string
	RegisterCodec(cdc *codec.LegacyAmino)
	RegisterInterfaceTypes(registry cdctypes.InterfaceRegistry)
}

type KeyManager interface {
	Sign(name, password string, data []byte) ([]byte, crypto.PubKey, error)
	Insert(name, password string) (string, string, error)
	Recover(name, password, mnemonic string) (string, error)
	Import(name, password string, privKeyArmor string) (address string, err error)
	Export(name, password string) (privKeyArmor string, err error)
	Delete(name, password string) error
	Find(name, password string) (address AccAddress, err error)
}
