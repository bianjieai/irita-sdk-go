package integration_test

import (
	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/v2/types"

	"github.com/bianjieai/irita-sdk-go/v2/modules/perm"
)

func (s IntegrationTestSuite) TestPerm() {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	acc := s.GetRandAccount()
	roles := []perm.Role{
		perm.RoleBlacklistAdmin,
	}

	//test AddRoles
	rs, err := s.Perm.AssignRoles(acc.Address.String(), roles, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// test QueryRoles
	roles2, err := s.Perm.QueryRoles(acc.Address.String())
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), roles2)
	require.EqualValues(s.T(), roles, roles2)

	// test RemoveRoles
	rs, err = s.Perm.UnassignRoles(acc.Address.String(), roles, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// test QueryRoles again
	roles2, err = s.Perm.QueryRoles(acc.Address.String())
	require.NoError(s.T(), err)
	require.Empty(s.T(), roles2)

	// test BlockAccount
	rs, err = s.Perm.BlockAccount(acc.Address.String(), baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// test QueryBlacklist
	bl, err := s.Perm.QueryBlacklist(1, 10)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), bl)
	require.EqualValues(s.T(), []string{acc.Address.String()}, bl)

	// test UnblockAccount
	rs, err = s.Perm.UnblockAccount(acc.Address.String(), baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	// test QueryBlacklist again
	bl, err = s.Perm.QueryBlacklist(1, 10)
	require.NoError(s.T(), err)
	require.Empty(s.T(), bl)
}
