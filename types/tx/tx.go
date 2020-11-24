package tx

import (
	"fmt"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

func Unwrap(tx sdk.Tx) (*sdk.UnwrappedTx, error) {
	txWrapper, ok := tx.(*wrapper)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, tx)
	}

	txBody := sdk.UnwrappedTxBody{
		Msgs:          txWrapper.tx.GetMsgs(),
		Memo:          txWrapper.tx.Body.Memo,
		TimeoutHeight: txWrapper.tx.Body.TimeoutHeight,
	}

	var signatures []sdk.UnwrappedSignature
	for i, pubKey := range txWrapper.GetPubKeys() {
		bech32PubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sdk.UnwrappedSignature{
			PubKey: sdk.UnwrappedPubKey{
				Type:  pubKey.Type(),
				Value: bech32PubKey,
			},
			Sig:      txWrapper.tx.Signatures[i],
			Sequence: txWrapper.tx.AuthInfo.SignerInfos[i].Sequence,
		})
	}

	authInfo := sdk.UnwrappedAuthInfo{
		Signatures: signatures,
		Fee: &sdk.StdFee{
			Amount: txWrapper.tx.AuthInfo.Fee.Amount,
			Gas:    txWrapper.tx.AuthInfo.Fee.GasLimit,
		},
	}

	return &sdk.UnwrappedTx{
		Body:     txBody,
		AuthInfo: authInfo,
	}, nil
}
