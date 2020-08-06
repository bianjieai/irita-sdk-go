package integration_test

import (
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/sm2"

	"github.com/bianjieai/irita-sdk-go/modules/identity"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/utils/uuid"
)

func (s IntegrationTestSuite) TestIdentity() {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	uuidGenerator, _ := uuid.NewV4()
	id := types.HexStringFrom(uuidGenerator.Bytes())

	testPubKeySM2 := sm2.GenPrivKey().PubKeySm2()
	testCredentials := "https://kyc.com/user/10001"
	testCertificate := `-----BEGIN CERTIFICATE-----
MIIDTDCCAjQCCQDvRoz+e/HRpDANBgkqhkiG9w0BAQsFADBoMQswCQYDVQQGEwJj
bjELMAkGA1UECAwCc2gxCzAJBgNVBAcMAnBkMQswCQYDVQQKDAJiajELMAkGA1UE
CwwCYmoxCzAJBgNVBAMMAmJqMRgwFgYJKoZIhvcNAQkBFgliakBiai5jb20wHhcN
MjAwNjEwMDkzNjMxWhcNMjAwNzEwMDkzNjMxWjBoMQswCQYDVQQGEwJjbjELMAkG
A1UECAwCc2gxCzAJBgNVBAcMAnBkMQswCQYDVQQKDAJiajELMAkGA1UECwwCYmox
CzAJBgNVBAMMAmJqMRgwFgYJKoZIhvcNAQkBFgliakBiai5jb20wggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQDUEifXes1/CXEjdH8SeSS+1x+ZlhktI8i8
9ncMeOr5oI1Mc7Kd7v85i0hrmjjZzUrHQy0Sdt2ltQjo6dtkq3wDsL4OgIqGO75z
OwG4EB0A1sJ/YTSX+fmWwy5ys19A2O5sTZOJEw3VFgiZHv1TZEiY+GVtpZ5Dti/1
t5ZzNTF+M0rpbICTxLh1GSpdhJs95yci1A8zqmPzPETVkxJwVCOg54WfpRQAiBqM
DKLjVXALuvlDDxVhB0u7kuvKAydZdV/pDs73HuY2srCOiDij3iVS01Ln02JNeMK8
IG9xRSw2eaSDp+fa1jtUXMDMmVNHCJqpQaFv0/1oN/ehUXb/DTMHAgMBAAEwDQYJ
KoZIhvcNAQELBQADggEBAKij8eUTcs+AJFPnzc3aolVZEApwvLum58WRjmoev44A
1528F4dXF7vJhIbqdOvEBy0YNQhNuNUs+JiHIFwuVvhNuAXDgXJNsvymx8fn0E5U
C90iTCiV9WhlL93S6fSelDj65sgD4Gw8Q4bBbNa/SRCu4+oBNS9BPjpcbrGllph9
7AkCGBiaabVLqGNyZJEKZpRQ3kOqdQzHYT/eHRC3hcO/KGf0vCOUTgEhHuYavMy/
JZOeFg1owNP2nZ8cD2TwDKS+T+T1rAG1ovnVp/PV7lbH1o8Kn2rwtj1S42O824Gr
2NyVhhdZkLI/uEX9mdmcFPB+oV6iiPnqEh/r2wswFgw=
-----END CERTIFICATE-----`

	pubKeyInfo := identity.PubkeyInfo{
		PubKey:     types.HexStringFrom(testPubKeySM2[:]),
		PubKeyAlgo: identity.SM2,
	}
	request := identity.CreateIdentityRequest{
		ID:          id,
		PubkeyInfo:  &pubKeyInfo,
		Certificate: testCertificate,
		Credentials: &testCredentials,
	}

	_, err := s.Identity.QueryIdentity(id)
	require.Error(s.T(), err, "not exist")

	rs, err := s.Identity.CreateIdentity(request, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	res, err := s.Identity.QueryIdentity(id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), res.Credentials, testCredentials)
	require.Contains(s.T(), res.Certificates, testCertificate)
	require.Contains(s.T(), res.PubkeyInfos, pubKeyInfo)

	test2PubKeySM2 := sm2.GenPrivKey().PubKeySm2()
	pubKeyInfo2 := identity.PubkeyInfo{
		PubKey:     types.HexStringFrom(test2PubKeySM2[:]),
		PubKeyAlgo: identity.SM2,
	}

	req2 := identity.UpdateIdentityRequest{
		ID:          id,
		PubkeyInfo:  &pubKeyInfo2,
		Certificate: testCertificate,
		Credentials: &testCredentials,
	}

	rs, err = s.Identity.UpdateIdentity(req2, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	res, err = s.Identity.QueryIdentity(id)
	require.NoError(s.T(), err)
	require.Len(s.T(), res.PubkeyInfos, 3)
}
