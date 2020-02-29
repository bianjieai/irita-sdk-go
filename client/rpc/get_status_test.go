package rpc

import (
	"github.com/bianjieai/irita-sdk-go/util"
	"testing"
)

var (
	c RPCClient
)

func TestMain(m *testing.M) {
	c1, err := NewClient("tcp://localhost:26657")
	if err != nil {
		panic(err)
	} else {
		c = c1
	}
	m.Run()
}

func TestClient_GetStatus(t *testing.T) {
	if res, err := c.GetStatus(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
