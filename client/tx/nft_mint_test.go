package tx

import (
	"github.com/bianjieai/irita-sdk-go/client/basic"
	"github.com/bianjieai/irita-sdk-go/client/lcd"
	"github.com/bianjieai/irita-sdk-go/client/rpc"
	"github.com/bianjieai/irita-sdk-go/keys"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/util"
	"testing"
)

func initNFTKM() TxClient {
	mnemonic := "crop ecology obey stone paper find eye accuse rely text sweet present sort prosper reform drop grow wave exchange garment draft amateur away law"
	password := ""
	fullPath := "44'/118'/0'/0/1"
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
