package identity

import (
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// Identity message types and params
const (
	TypeMsgCreateIdentity = "create_identity" // type for MsgCreateIdentity
	TypeMsgUpdateIdentity = "update_identity" // type for MsgUpdateIdentity

	IDLength     = 16  // size of the ID in bytes
	MaxURILength = 140 // maximum size of the URI

	DoNotModifyDesc = "[do-not-modify]" // description used to indicate not to modify a field

	ModuleName = "identity"

	RouterKey = ModuleName
)

var (
	_ sdk.Msg = &MsgCreateIdentity{}
	_ sdk.Msg = &MsgUpdateIdentity{}
)

// Route implements Msg
func (m MsgCreateIdentity) Route() string { return RouterKey }

// Type implements Msg
func (m MsgCreateIdentity) Type() string { return TypeMsgCreateIdentity }

// ValidateBasic implements Msg
func (m MsgCreateIdentity) ValidateBasic() error {
	return ValidateIdentityFields(
		m.Id,
		m.PubKey,
		m.Certificate,
		m.Credentials,
		m.Owner,
	)
}

// GetSignBytes implements Msg
func (m MsgCreateIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (m MsgCreateIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Owner)}
}

// Route implements m.
func (m MsgUpdateIdentity) Route() string { return RouterKey }

// Type implements m.
func (m MsgUpdateIdentity) Type() string { return TypeMsgUpdateIdentity }

// GetSignBytes implements m.
func (m MsgUpdateIdentity) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements m.
func (m MsgUpdateIdentity) ValidateBasic() error {
	return ValidateIdentityFields(
		m.Id,
		m.PubKey,
		m.Certificate,
		m.Credentials,
		m.Owner,
	)
}

// GetSigners implements m.
func (m MsgUpdateIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Owner)}
}

// ValidateIdentityFields validates the given identity fields
func ValidateIdentityFields(
	id string,
	pubKey *PubKeyInfo,
	certificate,
	credentials string,
	owner string,
) error {
	if len(owner) == 0 {
		return errors.New("owner missing")
	}

	if len(id) != IDLength*2 {
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
		pubKeyInfos = append(pubKeyInfos, PubkeyInfo{
			PubKey:     info.PubKey,
			PubKeyAlgo: info.Algorithm,
		})
	}

	return QueryIdentityResp{
		ID:           m.Id,
		PubkeyInfos:  pubKeyInfos,
		Certificates: m.Certificates,
		Credentials:  m.Credentials,
		Owner:        m.Owner,
		Data:         m.Data,
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
