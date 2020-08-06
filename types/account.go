package types

// BaseAccount defines the basic structure of the account
type BaseAccount struct {
	Address       AccAddress `json:"address"`
	Coins         Coins      `json:"coins"`
	PubKey        string     `json:"public_key"`
	AccountNumber uint64     `json:"account_number"`
	Sequence      uint64     `json:"sequence"`
}
