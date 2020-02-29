package tx

import (
	"github.com/bianjieai/irita-sdk-go/client/types"
	"github.com/bianjieai/irita/modules/record"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	RecordCreateReq struct {
		Contents []Content
		Creator  string
	}
	Content struct {
		Digest     string
		DigestAlgo string
		URI        string
		Meta       string
	}
)

func (c *client) CreateRecord(req RecordCreateReq, memo string, commit bool) (types.BroadcastTxResult, error) {
	var (
		result  types.BroadcastTxResult
		creator sdk.AccAddress
	)

	// get account info
	from := c.keyManager.GetAddr()
	account, err := c.liteClient.QueryAccount(from.String())
	if err != nil {
		return result, err
	}

	// build msg
	if v, err := sdk.AccAddressFromBech32(req.Creator); err != nil {
		return result, err
	} else {
		creator = v
	}
	msg := buildCreateRecordMsg(req.Contents, creator)

	stdSignMsg := buildStdSignMsg(c.chainId, memo, account, msg)

	return signAndBroadcastTx(c, stdSignMsg, commit)
}

func buildCreateRecordMsg(contents []Content, creator sdk.AccAddress) record.MsgCreateRecord {
	var msgContents []record.Content

	if len(contents) > 0 {
		for _, v := range contents {
			msgContents = append(msgContents, record.Content{
				Digest:     v.Digest,
				DigestAlgo: v.DigestAlgo,
				URI:        v.URI,
				Meta:       v.Meta,
			})
		}
	}

	return record.MsgCreateRecord{
		Contents: msgContents,
		Creator:  creator,
	}
}
