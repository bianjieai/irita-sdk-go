// Package bank is mainly used to transfer coins between accounts,query account balances, and implement interface rpc.BankI
//
//
// As a quick start:
//
// TransferNFT coins to other account
//
//  client := test.NewClient()
//  amt := types.NewIntWithDecimal(1, 18)
//  coins := types.NewCoins(types.NewCoin("point", amt))
//  to := "caa1rgnu8grzt6mwnjg7jss7w0sfyjn67g4em9njf5"
//  baseTx := types.BaseTx{
// 		From: "test1",
// 		Gas:  20000,
// 		Memo: "test",
// 		Mode: types.Commit,
//  }
//  result,err := client.BankI.Send(to,coins,baseTx)
//  fmt.Println(result)
//
// BurnNFT some coins from your account
//
//  client := test.NewClient()
//  amt := types.NewIntWithDecimal(1, 18)
//  coins := types.NewCoins(types.NewCoin("point", amt))
//  baseTx := types.BaseTx{
// 		From: "test1",
// 		Gas:  20000,
// 		Memo: "test",
// 		Mode: types.Commit,
//  }
//  result,err := client.BankI.BurnNFT(coins, baseTx)
//  fmt.Println(result)
//
// Set account memo
//
//  client := test.NewClient()
//  result,err := client.BankI.SetMemoRegexp("testMemo", baseTx)
//  fmt.Println(result)
//
// Queries account information
//
//  client := test.NewClient()
//  result,err := client.BankI.QueryAccount("caa1rgnu8grzt6mwnjg7jss7w0sfyjn67g4em9njf5")
//  fmt.Println(result)
//
// Queries the token information
//
//  client := test.NewClient()
//  result,err := client.BankI.QueryTokenStats("point")
//  fmt.Println(result)
//
package bank
