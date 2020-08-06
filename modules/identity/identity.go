package identity

import (
	"github.com/bianjieai/irita-sdk-go/codec"
	"github.com/bianjieai/irita-sdk-go/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/types"
)

type identityClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) IdentityI {
	return identityClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (i identityClient) RegisterCodec(cdc *codec.Codec) {
	registerCodec(cdc)
}

func (i identityClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateIdentity{},
		&MsgUpdateIdentity{},
	)
}

func (i identityClient) CreateIdentity(request CreateIdentityRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := i.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	id, e := sdk.HexBytesFrom(request.ID)
	if e != nil {
		return sdk.ResultTx{}, sdk.Wrap(e)
	}

	var pukKeyInfo PubKeyInfo
	if request.PubkeyInfo != nil {
		if len(request.PubkeyInfo.PubKey) > 0 {
			pubkey, e := sdk.HexBytesFrom(request.PubkeyInfo.PubKey)
			if e != nil {
				return sdk.ResultTx{}, sdk.Wrap(e)
			}
			pukKeyInfo.PubKey = pubkey
			pukKeyInfo.Algorithm = request.PubkeyInfo.PubKeyAlgo
		}
	}

	credentials := doNotModifyDesc
	if request.Credentials != nil {
		credentials = *request.Credentials
	}
	msg := MsgCreateIdentity{
		ID:          id,
		PubKey:      &pukKeyInfo,
		Certificate: request.Certificate,
		Credentials: credentials,
		Owner:       sender,
	}

	return i.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (i identityClient) UpdateIdentity(request UpdateIdentityRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := i.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	id, e := sdk.HexBytesFrom(request.ID)
	if e != nil {
		return sdk.ResultTx{}, sdk.Wrap(e)
	}

	var pukKeyInfo PubKeyInfo
	if request.PubkeyInfo != nil {
		if len(request.PubkeyInfo.PubKey) > 0 {
			pubkey, e := sdk.HexBytesFrom(request.PubkeyInfo.PubKey)
			if e != nil {
				return sdk.ResultTx{}, sdk.Wrap(e)
			}
			pukKeyInfo.PubKey = pubkey
			pukKeyInfo.Algorithm = request.PubkeyInfo.PubKeyAlgo
		}
	}

	credentials := doNotModifyDesc
	if request.Credentials != nil {
		credentials = *request.Credentials
	}

	msg := MsgUpdateIdentity{
		ID:          id,
		PubKey:      &pukKeyInfo,
		Certificate: request.Certificate,
		Credentials: credentials,
		Owner:       sender,
	}

	return i.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (i identityClient) QueryIdentity(id string) (QueryIdentityResponse, sdk.Error) {
	identityId, err := sdk.HexBytesFrom(id)
	if err != nil {
		return QueryIdentityResponse{}, sdk.Wrap(err)
	}

	param := struct{ ID sdk.HexBytes }{ID: identityId}

	var identity Identity
	if err := i.QueryWithResponse("custom/identity/identity", param, &identity); err != nil {
		return QueryIdentityResponse{}, sdk.Wrap(err)
	}

	return identity.Convert().(QueryIdentityResponse), nil
}
