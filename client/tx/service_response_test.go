package tx

import (
	"encoding/hex"
	"github.com/irisnet/explorer/backend/utils"
	"github.com/irisnet/sdk-go/client/basic"
	"github.com/irisnet/sdk-go/client/lcd"
	"github.com/irisnet/sdk-go/client/rpc"
	"github.com/irisnet/sdk-go/keys"
	commontypes "github.com/irisnet/sdk-go/types"
	"testing"
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
		t.Log(string(utils.MarshalJsonIgnoreErr(res)))
	}
}
