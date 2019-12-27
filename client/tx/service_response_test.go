package tx

import (
	"encoding/hex"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/basic"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/lcd"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/rpc"
	"gitlab.bianjie.ai/irita/irita-sdk-go/keys"
	commontypes "gitlab.bianjie.ai/irita/irita-sdk-go/types"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util"
	"testing"
)

func initSvcResponseKM() TxClient {
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

func TestClient_PostServiceResponse(t *testing.T) {
	c := initSvcResponseKM()
	response := ServiceResponse{
		ReqChainId: "irita-l1",
		RequestId:  "6527-6427-0",
		Provider:   "faa1mqvszlr9jfjw7dm5h9y8hf9yda2fg62uu4gxuk",
		Data:       hex.EncodeToString([]byte("service call response")),
		ErrorMsg:   "",
	}

	if res, err := c.PostServiceResponse(response, "test service response", true); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(util.ToJsonIgnoreErr(res)))
	}
}
