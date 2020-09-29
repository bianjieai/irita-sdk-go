package modules

import (
	"context"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/modules/token"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/utils/cache"
)

type tokenQuery struct {
	q sdk.Queries
	sdk.GRPCClient
	cdc codec.Marshaler
	log.Logger
	cache.Cache
}

func (l tokenQuery) QueryToken(denom string) (sdk.Token, error) {
	denom = strings.ToLower(denom)
	if t, err := l.Get(l.prefixKey(denom)); err == nil {
		return t.(sdk.Token), nil
	}

	conn, err := l.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return sdk.Token{}, sdk.Wrap(err)
	}

	response, err := token.NewQueryClient(conn).Token(
		context.Background(),
		&token.QueryTokenRequest{Denom: denom},
	)
	if err != nil {
		return sdk.Token{}, sdk.Wrap(err)
	}

	var srcToken token.TokenInterface
	if err = l.cdc.UnpackAny(response.Token, &srcToken); err != nil {
		return sdk.Token{}, sdk.Wrap(err)
	}
	token := srcToken.(*token.Token).Convert().(sdk.Token)
	l.SaveTokens(token)
	return token, nil
}

func (l tokenQuery) SaveTokens(tokens ...sdk.Token) {
	for _, t := range tokens {
		err1 := l.Set(l.prefixKey(t.Symbol), t)
		err2 := l.Set(l.prefixKey(t.MinUnit), t)
		if err1 != nil || err2 != nil {
			l.Debug("cache token failed", "symbol", t.Symbol)
		}
	}
}

func (l tokenQuery) ToMinCoin(coins ...sdk.DecCoin) (dstCoins sdk.Coins, err sdk.Error) {
	for _, coin := range coins {
		token, err := l.QueryToken(coin.Denom)
		if err != nil {
			return nil, sdk.Wrap(err)
		}

		minCoin, err := token.GetCoinType().ConvertToMinCoin(coin)
		if err != nil {
			return nil, sdk.Wrap(err)
		}
		dstCoins = append(dstCoins, minCoin)
	}
	return dstCoins.Sort(), nil
}

func (l tokenQuery) ToMainCoin(coins ...sdk.Coin) (dstCoins sdk.DecCoins, err sdk.Error) {
	for _, coin := range coins {
		token, err := l.QueryToken(coin.Denom)
		if err != nil {
			return dstCoins, sdk.Wrap(err)
		}

		mainCoin, err := token.GetCoinType().ConvertToMainCoin(coin)
		if err != nil {
			return dstCoins, sdk.Wrap(err)
		}
		dstCoins = append(dstCoins, mainCoin)
	}
	return dstCoins.Sort(), nil
}

func (l tokenQuery) prefixKey(symbol string) string {
	return fmt.Sprintf("token:%s", symbol)
}
