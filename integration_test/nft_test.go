package integration_test

import (
	"fmt"
	"strings"

	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/modules/nft"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

func (s IntegrationTestSuite) TestNFT() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	denomID := strings.ToLower(s.RandStringOfLength(4))
	denomName := strings.ToLower(s.RandStringOfLength(4))
	schema := strings.ToLower(s.RandStringOfLength(10))
	issueReq := nft.IssueDenomRequest{
		ID:     denomID,
		Name:   denomName,
		Schema: schema,
	}
	res, err := s.NFT.IssueDenom(issueReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	tokenID := strings.ToLower(s.RandStringOfLength(7))
	tokenName := strings.ToLower(s.RandStringOfLength(7))
	tokenData := strings.ToLower(s.RandStringOfLength(7))
	mintReq := nft.MintNFTRequest{
		Denom: denomID,
		ID:    tokenID,
		Name:  tokenName,
		URI:   fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
		Data:  tokenData,
	}
	res, err = s.NFT.MintNFT(mintReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	editReq := nft.EditNFTRequest{
		Denom: mintReq.Denom,
		ID:    mintReq.ID,
		URI:   fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
	}
	res, err = s.NFT.EditNFT(editReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	nftRes, err := s.NFT.QueryNFT(mintReq.Denom, mintReq.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), editReq.URI, nftRes.URI)

	supply, err := s.NFT.QuerySupply(mintReq.Denom, nftRes.Creator)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(1), supply)

	owner, err := s.NFT.QueryOwner(nftRes.Creator, mintReq.Denom)
	require.NoError(s.T(), err)
	require.Len(s.T(), owner.IDCs, 1)
	require.Len(s.T(), owner.IDCs[0].TokenIDs, 1)
	require.Equal(s.T(), tokenID, owner.IDCs[0].TokenIDs[0])

	uName := s.RandStringOfLength(10)
	pwd := "11111111"

	recipient, _, err := s.Key.Add(uName, pwd)
	require.NoError(s.T(), err)

	transferReq := nft.TransferNFTRequest{
		Recipient: recipient,
		Denom:     mintReq.Denom,
		ID:        mintReq.ID,
		URI:       fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
	}
	res, err = s.NFT.TransferNFT(transferReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	owner, err = s.NFT.QueryOwner(transferReq.Recipient, mintReq.Denom)
	require.NoError(s.T(), err)
	require.Len(s.T(), owner.IDCs, 1)
	require.Len(s.T(), owner.IDCs[0].TokenIDs, 1)
	require.Equal(s.T(), tokenID, owner.IDCs[0].TokenIDs[0])

	supply, err = s.NFT.QuerySupply(mintReq.Denom, transferReq.Recipient)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(1), supply)

	denoms, err := s.NFT.QueryDenoms()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), denoms)

	d, err := s.NFT.QueryDenom(denomID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), denomID, d.ID)
	require.Equal(s.T(), denomName, d.Name)
	require.Equal(s.T(), schema, d.Schema)

	col, err := s.NFT.QueryCollection(denomID)
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), d, col.Denom)
	require.Len(s.T(), col.NFTs, 1)
	require.Equal(s.T(), mintReq.ID, col.NFTs[0].ID)

	burnReq := nft.BurnNFTRequest{
		Denom: mintReq.Denom,
		ID:    mintReq.ID,
	}

	amount, e := sdk.ParseDecCoins("10point")
	require.NoError(s.T(), e)
	_, err = s.Bank.Send(recipient, amount, baseTx)
	require.NoError(s.T(), err)

	baseTx.From = uName
	baseTx.Password = pwd
	res, err = s.NFT.BurnNFT(burnReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	supply, err = s.NFT.QuerySupply(mintReq.Denom, transferReq.Recipient)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(0), supply)
}
