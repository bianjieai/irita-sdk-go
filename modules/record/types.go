package record

import (
	"fmt"

	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
)

const (
	ModuleName = "record"

	attributeKeyRecordID  = "record_id"
	eventTypeCreateRecord = "create_record"
)

var (
	_ sdk.Msg = &MsgCreateRecord{}

	recordKey = []byte{0x01} // record key
)

// Route implements Msg.
func (msg MsgCreateRecord) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgCreateRecord) Type() string { return "create_record" }

// GetSignBytes implements Msg.
func (msg MsgCreateRecord) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgCreateRecord) ValidateBasic() error {
	if len(msg.Contents) == 0 {
		return fmt.Errorf("contents missing")
	}
	if len(msg.Creator) == 0 {
		return fmt.Errorf("creator missing")
	}

	if err := sdk.ValidateAccAddress(msg.Creator); err != nil {
		return sdk.Wrap(err)
	}

	for i, content := range msg.Contents {
		if len(content.Digest) == 0 {
			return fmt.Errorf("content[%d] digest missing", i)
		}
		if len(content.DigestAlgo) == 0 {
			return fmt.Errorf("content[%d] digest algo missing", i)
		}
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgCreateRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (this Record) Convert() interface{} {
	return QueryRecordResp{
		Record: Data{
			TxHash:   this.TxHash,
			Contents: this.Contents,
			Creator:  this.Creator,
		},
	}
}

// GetRecordKey returns record key bytes
func GetRecordKey(recordID []byte) []byte {
	return append(recordKey, recordID...)
}
