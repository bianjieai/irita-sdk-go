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
	km, err := keys.NewKeyStoreKeyManager("./ks_1234567890.json", "1234567890")
	if err != nil {
		panic(err)
	}
	basicClient := basic.NewClient("http://irisnet-lcd.dev.rainbow.one")
	lite := lcd.NewClient(basicClient)
	rpcClient := rpc.NewClient("tcp://192.168.150.31:26657")

	c, err := NewClient("rainbow-dev", commontypes.Testnet, km, lite, rpcClient)
	if err != nil {
		panic(err)
	}
	return c
}

func TestClient_SendToken(t *testing.T) {
	c := initSendTokenKM()

	receiver := "faa1j3ufmgwe2cuumj7423jt4creqlcskltn6ht5w9"
	amount := fmt.Sprintf("%.0f", 0.12*math.Pow10(18))
	coins := []types.Coin{
		{
			Denom:  "iris-atto",
			Amount: amount,
		},
	}
	memo := "send from irisnet/sdk-go"
	if res, err := c.SendToken(receiver, coins, memo, false); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
