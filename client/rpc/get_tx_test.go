package rpc

import (
	"github.com/bianjieai/irita-sdk-go/util"
	"testing"
)

func TestClient_GetTx(t *testing.T) {
	hash := "95FBE21354D2FF9A8827D8103C974968684384B7EEED44BBEA71D90B2C56FA09"
	if res, err := c.GetTx(hash); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
