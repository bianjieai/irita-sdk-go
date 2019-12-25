package tx

import (
	"encoding/hex"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/basic"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/lcd"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/rpc"
	"gitlab.bianjie.ai/irita/irita-sdk-go/keys"
	commontypes "gitlab.bianjie.ai/irita/irita-sdk-go/types"
	"testing"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util"
)

func initSvcResponseKM() TxClient {
	km, err := keys.NewKeyStoreKeyManager("./ks_manifest_1234567890.json", "1234567890")
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

func TestClient_PostServiceResponse(t *testing.T) {
	c := initSvcResponseKM()
	response := ServiceResponse{
		ReqChainId: "rainbow-dev",
		RequestId:  "1194337-1194237-0",
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
