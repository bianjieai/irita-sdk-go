package rpc

import (
	"github.com/irisnet/sdk-go/util"
	"testing"
)

var (
	c RPCClient
)

func TestMain(m *testing.M) {
	c = NewClient("tcp://35.236.146.181:30657")
	m.Run()
}

func TestClient_GetStatus(t *testing.T) {
	if res, err := c.GetStatus(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
