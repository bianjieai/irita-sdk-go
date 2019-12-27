package tx

import (
	"gitlab.bianjie.ai/irita/irita-sdk-go/util"
	"testing"
)

func TestClient_TransferNFT(t *testing.T) {
	c := initNFTKM()
	req := NFTTransferReq{
		Denom:     "HelloKitty",
		ID:        "kitty_1",
		Recipient: "faa16y5ylrlnd9s3jc6xvg8y4k3mt3t84h5zgqcdn2",
		TokenURI:  "12345",
	}

	if res, err := c.TransferNFT(req, "transfer nft", true); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
