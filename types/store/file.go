package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/keyring"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/mitchellh/go-homedir"
	"github.com/mtibben/percent"
	"github.com/pkg/errors"

	cryptoamino "github.com/bianjieai/irita-sdk-go/crypto/codec"
	"github.com/bianjieai/irita-sdk-go/crypto/hd"
)

const (
	keyringFileDirName = "keyring-file"
)

var (
	_ KeyDAO = FileDAO{}

	filenameEscape = func(s string) string {
		return percent.Encode(s, "/")
	}
)

//Execute the local file system to realize the persistence of the key data, and the stored data is encrypted using `PBES2`.
//Can directly read the data of `iritacli` keys (--keyring-backend = file)
type FileDAO struct {
	dir string
}

func NewFileDAO(dir string) KeyDAO {
	fileDir := filepath.Join(dir, keyringFileDirName)
	return &FileDAO{dir: fileDir}
}

// Write will use user password to encrypt data and save to file, the file name is user name
func (f FileDAO) Write(name, password string, info KeyInfo) error {
	pubkey, err := PubKeyFromBytes(info.PubKey)
	if err != nil {
		return err
	}

	lInfo := localInfo{
		Name:         info.Name,
		PubKey:       pubkey,
		PrivKeyArmor: info.PrivKeyArmor,
		Algo:         hd.PubKeyType(info.Algo),
	}

	item := keyring.Item{
		Key:  name,
		Data: marshalInfo(lInfo),
	}

	bytes, err := json.Marshal(item)
	if err != nil {
		return err
	}

	token, err := jose.Encrypt(
		string(bytes), jose.PBES2_HS256_A128KW, jose.A256GCM, password,
		jose.Headers(map[string]interface{}{"created": time.Now().String()}),
	)
	if err != nil {
		return err
	}

	filename, err := f.filename(item.Key)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, []byte(token), 0600)
}

// Read will read encrypted data from file and decrypt with user password
func (f FileDAO) Read(name, password string) (KeyInfo, error) {
	filename, err := f.filename(name)
	if err != nil {
		return KeyInfo{}, err
	}

	if len(password) == 0 {
		return KeyInfo{}, fmt.Errorf("no password")
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return KeyInfo{}, errors.Wrap(err, "not found")
	}

	payload, _, err := jose.Decode(string(bytes), password)
	if err != nil {
		return KeyInfo{}, err
	}

	var decoded keyring.Item
	err = json.Unmarshal([]byte(payload), &decoded)
	if err != nil {
		return KeyInfo{}, err
	}

	info, err := unmarshalInfo(decoded.Data)
	if err != nil {
		return KeyInfo{}, err
	}

	i, ok := info.(localInfo)
	if !ok {
		return KeyInfo{}, fmt.Errorf("only support type KeyInfo")
	}

	return KeyInfo{
		Name:         i.Name,
		PubKey:       cryptoamino.MarshalPubkey(i.PubKey),
		PrivKeyArmor: i.PrivKeyArmor,
		Algo:         string(i.Algo),
	}, nil
}

// Delete will delete user data and use user password to verify permissions
func (f FileDAO) Delete(name, password string) error {
	//Perform security verification
	if _, err := f.Read(name, password); err != nil {
		return err
	}

	filename, err := f.filename(name)
	if err != nil {
		return err
	}

	return os.Remove(filename)
}

// Has returns whether the specified user name exists
func (f FileDAO) Has(name string) bool {
	filename, err := f.filename(name)
	if err != nil {
		return false
	}
	if _, err = os.Stat(filename); err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (f *FileDAO) filename(key string) (string, error) {
	dir, err := f.resolveDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, filenameEscape(string(infoKey(key)))), nil
}

func (f *FileDAO) resolveDir() (string, error) {
	if f.dir == "" {
		return "", fmt.Errorf("no directory provided for file keyring")
	}

	dir := f.dir

	// expand tilde for home directory
	if strings.HasPrefix(dir, "~") {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		dir = strings.Replace(dir, "~", home, 1)
	}

	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
	} else if err != nil && !stat.IsDir() {
		err = fmt.Errorf("%s is a file, not a directory", dir)
	}

	return dir, err
}
