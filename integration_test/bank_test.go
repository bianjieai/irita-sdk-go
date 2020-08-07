package integration_test

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/modules/bank"
	"github.com/bianjieai/irita-sdk-go/types"
)

func (s IntegrationTestSuite) TestBank() {
	cases := []SubTest{
		{
			"TestQueryAccount",
			queryAccount,
		},
		{
			"TestSend",
			send,
		},
		{
			"TestMultiSend",
			multiSend,
		},
		{
			"TestSimulate",
			simulate,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() {
			t.testCase(s)
		})
	}
}

func queryAccount(s IntegrationTestSuite) {
	acc, err := s.Bank.QueryAccount(s.Account().Address.String())
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), acc)
}

func send(s IntegrationTestSuite) {
	coins, err := types.ParseDecCoins("10point")
	s.NoError(err)
	to := s.GetRandAccount().Address.String()
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	ch := make(chan int)
	s.Bank.SubscribeSendTx(s.Account().Address.String(), to, func(send bank.EventDataMsgSend) {
		ch <- 1
	})
	result, err := s.Bank.Send(to, coins, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), result.Hash)

	resp, err := s.QueryTx(result.Hash)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.Result.Code, uint32(0))
	require.Equal(s.T(), resp.Height, result.Height)

	<-ch
}

func multiSend(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      800000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	coins, e := types.ParseDecCoins("1000point")
	require.NoError(s.T(), e)

	var accNum = 11
	var acc = make([]string, accNum)
	var receipts = make([]bank.Receipt, accNum)
	for i := 0; i < accNum; i++ {
		acc[i] = s.RandStringOfLength(10)
		addr, _, err := s.Key.Add(acc[i], "1234567890")

		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), addr)

		receipts[i] = bank.Receipt{
			Address: addr,
			Amount:  coins,
		}
	}

	_, err := s.Bank.MultiSend(bank.MultiSendRequest{Receipts: receipts}, baseTx)
	require.NoError(s.T(), err)

	coins, e = types.ParseDecCoins("1point")
	require.NoError(s.T(), e)

	to := s.GetRandAccount().Address.String()

	begin := time.Now()
	var wait sync.WaitGroup
	for i := 1; i <= 50; i++ {
		wait.Add(1)
		index := rand.Intn(accNum)
		go func() {
			defer wait.Done()
			_, err := s.Bank.Send(to, coins, types.BaseTx{
				From:     acc[index],
				Gas:      200000,
				Memo:     "test",
				Mode:     types.Async,
				Password: "1234567890",
			})
			require.NoError(s.T(), err)
		}()
	}
	wait.Wait()
	end := time.Now()
	fmt.Printf("total senconds:%s\n", end.Sub(begin).String())
}

func simulate(s IntegrationTestSuite) {
	coins, err := types.ParseDecCoins("10point")
	s.NoError(err)
	to := s.GetRandAccount().Address.String()
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
		Simulate: true,
	}

	result, err := s.Bank.Send(to, coins, baseTx)
	require.NoError(s.T(), err)
	require.Greater(s.T(), result.GasWanted, int64(0))
	fmt.Println(result)
}
