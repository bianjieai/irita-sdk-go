package keys

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

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
