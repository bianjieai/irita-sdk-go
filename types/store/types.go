package store

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"github.com/bianjieai/irita-sdk-go/crypto/hd"
)

var (
	_ Info = &localInfo{}
)

// KeyType reflects a human-readable type for key listing.
type KeyType uint

// Info KeyTypes
const (
	TypeLocal KeyType = 0
)

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

type Crypto interface {
	Encrypt(data string, password string) (string, error)
	Decrypt(data string, password string) (string, error)
}

// Info is the publicly exposed information about a keypair
type Info interface {
	// Human-readable type for key listing
	GetType() KeyType
	// Name of the key
	GetName() string
	// Public key
	GetPubKey() crypto.PubKey
	// Bip44 Path
	GetPath() (*hd.BIP44Params, error)
	// Algo
	GetAlgo() hd.PubKeyType
}

// localInfo is the public information about a locally stored key
// Note: Algo must be last field in struct for backwards amino compatibility
type localInfo struct {
	Name         string        `json:"name"`
	PubKey       crypto.PubKey `json:"pubkey"`
	PrivKeyArmor string        `json:"privkey.armor"`
	Algo         hd.PubKeyType `json:"algo"`
}

// GetType implements Info interface
func (i localInfo) GetType() KeyType {
	return TypeLocal
}

// GetType implements Info interface
func (i localInfo) GetName() string {
	return i.Name
}

// GetType implements Info interface
func (i localInfo) GetPubKey() crypto.PubKey {
	return i.PubKey
}

// GetType implements Info interface
func (i localInfo) GetAlgo() hd.PubKeyType {
	return i.Algo
}

// GetType implements Info interface
func (i localInfo) GetPath() (*hd.BIP44Params, error) {
	return nil, fmt.Errorf("BIP44 Paths are not available for this type")
}

// encoding info
func marshalInfo(i Info) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(i)
}

// decoding info
func unmarshalInfo(bz []byte) (info Info, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(bz, &info)
	return
}
