package tx

import (
	"github.com/irisnet/modules/incubator/nft"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/types"
)

type (
	NFTEditReq struct {
		Denom    string `json:"denom"`
		ID       string `json:"id"`
		TokenURI string `json:"tokenURI"`
	}
)

func (c *client) EditNFT(r NFTEditReq, memo string, commit bool) (types.BroadcastTxResult, error) {
	var (
		result types.BroadcastTxResult
	)

	// get account info
	sender := c.keyManager.GetAddr()
	account, err := c.liteClient.QueryAccount(sender.String())
	if err != nil {
		return result, err
	}

	// build stdSignMsg
	msg := nft.NewMsgEditNFTMetadata(sender, r.ID, r.Denom, r.TokenURI)
	stdSignMsg := buildStdSignMsg(c.chainId, memo, account, msg)

	return signAndBroadcastTx(c, stdSignMsg, commit)
}
