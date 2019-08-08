package lcd

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/sdk-go/client/types"
)

type (
	AccountInfo struct {
		Type  string           `json:"type"`
		Value AccountInfoValue `json:"value"`
	}

	AccountInfoValue struct {
		AccountNumber string       `json:"account_number"`
		Address       string       `json:"address"`
		Sequence      string       `json:"sequence"`
		Coins         []types.Coin `json:"coins"`
	}
)

func (c *client) QueryAccount(address string) (AccountInfo, error) {
	var (
		accountInfo AccountInfo
	)
	path := fmt.Sprintf(UriQueryAccount, address)

	if _, body, err := c.httpClient.Get(path, nil); err != nil {
		return accountInfo, err
	} else {
		if err := json.Unmarshal(body, &accountInfo); err != nil {
			return accountInfo, err
		} else {
			return accountInfo, nil
		}
	}
}
