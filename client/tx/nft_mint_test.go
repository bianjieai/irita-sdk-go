package tx

import (
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/basic"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/lcd"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/rpc"
	"gitlab.bianjie.ai/irita/irita-sdk-go/keys"
	"gitlab.bianjie.ai/irita/irita-sdk-go/types"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util"
	"testing"
)

func initNFTKM() TxClient {
	mnemonic := "situate wink injury solar orange ugly behave elite roast ketchup sand elephant monitor inherit canal menu demand hockey dose clap illness hurdle elbow high"
	password := ""
	fullPath := "44'/118'/0'/0/0"
	km, err := keys.NewKeyManagerFromMnemonic(mnemonic, password, fullPath)
	if err != nil {
		panic(err)
	}
	basicClient := basic.NewClient("http://localhost:1317")
	lite := lcd.NewClient(basicClient)
	rpcClient := rpc.NewClient("tcp://localhost:26657")

	c, err := NewClient("irita-l1", types.Testnet, km, lite, rpcClient)
	if err != nil {
		panic(err)
	}
	return c
}

func TestClient_MintNFT(t *testing.T) {
	c := initNFTKM()

	req := NFTMintReq{
		Recipient: "faa1s4p3m36dcw5dga5z8hteeznvd8827ulhmm857j",
		Denom:     "HelloKitty",
		ID:        "kitty_1",
		TokenURI:  "https://irita.org/1",
	}

	if res, err := c.MintNFT(req, "send from sdk", false); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
