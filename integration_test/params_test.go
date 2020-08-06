package integration_test

import (
	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/modules/params"
	"github.com/bianjieai/irita-sdk-go/types"
)

func (s IntegrationTestSuite) TestParams() {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	var request = []params.UpdateParamRequest{
		{
			Module: "service",
			Key:    "MaxRequestTimeout",
			Value:  "200",
		},
	}
	rs, err := s.Params.UpdateParams(request, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), rs.Hash)

	p, err := s.Service.QueryParams()
	require.NoError(s.T(), err)
	require.Equal(s.T(), int64(200), p.MaxRequestTimeout)
}
