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

func initSvcRequestKM() TxClient {
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

func TestClient_PostServiceRequest(t *testing.T) {
	c := initSvcRequestKM()

	data := "This is request data"
	requestData := hex.EncodeToString([]byte(data))
	request := ServiceRequest{
		ServiceName: "Material_Accept_Confirmation",
		BindChainId: "irita-l1",
		DefChainId:  "irita-l1",
		MethodId:    1,
		Provider:    "faa1q2mevumtwk9vw2ejq6drm2f098ehaapwkye38a",
		Consumer:    "faa1s4p3m36dcw5dga5z8hteeznvd8827ulhmm857j",
		ServiceFee:  "1irita",
		Data:        requestData,
		Profiling:   false,
	}
	res, err := c.PostServiceRequest(request, "service call test", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(util.ToJsonIgnoreErr(res)))
}
