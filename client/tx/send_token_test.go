package tx

import (
	"fmt"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/basic"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/lcd"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/rpc"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/types"
	"gitlab.bianjie.ai/irita/irita-sdk-go/keys"
	commontypes "gitlab.bianjie.ai/irita/irita-sdk-go/types"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util"
	"math"
	"testing"
)

func initSendTokenKM() TxClient {
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

	c, err := NewClient("irita-l1", commontypes.Testnet, km, lite, rpcClient)
	if err != nil {
		panic(err)
	}
	return c
}

func TestClient_SendToken(t *testing.T) {
	c := initSendTokenKM()

	receiver := "faa1j3ufmgwe2cuumj7423jt4creqlcskltn6ht5w9"
	amount := fmt.Sprintf("%.0f", 0.12*math.Pow10(2))
	coins := []types.Coin{
		{
			Denom:  "irita",
			Amount: amount,
		},
	}
	memo := "send from irita/sdk-go"
	if res, err := c.SendToken(receiver, coins, memo, false); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
