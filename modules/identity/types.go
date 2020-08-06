package identity

import (
	"encoding/json"
	"errors"
	"fmt"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// Identity message types and params
const (
	TypeMsgCreateIdentity = "create_identity" // type for MsgCreateIdentity
	TypeMsgUpdateIdentity = "update_identity" // type for MsgUpdateIdentity

	IDLength     = 16  // size of the ID in bytes
	MaxURILength = 140 // maximum size of the URI

	// ModuleName is the name of the identity module
	ModuleName = "identity"

	doNotModifyDesc = "[do-not-modify]" // description used to indicate not to modify a field
)

var (
	_ sdk.Msg = MsgCreateIdentity{}
	_ sdk.Msg = MsgUpdateIdentity{}

	amino     = codec.New()
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	registerCodec(amino)
}

// Route implements Msg
func (msg MsgCreateIdentity) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgCreateIdentity) Type() string { return TypeMsgCreateIdentity }

// ValidateBasic implements Msg
func (msg MsgCreateIdentity) ValidateBasic() error {
	return ValidateIdentityFields(
		msg.ID,
		msg.PubKey,
		msg.Certificate,
		msg.Credentials,
		msg.Owner,
	)
}

// GetSignBytes implements Msg
func (msg MsgCreateIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgCreateIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// Route implements Msg.
func (msg MsgUpdateIdentity) Route() string {
	return ModuleName
}

// Type implements Msg.
func (msg MsgUpdateIdentity) Type() string {
	return TypeMsgUpdateIdentity
}

// GetSignBytes implements Msg.
func (msg MsgUpdateIdentity) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateIdentity) ValidateBasic() error {
	return ValidateIdentityFields(
		msg.ID,
		msg.PubKey,
		msg.Certificate,
		msg.Credentials,
		msg.Owner,
	)
}

// GetSigners implements Msg.
func (msg MsgUpdateIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateIdentityFields validates the given identity fields
func ValidateIdentityFields(
	id tmbytes.HexBytes,
	pubKey *PubKeyInfo,
	certificate,
	credentials string,
	owner sdk.AccAddress,
) error {
	if owner.Empty() {
		return errors.New("owner missing")
	}

	if len(id) != IDLength {
		return fmt.Errorf("size of the ID must be %d in bytes", IDLength)
	}

	if len(credentials) > MaxURILength {
		return fmt.Errorf("length of the credentials uri must not be greater than %d", MaxURILength)
	}

	return nil
}

func (m Identity) Convert() interface{} {
	var pubKeyInfos []PubkeyInfo
	for _, info := range m.PubKeys {
		pubKeyInfos = append(
			pubKeyInfos,
			PubkeyInfo{
				PubKey:     info.PubKey.String(),
				PubKeyAlgo: info.Algorithm,
			},
		)
	}

	return QueryIdentityResponse{
		ID:           m.ID.String(),
		PubkeyInfos:  pubKeyInfos,
		Certificates: m.Certificates,
		Credentials:  m.Credentials,
		Owner:        m.Owner.String(),
	}
}

// MarshalJSON returns the JSON representation
func (p PubKeyAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(PubKeyAlgorithm_name[int32(p)])
}

// UnmarshalJSON unmarshals raw JSON bytes into a PubKeyAlgorithm
func (p *PubKeyAlgorithm) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}

	algo := PubKeyAlgorithm_value[s]
	*p = PubKeyAlgorithm(algo)
	return nil
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateIdentity{}, "irita/MsgCreateIdentity", nil)
	cdc.RegisterConcrete(MsgUpdateIdentity{}, "irita/MsgUpdateIdentity", nil)
}
