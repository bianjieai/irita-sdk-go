package types

import (
	"errors"
	"strings"
)

type Unit struct {
	Denom string `json:"denom"` //denom of unit
	Scale uint8  `json:"scale"` //scale of unit
}

func NewUnit(denom string, scale uint8) Unit {
	return Unit{
		Denom: denom,
		Scale: scale,
	}
}

//GetScaleFactor return 1 * 10^scale
func (u Unit) GetScaleFactor() Int {
	return NewIntWithDecimal(1, int(u.Scale))
}

type CoinType struct {
	Name     string `json:"name"`      //description name of CoinType
	MinUnit  Unit   `json:"min_unit"`  //the min unit of CoinType
	MainUnit Unit   `json:"main_unit"` //the max unit of CoinType
	Desc     string `json:"desc"`      //the description of CoinType
}

//ToMainCoin return the main denom coin from args
func (ct CoinType) ConvertToMainCoin(coin Coin) (DecCoin, error) {
	if !ct.hasUnit(coin.Denom) {
		return DecCoin{}, errors.New("coinType unit (%s) not defined" + coin.Denom)
	}

	if ct.isMainUnit(coin.Denom) {
		return DecCoin{}, nil
	}

	// dest amount = src amount * (10^(dest scale) / 10^(src scale))
	dstScale := NewDecFromInt(ct.MainUnit.GetScaleFactor())
	srcScale := NewDecFromInt(ct.MinUnit.GetScaleFactor())
	amount := NewDecFromInt(coin.Amount)

	amt := amount.Mul(dstScale).Quo(srcScale)
	return NewDecCoinFromDec(ct.MainUnit.Denom, amt), nil
}

//ToMinCoin return the min denom coin from args
func (ct CoinType) ConvertToMinCoin(coin DecCoin) (newCoin Coin, err error) {
	if !ct.hasUnit(coin.Denom) {
		return newCoin, errors.New("coinType unit (%s) not defined" + coin.Denom)
	}

	if ct.isMinUnit(coin.Denom) {
		newCoin, _ := coin.TruncateDecimal()
		return newCoin, nil
	}

	// dest amount = src amount * (10^(dest scale) / 10^(src scale))
	srcScale := NewDecFromInt(ct.MainUnit.GetScaleFactor())
	dstScale := NewDecFromInt(ct.MinUnit.GetScaleFactor())
	amount := coin.Amount

	amt := amount.Mul(dstScale).Quo(srcScale)
	return NewCoin(ct.MinUnit.Denom, amt.RoundInt()), nil
}

func (ct CoinType) isMainUnit(name string) bool {
	return ct.MainUnit.Denom == strings.TrimSpace(name)
}

func (ct CoinType) isMinUnit(name string) bool {
	return ct.MinUnit.Denom == strings.TrimSpace(name)
}

func (ct CoinType) hasUnit(name string) bool {
	if ct.isMainUnit(name) || ct.isMinUnit(name) {
		return true
	}
	return false
}
