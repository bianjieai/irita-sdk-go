package params

import (
	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type paramsClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) ParamsI {
	return paramsClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (p paramsClient) RegisterCodec(cdc *codec.Codec) {
	registerCodec(cdc)
}

func (p paramsClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
}

func (p paramsClient) UpdateParams(requests []UpdateParamRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := p.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s not found", baseTx.From)
	}

	var changes []ParamChange
	for _, req := range requests {
		v, err := p.MarshalJSON(req.Value)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		changes = append(changes, ParamChange{
			Subspace: req.Module,
			Key:      req.Key,
			Value:    string(v),
		})
	}

	msg := MsgUpdateParams{
		Changes:  changes,
		Operator: sender,
	}

	return p.BuildAndSend([]sdk.Msg{msg}, baseTx)
}
