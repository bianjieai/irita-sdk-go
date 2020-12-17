package oracle

import (
	"errors"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

const (
	ModuleName = "oracle"
)

var (
	_ sdk.Msg = &MsgCreateFeed{}
	_ sdk.Msg = &MsgStartFeed{}
	_ sdk.Msg = &MsgPauseFeed{}
	_ sdk.Msg = &MsgEditFeed{}
)

func (m MsgCreateFeed) Route() string {
	return ModuleName
}

func (m MsgCreateFeed) Type() string {
	return "create_feed"
}

func (m MsgCreateFeed) ValidateBasic() error {
	if len(m.FeedName) == 0 {
		return errors.New("feedName missing")
	}

	if len(m.Providers) == 0 {
		return errors.New("providers missing")
	}

	if len(m.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(m.AggregateFunc) == 0 {
		return errors.New("aggregateFunc missing")
	}

	if len(m.ValueJsonPath) == 0 {
		return errors.New("valueJsonPath missing")
	}

	return nil
}

func (m MsgCreateFeed) GetSignBytes() []byte {
	if len(m.Providers) == 0 {
		m.Providers = nil
	}
	if len(m.ServiceFeeCap) == 0 {
		m.ServiceFeeCap = nil
	}

	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCreateFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Creator)}
}

func (m MsgStartFeed) Route() string {
	return ModuleName
}

func (m MsgStartFeed) Type() string {
	return "start_feed"
}

func (m MsgStartFeed) ValidateBasic() error {
	if len(m.FeedName) == 0 {
		return errors.New("feedName missing")
	}
	return nil
}

func (m MsgStartFeed) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgStartFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Creator)}
}

func (m MsgPauseFeed) Route() string {
	return ModuleName
}

func (m MsgPauseFeed) Type() string {
	return "pause_feed"
}

func (m MsgPauseFeed) ValidateBasic() error {
	if len(m.FeedName) == 0 {
		return errors.New("feedName missing")
	}
	return nil
}

func (m MsgPauseFeed) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgPauseFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Creator)}
}

func (m MsgEditFeed) Route() string {
	return ModuleName
}

func (m MsgEditFeed) Type() string {
	return "edit_feed"
}

func (m MsgEditFeed) ValidateBasic() error {
	if len(m.FeedName) == 0 {
		return errors.New("feedName missing")
	}
	return nil
}

func (m MsgEditFeed) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgEditFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Creator)}
}

func (f FeedContext) Convert() interface{} {
	return QueryFeedResp{
		Feed:              f.Feed,
		ServiceName:       f.ServiceName,
		Providers:         f.Providers,
		Input:             f.Input,
		Timeout:           f.Timeout,
		ServiceFeeCap:     f.ServiceFeeCap,
		RepeatedFrequency: f.RepeatedFrequency,
		ResponseThreshold: f.ResponseThreshold,
		State:             f.State,
	}
}

type Feeds []FeedContext

func (fs Feeds) Convert() interface{} {
	var frs []QueryFeedResp
	for _, f := range fs {
		frs = append(frs, QueryFeedResp{
			Feed:              f.Feed,
			ServiceName:       f.ServiceName,
			Providers:         f.Providers,
			Input:             f.Input,
			Timeout:           f.Timeout,
			ServiceFeeCap:     f.ServiceFeeCap,
			RepeatedFrequency: f.RepeatedFrequency,
			ResponseThreshold: f.ResponseThreshold,
			State:             f.State,
		})
	}
	return frs
}

type FeedValues []FeedValue

func (fvs FeedValues) Convert() interface{} {
	var fvrs []QueryFeedValueResp
	for _, fv := range fvs {
		fvrs = append(fvrs, QueryFeedValueResp{
			Data:      fv.Data,
			Timestamp: fv.Timestamp,
		})
	}
	return fvrs
}

//type validators []Validator
//
//func (vs validators) Convert() interface{} {
//	var vrs []QueryValidatorResp
//	for _, v := range vs {
//		vrs = append(vrs, QueryValidatorResp{
//			ID:          v.Id,
//			Name:        v.Name,
//			Pubkey:      v.Pubkey,
//			Certificate: v.Certificate,
//			Power:       v.Power,
//			Details:     v.Description,
//			Jailed:      v.Jailed,
//			Operator:    v.Operator,
//		})
//	}
//	return vrs
//}
//
//func (n Node) Convert() interface{} {
//	return QueryNodeResp{
//		ID:          n.Id,
//		Name:        n.Name,
//		Certificate: n.Certificate,
//	}
//}
//
//type nodes []Node
//
//func (ns nodes) Convert() interface{} {
//	var nrs []QueryNodeResp
//	for _, n := range ns {
//		nrs = append(nrs, QueryNodeResp{
//			ID:          n.Id,
//			Name:        n.Name,
//			Certificate: n.Certificate,
//		})
//	}
//	return nrs
//}
//
//func (p Params) Convert() interface{} {
//	return QueryParamsResp(p)
//}
