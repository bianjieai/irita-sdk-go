package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bianjieai/irita-sdk-go/utils/bech32"
)

func TestGetFromBech32(t *testing.T) {
	addBz, err := bech32.GetFromBech32("caa1rgnu8grzt6mwnjg7jss7w0sfyjn67g4em9njf5", "caa")
	require.NoError(t, err)
	addr, err := bech32.ConvertAndEncode("iaa", addBz)
	require.NoError(t, err)
	fmt.Println(addr)

	addBz, err = bech32.GetFromBech32("ccp1ulx45dfpqdnyust0a8ezcdhjjn3cekla9wx9lnpqmer9httzvstx64mns5njqzr4fdd", "ccp")
	require.NoError(t, err)

	addr, err = bech32.ConvertAndEncode("icp", addBz)
	require.NoError(t, err)
	fmt.Println(addr)
}
