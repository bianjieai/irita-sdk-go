package integration_test

import (
	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/modules/node"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

func (s IntegrationTestSuite) TestValidator() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	cert := `-----BEGIN CERTIFICATE-----
MIICBDCCAaugAwIBAgIUOXW65eK84IOoyuQ3EP4Ukq+XevAwCgYIKoEcz1UBg3Uw
WDELMAkGA1UEBhMCQ04xDTALBgNVBAgMBHJvb3QxDTALBgNVBAcMBHJvb3QxDTAL
BgNVBAoMBHJvb3QxDTALBgNVBAsMBHJvb3QxDTALBgNVBAMMBHJvb3QwHhcNMjEw
MzE2MTAxNjI5WhcNMjIwMzE2MTAxNjI5WjBYMQswCQYDVQQGEwJDTjENMAsGA1UE
CAwEcm9vdDENMAsGA1UEBwwEcm9vdDENMAsGA1UECgwEcm9vdDENMAsGA1UECwwE
cm9vdDENMAsGA1UEAwwEcm9vdDBZMBMGByqGSM49AgEGCCqBHM9VAYItA0IABJDC
ShvKTgOt8XXHgXVyyp/iEUA/VXQO85Yc6x9jsNfBEq425ibBEFLIMe/6F2m5d7bu
06jz1DTqNBLUCw8h2ROjUzBRMB0GA1UdDgQWBBS1B/RJhawZ0qm6xv8ZAo2inUSe
JjAfBgNVHSMEGDAWgBS1B/RJhawZ0qm6xv8ZAo2inUSeJjAPBgNVHRMBAf8EBTAD
AQH/MAoGCCqBHM9VAYN1A0cAMEQCIEILyjp7cOHZwuqDp3yhXeEu3p81b/q59v6q
D0qo30+nAiBu4gilD42Ny2SjJB5r1zjLJT9LZ/vIlcyptCI5IfDQzw==
-----END CERTIFICATE-----
`

	createReq := node.CreateValidatorRequest{
		Name:        "test1",
		Certificate: cert,
		Power:       10,
		Details:     "this is a test",
	}

	rs, err := s.Node.CreateValidator(createReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	validatorID, er := rs.Events.GetValue("create_validator", "validator")
	require.NoError(s.T(), er)

	v, err := s.Node.QueryValidator(validatorID)
	require.NoError(s.T(), err)

	vs, err := s.Node.QueryValidators(nil)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), vs)

	updateReq := node.UpdateValidatorRequest{
		ID:          validatorID,
		Name:        "test2",
		Certificate: cert,
		Power:       10,
		Details:     "this is a updated test",
	}
	rs, err = s.Node.UpdateValidator(updateReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	v, err = s.Node.QueryValidator(validatorID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), updateReq.Name, v.Name)
	require.Equal(s.T(), updateReq.Details, v.Details)

	rs, err = s.Node.RemoveValidator(validatorID, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	v, err = s.Node.QueryValidator(validatorID)
	require.Error(s.T(), err)

	grantNodeReq := node.GrantNodeRequest{
		Name:        "test3",
		Certificate: cert,
		Details:     "this is a grantNode test",
	}
	rs, err = s.Node.GrantNode(grantNodeReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	noid, e := rs.Events.GetValue("grant_node", "id")
	require.NoError(s.T(), e)

	n, err := s.Node.QueryNode(noid)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), n)

	ns, err := s.Node.QueryNodes(nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(ns))

	rs, err = s.Node.RevokeNode(noid, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

}
