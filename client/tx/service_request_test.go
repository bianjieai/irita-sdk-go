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

func initSvcRequestKM() TxClient {
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

func TestClient_PostServiceRequest(t *testing.T) {
	c := initSvcRequestKM()

	data := "This is request data"
	requestData := hex.EncodeToString([]byte(data))
	request := ServiceRequest{
		ServiceName: "Material_Accept_Confirmation",
		BindChainId: "rainbow-dev",
		DefChainId:  "rainbow-dev",
		MethodId:    1,
		Provider:    "faa1mqvszlr9jfjw7dm5h9y8hf9yda2fg62uu4gxuk",
		Consumer:    "faa1t7jpg8pue93nzuxa6cr30ax0n0hu99p43mk6am",
		ServiceFee:  "1iris",
		Data:        requestData,
		Profiling:   false,
	}
	res, err := c.PostServiceRequest(request, "service call test", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(res)))
}
