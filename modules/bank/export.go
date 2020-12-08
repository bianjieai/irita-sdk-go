package bank

import (
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

// Client export a group api for Admin module
type Client interface {
	sdk.Module

	Send(to string, amount sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	MultiSend(receipts MultiSendRequest, baseTx sdk.BaseTx) ([]sdk.ResultTx, sdk.Error)
	SubscribeSendTx(from, to string, callback EventMsgSendCallback) sdk.Subscription

	QueryAccount(address string) (sdk.BaseAccount, sdk.Error)
}

type Receipt struct {
	Address string       `json:"address"`
	Amount  sdk.DecCoins `json:"amount"`
}

type MultiSendRequest struct {
	Receipts []Receipt
}

func (msr MultiSendRequest) Len() int {
	return len(msr.Receipts)
}

func (msr MultiSendRequest) Sub(begin, end int) sdk.SplitAble {
	return MultiSendRequest{Receipts: msr.Receipts[begin:end]}
}

type EventDataMsgSend struct {
	Height int64      `json:"height"`
	Hash   string     `json:"hash"`
	From   string     `json:"from"`
	To     string     `json:"to"`
	Amount []sdk.Coin `json:"amount"`
}

type EventMsgSendCallback func(EventDataMsgSend)
