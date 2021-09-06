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
			"TestConcurrency",
			concurrencySend,
		},
		{
			"TestMultiSend",
			multiSend,
		},
		{
			"TestSendBatch",
			sendBatch,
		},
		{
			"buildAndSigned",
			buildAndSigned,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() {
			t.testCase(s)
		})
	}
}

func queryAccount(s IntegrationTestSuite) {
	acc, err := s.Bank.QueryAccount("iaa1f7jnpkt2yxd8h72w6hf03juy3uk6m2sur845kq")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), acc)
}

func send(s IntegrationTestSuite) {
	coins, err := types.ParseDecCoins("10point")
	s.NoError(err)
	to := s.GetRandAccount().Address.String()
	baseTx := types.BaseTx{
		From:               s.Account().Name,
		Gas:                200000,
		Memo:               "test",
		Mode:               types.Commit,
		Password:           s.Account().Password,
		SimulateAndExecute: false,
		GasAdjustment:      1.5,
	}

	ch := make(chan int)
	s.Bank.SubscribeSendTx(s.Account().Address.String(), to, func(send bank.EventDataMsgSend) {
		ch <- 1
	})
	result, err := s.Bank.Send(to, coins, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), result.Hash)
	time.Sleep(1 * time.Second)

	resp, err := s.Manager().QueryTx(result.Hash)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.Result.Code, uint32(0))
	require.Equal(s.T(), resp.Height, result.Height)

	<-ch
}

func concurrencySend(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      500000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	coins, err := types.ParseDecCoins("1point")
	s.NoError(err)
	to := s.GetRandAccount().Address.String()

	var lock sync.WaitGroup
	for i := 1; i <= 20; i++ {
		lock.Add(1)
		go func() {
			result, err := s.Bank.Send(to, coins, baseTx)
			require.NoError(s.T(), err)
			require.NotEmpty(s.T(), result.Hash)
			lock.Done()
		}()
	}
	lock.Wait()
}

func multiSend(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      500000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	coins, e := types.ParseDecCoins("1000point")
	require.NoError(s.T(), e)

	var accNum = 4
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

func sendBatch(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:               s.Account().Name,
		Gas:                500000,
		Memo:               "test",
		Mode:               types.Commit,
		Password:           s.Account().Password,
		SimulateAndExecute: true,
	}

	coins, e := types.ParseDecCoins("1point")
	require.NoError(s.T(), e)

	var accNum = 100
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
}

func buildAndSigned(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      500000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}
	amount, err := types.ParseCoin("1000uirita")
	require.NoError(s.T(), err)
	sendMsg := bank.MsgSend{
		FromAddress: "",
		ToAddress:   "",
		Amount:      types.Coins{amount},
	}
	res, err := s.BaseClient.BuildAndSign([]types.Msg{&sendMsg}, baseTx)
	require.NoError(s.T(), err)
	fmt.Println(string(res))
}
