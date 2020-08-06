package integration_test

import (
	"fmt"

	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/modules/record"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

func (s IntegrationTestSuite) TestRecord() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	num := 5
	contents := make([]record.Content, num)
	for i := 0; i < num; i++ {
		contents[i] = record.Content{
			Digest:     s.RandStringOfLength(10),
			DigestAlgo: s.RandStringOfLength(5),
			URI:        fmt.Sprintf("https://%s", s.RandStringOfLength(10)),
			Meta:       s.RandStringOfLength(20),
		}
	}

	req := record.CreateRecordRequest{
		Contents: contents,
	}

	recordID, err := s.Record.CreateRecord(req, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), recordID)

	record, err := s.Record.QueryRecord(recordID)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), record.Contents)

	for i := 0; i < num; i++ {
		require.EqualValues(s.T(), contents[i], record.Contents[i])
	}
}
