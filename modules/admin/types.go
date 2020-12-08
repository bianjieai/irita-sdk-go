package admin

import (
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	// ModuleName is the name of the admin module
	ModuleName = "admin"
)

var (
	_ sdk.Msg = &MsgAddRoles{}
	_ sdk.Msg = &MsgRemoveRoles{}
	_ sdk.Msg = &MsgBlockAccount{}
	_ sdk.Msg = &MsgUnblockAccount{}
)

func (m MsgAddRoles) Route() string {
	return ModuleName
}

func (m MsgAddRoles) Type() string {
	return "add_roles"
}

func (m MsgAddRoles) ValidateBasic() error {
	if len(m.Address) == 0 {
		return errors.New("address missing")
	}
	if len(m.Operator) == 0 {
		return errors.New("operator missing")
	}
	if len(m.Roles) == 0 {
		return errors.New("roles missing")
	}
	return nil
}

func (m MsgAddRoles) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgAddRoles) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
}

func (m MsgRemoveRoles) Route() string {
	return ModuleName
}

func (m MsgRemoveRoles) Type() string {
	return "remove_roles"
}

func (m MsgRemoveRoles) ValidateBasic() error {
	if len(m.Address) == 0 {
		return errors.New("address missing")
	}
	if len(m.Operator) == 0 {
		return errors.New("operator missing")
	}
	if len(m.Roles) == 0 {
		return errors.New("roles missing")
	}
	return nil
}

func (m MsgRemoveRoles) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgRemoveRoles) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
}

func (m MsgBlockAccount) Route() string {
	return ModuleName
}

func (m MsgBlockAccount) Type() string {
	return "block_account"
}

func (m MsgBlockAccount) ValidateBasic() error {
	if len(m.Address) == 0 {
		return errors.New("address missing")
	}
	if len(m.Operator) == 0 {
		return errors.New("operator missing")
	}
	return nil
}

func (m MsgBlockAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgBlockAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
}

func (m MsgUnblockAccount) Route() string {
	return ModuleName
}

func (m MsgUnblockAccount) Type() string {
	return "unblock_account"
}

func (m MsgUnblockAccount) ValidateBasic() error {
	if len(m.Address) == 0 {
		return errors.New("address missing")
	}
	if len(m.Operator) == 0 {
		return errors.New("operator missing")
	}
	return nil
}

func (m MsgUnblockAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgUnblockAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Operator)}
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
