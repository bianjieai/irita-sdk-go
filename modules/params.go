package modules

import (
	"fmt"
	"time"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/modules/service"
	"github.com/bianjieai/irita-sdk-go/modules/token"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/utils/cache"
	"github.com/bianjieai/irita-sdk-go/utils/log"
)

type paramsQuery struct {
	sdk.Queries
	*log.Logger
	cache.Cache
	cdc        codec.Marshaler
	expiration time.Duration
}

func (p paramsQuery) prefixKey(module string) string {
	return fmt.Sprintf("params:%s", module)
}

func (p paramsQuery) QueryParams(module string, res sdk.Response) sdk.Error {
	param, err := p.Get(p.prefixKey(module))
	if err == nil {
		bz := param.([]byte)
		if err = p.cdc.UnmarshalJSON(bz, res); err != nil {
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

	if err = p.cdc.UnmarshalJSON(bz, res); err != nil {
		return sdk.Wrap(err)
	}

	if err := p.SetWithExpire(p.prefixKey(module), bz, p.expiration); err != nil {
		p.Warn().
			Str("module", module).
			Msg("params cache failed")
	}

	return nil
}