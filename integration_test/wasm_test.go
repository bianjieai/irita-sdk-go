package integration_test

import (
	"fmt"

	"github.com/bianjieai/irita-sdk-go/modules/wasm"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func (its IntegrationTestSuite) TestWasm() {
	baseTx := types.BaseTx{
		From:     its.Account().Name,
		Gas:      4000000,
		Fee:      types.NewDecCoins(types.NewInt64DecCoin("point", 24)),
		Memo:     "test",
		Mode:     types.Commit,
		Password: its.Account().Password,
	}

	request := wasm.StoreRequest{
		WASMFile: "./election.wasm",
	}

	codeID, err := its.WASM.Store(request, baseTx)
	require.NoError(its.T(), err)
	require.NotEmpty(its.T(), codeID)

	args := wasm.NewArgs().
		Put("start", 1).
		Put("end", 100).
		Put("candidates", []string{"iaa1qvty8x0c78am8c44zv2n7tgm6gfqt78j0verqa", "iaa1zk2tse0pkk87p2v8tcsfs0ytfw3t88kejecye5"})

	initReq := wasm.InstantiateRequest{
		CodeID:  codeID,
		Label:   "test wasm",
		InitMsg: args,
	}

	contractAddress, err := its.WASM.Instantiate(initReq, baseTx)
	require.NoError(its.T(), err)
	require.NotEmpty(its.T(), contractAddress)

	info, err := its.WASM.QueryContractInfo(contractAddress)
	require.NoError(its.T(), err)
	require.NotEmpty(its.T(), info)
	require.Equal(its.T(), fmt.Sprintf("%d", info.CodeID), codeID)
	require.Equal(its.T(), info.Label, initReq.Label)

	execAbi := wasm.NewContractABI().
		WithMethod("vote").
		WithArgs("candidate", "iaa1qvty8x0c78am8c44zv2n7tgm6gfqt78j0verqa")
	_, err = its.WASM.Execute(contractAddress, execAbi, nil, baseTx)
	require.NoError(its.T(), err)

	queryAbi := wasm.NewContractABI().
		WithMethod("get_vote_info")
	bz, err := its.WASM.QueryContract(contractAddress, queryAbi)
	require.NoError(its.T(), err)
	require.NotEmpty(its.T(), bz)
}
