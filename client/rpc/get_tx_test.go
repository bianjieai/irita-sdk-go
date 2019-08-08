package rpc

import (
	"github.com/irisnet/sdk-go/util"
	"testing"
)

func TestClient_GetTx(t *testing.T) {
	hash := "7854CD857E686550B679A9BF3118BDA7281A8B9979C215E0AEABACDA8C76B10B"
	if res, err := c.GetTx(hash); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
