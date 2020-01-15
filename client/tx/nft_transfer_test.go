package tx

import (
	"github.com/bianjieai/irita-sdk-go/util"
	"testing"
)

func TestClient_TransferNFT(t *testing.T) {
	c := initNFTKM()
	req := NFTTransferReq{
		Denom:     "crypto-kitties",
		ID:        "we123",
		Recipient: "faa1f3npynk6ngzhz4d3jsr8n0704e4lgl55u0muv3",
		TokenURI:  "[do-not-modify]",
	}

	if res, err := c.TransferNFT(req, "transfer nft", true); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
