package modules

import (
	"fmt"
	"time"

	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/modules/bank"
	sdk "github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/utils/cache"
	"github.com/bianjieai/irita-sdk-go/utils/log"
)

// Must be used with locker, otherwise there are thread safety issues
type accountQuery struct {
	sdk.Queries
	*log.Logger
	cache.Cache
	cdc        codec.Marshaler
	km         sdk.KeyManager
	expiration time.Duration
}

func (a accountQuery) QueryAndRefreshAccount(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := a.Get(a.prefixKey(address))
	if err != nil {
		return a.refresh(address)
	}

	acc := account.(accountInfo)
	baseAcc := sdk.BaseAccount{
		Address:       sdk.MustAccAddressFromBech32(address),
		AccountNumber: acc.N,
		Sequence:      acc.S + 1,
	}
	a.saveAccount(baseAcc)

	a.Debug().
		Str("address", address).
		Msg("query account from cache")

	return baseAcc, nil
}

func (a accountQuery) QueryAccount(address string) (sdk.BaseAccount, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	param := struct {
		Address sdk.AccAddress `json:"account"`
	}{
		Address: addr,
	}

	bz, er := a.Query("custom/auth/account", param)
	if er != nil {
		return sdk.BaseAccount{}, sdk.Wrap(er)
	}
	var account bank.BaseAccount
	if err := a.cdc.UnmarshalJSON(bz, &account); err != nil {
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	a.Debug().
		Str("address", address).
		Msg("query account from chain")
	return account.Convert().(sdk.BaseAccount), nil
}

func (a accountQuery) QueryAddress(name, password string) (sdk.AccAddress, sdk.Error) {
	addr, err := a.Get(a.prefixKey(name))
	if err == nil {
		address, err := sdk.AccAddressFromBech32(addr.(string))
		if err != nil {
			a.Warn().
				Str("name", name).
				Msg("invalid address")
			_ = a.Remove(a.prefixKey(name))
		} else {
			return address, nil
		}
	}

	address, err := a.km.Find(name, password)
	if err != nil {
		a.Warn().
			Str("name", name).
			Msg("can't find account")
		return address, sdk.Wrap(err)
	}

	if err := a.SetWithExpire(a.prefixKey(name), address.String(), a.expiration); err != nil {
		a.Warn().
			Str("name", name).
			Msg("cache user failed")
	}
	a.Debug().
		Str("name", name).
		Str("address", address.String()).
		Msg("query user from cache")

	return address, nil
}

func (a accountQuery) removeCache(address string) bool {
	return a.Remove(a.prefixKey(address))
}

func (a accountQuery) refresh(address string) (sdk.BaseAccount, sdk.Error) {
	account, err := a.QueryAccount(address)
	if err != nil {
		a.Err(err).
			Str("address", address).
			Msg("update cache failed")
		return sdk.BaseAccount{}, sdk.Wrap(err)
	}

	a.saveAccount(account)
	return account, nil
}

func (a accountQuery) saveAccount(account sdk.BaseAccount) {
	address := account.Address.String()
	info := accountInfo{
		N: account.AccountNumber,
		S: account.Sequence,
	}
	if err := a.SetWithExpire(a.prefixKey(address), info, a.expiration); err != nil {
		a.Warn().
			Str("address", address).
			Msg("cache account failed")
		return
	}
	a.Debug().
		Str("address", address).
		Msgf("cache account %s", a.expiration.String())
}

func (a accountQuery) prefixKey(address string) string {
	return fmt.Sprintf("account:%s", address)
}

type accountInfo struct {
	N uint64 `json:"n"`
	S uint64 `json:"s"`
}
