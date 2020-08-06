package utils

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/bianjieai/irita-sdk-go/types"
)

func TestSplitArray(t *testing.T) {
	data := Ints{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	subData := SubArray(4, data)
	require.Len(t, subData, 3)
}

type Ints []int

func (i Ints) Len() int {
	return len(i)
}

func (i Ints) Sub(begin, end int) sdk.SplitAble {
	return i[begin:end]
}
