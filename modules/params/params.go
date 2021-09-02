package params

import (
	"github.com/bianjieai/irita-sdk-go/v2/codec"
	"github.com/bianjieai/irita-sdk-go/v2/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
)

type paramsClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return paramsClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (p paramsClient) Name() string {
	return ModuleName
}

func (p paramsClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (p paramsClient) UpdateParams(requests []UpdateParamRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := p.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrapf("%s not found", baseTx.From)
	}

	var changes []ParamChange
	for _, req := range requests {
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		changes = append(changes, ParamChange{
			Subspace: req.Module,
			Key:      req.Key,
			Value:    req.Value,
		})
	}

	msg := &MsgUpdateParams{
		Changes:  changes,
		Operator: sender.String(),
	}
	return p.BuildAndSend([]sdk.Msg{msg}, baseTx)
}
