package admin

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	// ModuleName is the name of the admin module
	ModuleName = "admin"
)

var (
	_ sdk.Msg = MsgAddRoles{}
	_ sdk.Msg = MsgRemoveRoles{}
	_ sdk.Msg = MsgBlockAccount{}
	_ sdk.Msg = MsgUnblockAccount{}

	amino = codec.New()

	// ModuleCdc references the global admin module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to admin and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	registerCodec(amino)
}

func (m MsgAddRoles) Route() string {
	return ModuleName
}

func (m MsgAddRoles) Type() string {
	return "add_roles"
}

func (m MsgAddRoles) ValidateBasic() error {
	if m.Address.Empty() {
		return errors.New("address missing")
	}
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	if len(m.Roles) == 0 {
		return errors.New("roles missing")
	}
	return nil
}

func (m MsgAddRoles) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgAddRoles) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

func (m MsgRemoveRoles) Route() string {
	return ModuleName
}

func (m MsgRemoveRoles) Type() string {
	return "remove_roles"
}

func (m MsgRemoveRoles) ValidateBasic() error {
	if m.Address.Empty() {
		return errors.New("address missing")
	}
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	if len(m.Roles) == 0 {
		return errors.New("roles missing")
	}
	return nil
}

func (m MsgRemoveRoles) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgRemoveRoles) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

func (m MsgBlockAccount) Route() string {
	return ModuleName
}

func (m MsgBlockAccount) Type() string {
	return "block_account"
}

func (m MsgBlockAccount) ValidateBasic() error {
	if m.Address.Empty() {
		return errors.New("address missing")
	}
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	return nil
}

func (m MsgBlockAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgBlockAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

func (m MsgUnblockAccount) Route() string {
	return ModuleName
}

func (m MsgUnblockAccount) Type() string {
	return "unblock_account"
}

func (m MsgUnblockAccount) ValidateBasic() error {
	if m.Address.Empty() {
		return errors.New("address missing")
	}
	if m.Operator.Empty() {
		return errors.New("operator missing")
	}
	return nil
}

func (m MsgUnblockAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUnblockAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Operator}
}

// RoleFromstring turns a string into a Auth
func roleFromString(str string) (Role, error) {
	switch str {
	case "RootAdmin":
		return RoleRootAdmin, nil

	case "PermAdmin":
		return RolePermAdmin, nil

	case "BlacklistAdmin":
		return RoleBlacklistAdmin, nil

	case "NodeAdmin":
		return RoleNodeAdmin, nil

	case "ParamAdmin":
		return RoleParamAdmin, nil

	case "PowerUser":
		return RolePowerUser, nil

	case "RelayerUser":
		return RoleRelayerUser, nil

	case "RoleIDAdmin":
		return RoleIDAdmin, nil
	default:
		return Role(0xff), fmt.Errorf("'%s' is not a valid role", str)
	}
}

// Marshal needed for protobuf compatibility
func (r Role) Marshal() ([]byte, error) {
	return []byte{byte(r)}, nil
}

// Unmarshal needed for protobuf compatibility
func (r *Role) Unmarshal(data []byte) error {
	*r = Role(data[0])
	return nil
}

// MarshalJSON Marshals to JSON using string representation of the status
func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.string())
}

// UnmarshalJSON Unmarshals from JSON assuming Bech32 encoding
func (r *Role) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := roleFromString(s)
	if err != nil {
		return err
	}

	*r = bz2
	return nil
}

// string implements the stringer interface.
func (r Role) string() string {
	switch r {
	case RoleRootAdmin:
		return "RootAdmin"

	case RolePermAdmin:
		return "PermAdmin"

	case RoleBlacklistAdmin:
		return "BlacklistAdmin"

	case RoleNodeAdmin:
		return "NodeAdmin"

	case RoleParamAdmin:
		return "ParamAdmin"

	case RolePowerUser:
		return "PowerUser"

	case RoleRelayerUser:
		return "RelayerUser"

	case RoleIDAdmin:
		return "RoleIDAdmin"

	default:
		return ""
	}
}

func registerCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddRoles{}, "irita/modules/MsgAddRoles", nil)
	cdc.RegisterConcrete(MsgRemoveRoles{}, "irita/modules/MsgRemoveRoles", nil)
	cdc.RegisterConcrete(MsgBlockAccount{}, "irita/modules/MsgBlockAccount", nil)
	cdc.RegisterConcrete(MsgUnblockAccount{}, "irita/modules/MsgUnblockAccount", nil)
}
