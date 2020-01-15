package tx

import (
	"github.com/bianjieai/irita-sdk-go/util"
	"testing"
)

func TestClient_EditNFT(t *testing.T) {
	c := initNFTKM()

	req := NFTEditReq{
		Denom:    "HelloKitty",
		ID:       "kitty_1",
		TokenURI: "https://irisnet.org",
	}
	if res, err := c.EditNFT(req, "hello", true); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
