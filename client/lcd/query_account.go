package lcd

import (
	"encoding/json"
	"fmt"
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/types"
)

type (
	AccQueryResult struct {
		Height string      `json:"height"`
		Result AccountInfo `json:"result"`
	}
	AccountInfo struct {
		Type  string           `json:"type"`
		Value AccountInfoValue `json:"value"`
	}

	AccountInfoValue struct {
		AccountNumber uint64       `json:"account_number"`
		Address       string       `json:"address"`
		Sequence      uint64       `json:"sequence"`
		Coins         []types.Coin `json:"coins"`
	}
)

func (c *client) QueryAccount(address string) (AccountInfo, error) {
	var (
		accQueryResult AccQueryResult
		accountInfo    AccountInfo
	)
	path := fmt.Sprintf(UriQueryAccount, address)

	if _, body, err := c.httpClient.Get(path, nil); err != nil {
		return accountInfo, err
	} else {
		if err := json.Unmarshal(body, &accQueryResult); err != nil {
			return accountInfo, err
		} else {
			return accQueryResult.Result, nil
		}
	}
}
