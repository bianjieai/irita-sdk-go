package bank

import (
	"context"
	"strings"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	utils "github.com/bianjieai/irita-sdk-go/utils"
)

type bankClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return bankClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (b bankClient) Name() string {
	return ModuleName
}

func (b bankClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

// QueryAccount return account information specified address
func (b bankClient) QueryAccount(address string) (sdk.BaseAccount, sdk.Error) {
	conn, err := b.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	if err := sdk.ValidateAccAddress(address); err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	account, err := b.BaseClient.QueryAccount(address)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	request := &QueryAllBalancesRequest{
		Address:    address,
		Pagination: nil,
	}
	balances, err := NewQueryClient(conn).AllBalances(context.Background(), request)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	account.Coins = balances.Balances
	return account, nil
}

// Send is responsible for transferring tokens from `From` to `to` account
func (b bankClient) Send(to string, amount sdk.DecCoins, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	if err := sdk.ValidateAccAddress(to); err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s not found", baseTx.From)
	}

	amt, err := b.ToMinCoin(amount...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := NewMsgSend(sender.String(), to, amt)
	return b.BuildAndSend([]sdk.Msg{&msg}, baseTx)
}

func (b bankClient) MultiSend(request MultiSendRequest, baseTx sdk.BaseTx) (resTxs []sdk.ResultTx, err sdk.Error) {
	sender, err := b.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, sdk.Wrapf("%s not found", baseTx.From)
	}

	if len(request.Receipts) > maxMsgLen {
		return b.SendBatch(sender, request, baseTx)
	}

	var inputs = make([]Input, len(request.Receipts))
	var outputs = make([]Output, len(request.Receipts))
	for i, receipt := range request.Receipts {
		if err := sdk.ValidateAccAddress(receipt.Address); err != nil {
			return nil, sdk.Wrap(err)
		}

		amt, err := b.ToMinCoin(receipt.Amount...)
		if err != nil {
			return nil, sdk.Wrap(err)
		}

		inputs[i] = NewInput(sender.String(), amt)
		outputs[i] = NewOutput(receipt.Address, amt)
	}

	msg := NewMsgMultiSend(inputs, outputs)
	res, err := b.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resTxs = append(resTxs, res)
	return
}

func (b bankClient) SendBatch(sender sdk.AccAddress,
	request MultiSendRequest, baseTx sdk.BaseTx) ([]sdk.ResultTx, sdk.Error) {
	batchReceipts := utils.SubArray(maxMsgLen, request)

	var msgs sdk.Msgs
	for _, receipts := range batchReceipts {
		req := receipts.(MultiSendRequest)
		var inputs = make([]Input, len(req.Receipts))
		var outputs = make([]Output, len(req.Receipts))
		for i, receipt := range req.Receipts {
			if err := sdk.ValidateAccAddress(receipt.Address); err != nil {
				return nil, sdk.Wrap(err)
			}

			amt, err := b.ToMinCoin(receipt.Amount...)
			if err != nil {
				return nil, sdk.Wrap(err)
			}

			inputs[i] = NewInput(sender.String(), amt)
			outputs[i] = NewOutput(receipt.Address, amt)
		}
		msgs = append(msgs, NewMsgMultiSend(inputs, outputs))
	}
	return b.BaseClient.SendBatch(msgs, baseTx)
}

// SubscribeSendTx Subscribe MsgSend event and return subscription
func (b bankClient) SubscribeSendTx(from, to string, callback EventMsgSendCallback) sdk.Subscription {
	var builder = sdk.NewEventQueryBuilder()

	from = strings.TrimSpace(from)
	if len(from) != 0 {
		builder.AddCondition(sdk.NewCond(sdk.EventTypeMessage,
			sdk.AttributeKeySender).EQ(sdk.EventValue(from)))
	}

	to = strings.TrimSpace(to)
	if len(to) != 0 {
		builder.AddCondition(sdk.Cond("transfer.recipient").EQ(sdk.EventValue(to)))
	}

	subscription, _ := b.SubscribeTx(builder, func(data sdk.EventDataTx) {
		for _, msg := range data.Tx.GetMsgs() {
			if value, ok := msg.(*MsgSend); ok {
				callback(EventDataMsgSend{
					Height: data.Height,
					Hash:   data.Hash,
					From:   value.FromAddress,
					To:     value.ToAddress,
					Amount: value.Amount,
				})
			}
		}
	})
	return subscription
}
