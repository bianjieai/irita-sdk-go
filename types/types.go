package types

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type NetworkType int

const (
	_ NetworkType = iota
	Testnet
	Mainnet
)

type (
	MsgBankSend = bank.MsgMultiSend
	StdFee      = auth.StdFee
)
