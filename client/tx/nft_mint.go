package tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/modules/incubator/nft"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/types"
)

type (
	NFTMintReq struct {
		Recipient string `json:"recipient"`
		Denom     string `json:"denom"`
		ID        string `json:"id"`
		TokenURI  string `json:"tokenURI"`
	}
)

func (c *client) MintNFT(r NFTMintReq, memo string, commit bool) (types.BroadcastTxResult, error) {
	var (
		result types.BroadcastTxResult

		recipient sdk.AccAddress
	)

	// get account info
	sender := c.keyManager.GetAddr()
	account, err := c.liteClient.QueryAccount(sender.String())
	if err != nil {
		return result, err
	}

	// build msg
	if v, err := sdk.AccAddressFromBech32(r.Recipient); err != nil {
		return result, err
	} else {
		recipient = v
	}
	msg := nft.NewMsgMintNFT(sender, recipient, r.ID, r.Denom, r.TokenURI)
	stdSignMsg := buildStdSignMsg(c.chainId, memo, account, msg)

	return signAndBroadcastTx(c, stdSignMsg, commit)
}
