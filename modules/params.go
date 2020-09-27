package modules

import (
	"fmt"
	"time"

	"github.com/bianjieai/irita-sdk-go/modules/service"
	"github.com/bianjieai/irita-sdk-go/modules/token"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bianjieai/irita-sdk-go/codec"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/utils/cache"
)

type paramsQuery struct {
	sdk.Queries
	log.Logger
	cache.Cache
	cdc         codec.Marshaler
	legacyAmino *codec.LegacyAmino
	expiration  time.Duration
}

func (p paramsQuery) prefixKey(module string) string {
	return fmt.Sprintf("params:%s", module)
}

func (p paramsQuery) QueryParams(module string, res sdk.Response) sdk.Error {
	param, err := p.Get(p.prefixKey(module))
	if err == nil {
		bz := param.([]byte)
		err = p.legacyAmino.UnmarshalJSON(bz, res)
		if err != nil {
			return sdk.Wrap(err)
		}
		return nil
	}

	var path string
	switch module {
	case service.ModuleName:
		path = fmt.Sprintf("custom/%s/parameters", module)
	case token.ModuleName:
		path = fmt.Sprintf("custom/%s/params", module)
	//case "auth":
	//	path = fmt.Sprintf("custom/%s/params",module)

	default:
		return sdk.Wrapf("unsupported param query")
	}

	//path := fmt.Sprintf("custom/%s/parameters", module)
	bz, err := p.Query(path, nil)
	if err != nil {
		return sdk.Wrap(err)
	}

	err = p.legacyAmino.UnmarshalJSON(bz, res)
	if err != nil {
		return sdk.Wrap(err)
	}

	if err := p.SetWithExpire(p.prefixKey(module), bz, p.expiration); err != nil {
		p.Debug("params cache failed", "module", module)
	}
	return nil
}
