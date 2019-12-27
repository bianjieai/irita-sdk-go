package lcd

import (
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/basic"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util"
	"testing"
)

var (
	c LiteClient
)

func TestMain(m *testing.M) {
	baseClient := basic.NewClient("http://localhost:1317")
	c = NewClient(baseClient)
	m.Run()
}

func TestClient_QueryAccount(t *testing.T) {
	address := "faa1gq9ccwgx92lzxzukmx8lkm8xlpangwtqec5u2x"
	if res, err := c.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
