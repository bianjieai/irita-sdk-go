package types

import (
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
)

type NetworkType int

const (
	_ NetworkType = iota
	Testnet
	Mainnet
)

type (
	MsgBankSend = bank.MsgSend
	StdFee      = auth.StdFee
)
