package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bianjieai/irita-sdk-go/v2/crypto"
	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
)

func TestNewMnemonicKeyManager(t *testing.T) {
	mnemonic := "nerve leader thank marriage spice task van start piece crowd run hospital control outside cousin romance left choice poet wagon rude climb leisure spring"

	km, err := crypto.NewMnemonicKeyManager(mnemonic, "sm2")
	assert.NoError(t, err)

	pubKey := km.ExportPubKey()
	pubkeyBech32, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey)
	assert.NoError(t, err)
	assert.Equal(t, "iap1qg0f84wcp06t5ajvdf6dxhnwu0hhgjv3ulvpvy9cklqp374t5sty55dw3ps", pubkeyBech32)

	address := sdk.AccAddress(pubKey.Address()).String()
	assert.Equal(t, "iaa1yh6ke44anmv92g9w3r3rf0lpaxhjrenrshc4am", address)
}
