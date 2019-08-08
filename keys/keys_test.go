package keys

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
