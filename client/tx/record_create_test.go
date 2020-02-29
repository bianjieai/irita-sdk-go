package tx

import (
	"github.com/bianjieai/irita-sdk-go/client/basic"
	"github.com/bianjieai/irita-sdk-go/client/lcd"
	"github.com/bianjieai/irita-sdk-go/client/rpc"
	"github.com/bianjieai/irita-sdk-go/keys"
	commontypes "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/util"
	"testing"
)

func initCreateRecordKM() TxClient {
	mnemonic := "situate wink injury solar orange ugly behave elite roast ketchup sand elephant monitor inherit canal menu demand hockey dose clap illness hurdle elbow high"
	password := ""
	fullPath := "44'/118'/0'/0/0"
	km, err := keys.NewKeyManagerFromMnemonic(mnemonic, password, fullPath)
	if err != nil {
		panic(err)
	}
	basicClient := basic.NewClient("http://localhost:1317")
	lite := lcd.NewClient(basicClient)
	rpcClient, err := rpc.NewClient("tcp://localhost:26657")
	if err != nil {
		panic(err)
	}

	c, err := NewClient("irita", commontypes.Testnet, km, lite, rpcClient)
	if err != nil {
		panic(err)
	}
	return c
}

func TestClient_CreateRecord(t *testing.T) {
	c := initCreateRecordKM()

	req := RecordCreateReq{
		Contents: []Content{
			{
				Digest:     "hash123",
				DigestAlgo: "md5",
				URI:        "www.rainbow.one",
				Meta:       "rainbow app",
			},
		},
		Creator: "faa1s4p3m36dcw5dga5z8hteeznvd8827ulhmm857j",
	}

	if res, err := c.CreateRecord(req, "test record create", true); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
