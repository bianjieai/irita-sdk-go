package tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/modules/incubator/nft"
	"github.com/bianjieai/irita-sdk-go/client/types"
)

type (
	NFTTransferReq struct {
		Denom     string `json:"denom"`
		ID        string `json:"id"`
		Recipient string `json:"recipient"`
		TokenURI  string `json:"tokenURI"` // use "[do-not-modify]" will not change tokenURI when transfer nft
	}
)

func (c *client) TransferNFT(r NFTTransferReq, memo string, commit bool) (types.BroadcastTxResult, error) {
	var (
		result types.BroadcastTxResult
	)

	// get send info
	sender := c.keyManager.GetAddr()
	accInfo, err := c.liteClient.QueryAccount(sender.String())
	if err != nil {
		return result, err
	}

	// build stdSignMsg
	recipient, err := sdk.AccAddressFromBech32(r.Recipient)
	if err != nil {
		return result, err
	}
	msg := nft.NewMsgTransferNFT(sender, recipient, r.Denom, r.ID, r.TokenURI)
	stdSignMsg := buildStdSignMsg(c.chainId, memo, accInfo, msg)

	return signAndBroadcastTx(c, stdSignMsg, commit)
}
