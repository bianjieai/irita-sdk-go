package params

import (
	"errors"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	ModuleName = "params"
)

var (
	_ sdk.Msg = MsgUpdateParams{}

	amino = codec.New()

	// ModuleCdc references the global params module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to params and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	registerCodec(amino)
}

func (m MsgUpdateParams) Route() string {
	return ModuleName
}

func (m MsgUpdateParams) Type() string {
	return "update_params"
}

func (m MsgUpdateParams) ValidateBasic() error {
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	return validateChanges(m.Changes)
}

func (m MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUpdateParams{}, "irita/modules/MsgUpdateParams", nil)
}

// ValidateChanges performs basic validation checks over a set of ParamChange. It
// returns an error if any ParamChange is invalid.
func validateChanges(changes []ParamChange) error {
	if len(changes) == 0 {
		return errors.New("no change params")
	}

	for _, pc := range changes {
		if len(pc.Subspace) == 0 {
			return errors.New("empty subspace")
		}
		if len(pc.Key) == 0 {
			return errors.New("empty params key")
		}
		if len(pc.Value) == 0 {
			return errors.New("empty params value")
		}
	}

	return nil
}
