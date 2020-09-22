package store

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	dbm "github.com/tendermint/tm-db"
)

const (
	keyDBName  = "keys"
	infoSuffix = "info"
)

var (
	_ KeyDAO = LevelDBDAO{}
)

type LevelDBDAO struct {
	db dbm.DB
	Crypto
}

// NewLevelDB initialize a keybase based on the configuration.
// Use leveldb as storage
func NewLevelDB(rootDir string, crypto Crypto) (KeyDAO, error) {
	db, err := dbm.NewGoLevelDB(keyDBName, filepath.Join(rootDir, "keys"))
	if err != nil {
		return nil, err
	}

	if crypto == nil {
		crypto = AES{}
	}

	levelDB := LevelDBDAO{
		db:     db,
		Crypto: crypto,
	}
	return levelDB, nil
}

// Write add a key information to the local store
func (k LevelDBDAO) Write(name, password string, info KeyInfo) error {
	if k.Has(name) {
		return fmt.Errorf("name %s has exist", name)
	}

	privStr, err := k.Encrypt(info.PrivKeyArmor, password)
	if err != nil {
		return err
	}

	info.PrivKeyArmor = privStr

	bz, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return k.db.SetSync(infoKey(name), bz)
}

// Read read a key information from the local store
func (k LevelDBDAO) Read(name, password string) (store KeyInfo, err error) {
	bz, err := k.db.Get(infoKey(name))
	if bz == nil || err != nil {
		return store, err
	}

	if err := json.Unmarshal(bz, &store); err != nil {
		return store, err
	}

	if len(password) > 0 {
		privStr, err := k.Decrypt(store.PrivKeyArmor, password)
		if err != nil {
			return store, err
		}
		store.PrivKeyArmor = privStr
	}
	return
}

// ReadMetadata read a key information from the local store
func (k LevelDBDAO) ReadMetadata(name string) (store KeyInfo, err error) {
	bz, err := k.db.Get(infoKey(name))
	if bz == nil || err != nil {
		return store, err
	}

	if err := json.Unmarshal(bz, &store); err != nil {
		return store, err
	}
	return
}

// Delete delete a key from the local store
func (k LevelDBDAO) Delete(name, password string) error {
	_, err := k.Read(name, password)
	if err != nil {
		return err
	}
	return k.db.DeleteSync(infoKey(name))
}

// Delete delete a key from the local store
func (k LevelDBDAO) Has(name string) bool {
	existed, err := k.db.Has(infoKey(name))
	if err != nil {
		return false
	}
	return existed
}

func infoKey(name string) []byte {
	return []byte(fmt.Sprintf("%s.%s", name, infoSuffix))
}
