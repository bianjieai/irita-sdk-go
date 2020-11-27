package integration_test

import (
	"github.com/stretchr/testify/require"

	sdk "github.com/bianjieai/irita-sdk-go/types"

	"github.com/bianjieai/irita-sdk-go/modules/validator"
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
MIIBkzCCATkCFGkwIuNrP0KGhcypL+10kY9aA2dUMAoGCCqBHM9VAYN1MEUxCzAJ
BgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5l
dCBXaWRnaXRzIFB0eSBMdGQwHhcNMjAxMTI2MTA1NDE5WhcNMjAxMjI2MTA1NDE5
WjBTMQswCQYDVQQGEwJDTjELMAkGA1UECAwCU0gxCzAJBgNVBAcMAlNIMQswCQYD
VQQKDAJCSjELMAkGA1UECwwCQkoxEDAOBgNVBAMMB2NzY2hhaW4wWTATBgcqhkjO
PQIBBggqgRzPVQGCLQNCAAQBYQNhduxB79KDhBMAV4SaNu8Tc2wHEDVLnOMy2lta
j4e0guc/xEhoV+hKYBAkYFjsbBm6Oi2Yx+bkPV96kugUMAoGCCqBHM9VAYN1A0gA
MEUCIBav+UKwrL8ChHCF7AwSbuKmF2y+qDFzHSTjK2QC2k8QAiEA5CP3hvMc6qvX
LTWCInii/I8Skv+Nuk034CK3u1fThnk=
-----END CERTIFICATE-----
`

	createReq := validator.CreateValidatorRequest{
		Name:        "test1",
		Certificate: cert,
		Power:       10,
		Details:     "this is a test",
	}

	validatorClient := validator.NewClient(s.IRITAClient.BaseClient, s.IRITAClient.AppCodec())
	s.IRITAClient.RegisterModule(validatorClient)

	rs, err := validatorClient.CreateValidator(createReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	validatorID, er := rs.Events.GetValue("create_validator", "validator")
	require.NoError(s.T(), er)

	v, err := validatorClient.QueryValidator(validatorID)
	require.NoError(s.T(), err)

	vs, err := validatorClient.QueryValidators(nil, 0, 0, false)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), vs)

	updateReq := validator.UpdateValidatorRequest{
		ID:          validatorID,
		Name:        "test2",
		Certificate: cert,
		Power:       10,
		Details:     "this is a updated test",
	}
	rs, err = validatorClient.UpdateValidator(updateReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	v, err = validatorClient.QueryValidator(validatorID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), updateReq.Name, v.Name)
	require.Equal(s.T(), updateReq.Details, v.Details)

	rs, err = validatorClient.RemoveValidator(validatorID, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	v, err = validatorClient.QueryValidator(validatorID)
	require.Error(s.T(), err)
}
