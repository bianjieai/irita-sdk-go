package modules

import (
	"errors"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func (base baseClient) QueryWithData(path string, key []byte) ([]byte, int64, error) {
	opts := rpcclient.ABCIQueryOptions{
		Prove: true,
	}

	result, err := base.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return nil, 0, err
	}

	if !result.Response.IsOK() {
		return nil, 0, errors.New(result.Response.Log)
	}
	return result.Response.Value, result.Response.Height, nil
}
