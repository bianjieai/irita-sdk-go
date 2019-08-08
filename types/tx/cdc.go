package tx

import (
	"github.com/irisnet/irishub/app/v1"
	"github.com/tendermint/go-amino"
)

var Cdc *amino.Codec

func init() {
	Cdc = v1.MakeCodec()
}
