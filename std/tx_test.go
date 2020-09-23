package std_test

import (
	"testing"

	"github.com/bianjieai/irita-sdk-go/modules/bank"
	"github.com/bianjieai/irita-sdk-go/std"
	sdk "github.com/bianjieai/irita-sdk-go/types"

	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	f := std.NewStdFee(100, sdk.NewCoins(sdk.NewInt64Coin("stake", 50)))
	m := "hello world"
	acc1 := sdk.AccAddress("from")
	acc2 := sdk.AccAddress("to")
	msg1 := bank.NewMsgSend(acc1, acc2, sdk.NewCoins(sdk.NewInt64Coin("stake", 100000)))
	sdkMsgs := []sdk.Msg{&msg1}

	tx, err := std.NewTransaction(f, m, sdkMsgs)
	require.NoError(t, err)
	require.NotNil(t, tx)
	require.Equal(t, tx.GetMsgs(), sdkMsgs)
	require.Equal(t, tx.GetSigners(), []sdk.AccAddress{acc1})
	require.Equal(t, tx.GetFee(), f)
	require.Equal(t, tx.GetMemo(), m)

	// no signatures; validation should fail
	require.Empty(t, tx.GetSignatures())
	require.Error(t, tx.ValidateBasic())

	signDocJSON := `{"base":{"accountNumber":"1","chainId":"chain-test","fee":{"amount":[{"amount":"50","denom":"stake"}],"gas":"100"},"memo":"hello world","sequence":"21"},"msgs":[{"msgSend":{"amount":[{"amount":"100000","denom":"stake"}],"fromAddress":"iaa1veex7mgwf54pe","toAddress":"iaa1w3hsdwlc9j"}}]}`
	bz, err := tx.CanonicalSignBytes("chain-test", 1, 21)
	require.NoError(t, err)
	require.Equal(t, signDocJSON, string(bz))
}
