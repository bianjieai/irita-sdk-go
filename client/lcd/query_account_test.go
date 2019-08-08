package lcd

import (
	"github.com/irisnet/sdk-go/client/basic"
	"github.com/irisnet/sdk-go/util"
	"testing"
)

var (
	c LiteClient
)

func TestMain(m *testing.M) {
	baseClient := basic.NewClient("http://v2.irisnet-lcd.dev.rainbow.one")
	c = NewClient(baseClient)
	m.Run()
}

func TestClient_QueryAccount(t *testing.T) {
	address := "faa1282eufkw9qgm55symgqqg38nremslvggpylkht"
	if res, err := c.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
