package integration_test

import (
	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/types"

	"github.com/bianjieai/irita-sdk-go/modules/admin"

)

func (s IntegrationTestSuite) TestAdmin() {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	acc := s.GetRandAccount()
	roles := []admin.Role{
		admin.RoleBlacklistAdmin,
	}

	adminClient := admin.NewClient(s.IRITAClient.BaseClient, s.IRITAClient.AppCodec())
	s.IRITAClient.RegisterModule(adminClient)

	//test AddRoles
	rs, err := adminClient.AddRoles(acc.Address.String(), roles, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// test QueryRoles
	roles2, err := adminClient.QueryRoles(acc.Address.String())
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), roles2)
	require.EqualValues(s.T(), roles, roles2)

	// test RemoveRoles
	rs, err = adminClient.RemoveRoles(acc.Address.String(), roles, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// test QueryRoles again
	roles2, err = adminClient.QueryRoles(acc.Address.String())
	require.NoError(s.T(), err)
	require.Empty(s.T(), roles2)

	// test BlockAccount
	rs, err = adminClient.BlockAccount(acc.Address.String(), baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// test QueryBlacklist
	bl, err := adminClient.QueryBlacklist(1, 10)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), bl)
	require.EqualValues(s.T(), []string{acc.Address.String()}, bl)

	// test UnblockAccount
	rs, err = adminClient.UnblockAccount(acc.Address.String(), baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// test QueryBlacklist again
	bl, err = adminClient.QueryBlacklist(1, 10)
	require.NoError(s.T(), err)
	require.Empty(s.T(), bl)
}
