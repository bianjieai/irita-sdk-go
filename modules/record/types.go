package record

import (
	"fmt"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	ModuleName = "record"

	attributeKeyRecordID  = "record_id"
	eventTypeCreateRecord = "create_record"
)

var (
	_ sdk.Msg = MsgCreateRecord{}

	amino = codec.New()

	// ModuleCdc references the global record module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to record and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	registerCodec(amino)
}

// Route implements Msg.
func (msg MsgCreateRecord) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgCreateRecord) Type() string { return "create_record" }

// GetSignBytes implements Msg.
func (msg MsgCreateRecord) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
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
	if msg.Creator.Empty() {
		return fmt.Errorf("creator missing")
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
	return []sdk.AccAddress{msg.Creator}
}

func (this Record) Convert() interface{} {
	return QueryRecordResponse{
		TxHash:   this.TxHash.String(),
		Contents: this.Contents,
		Creator:  this.Creator.String(),
	}
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateRecord{}, "irismod/record/MsgCreateRecord", nil)
}
