package lcd

import (
	"github.com/bianjieai/irita-sdk-go/client/basic"
	"github.com/bianjieai/irita-sdk-go/util"
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
	address := "faa16y5ylrlnd9s3jc6xvg8y4k3mt3t84h5zgqcdn2"
	if res, err := c.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
