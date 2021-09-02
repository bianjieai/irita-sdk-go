package identity

import (
	"context"

	"github.com/bianjieai/irita-sdk-go/v2/codec"
	"github.com/bianjieai/irita-sdk-go/v2/codec/types"
	sdk "github.com/bianjieai/irita-sdk-go/v2/types"
)

type identityClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler) Client {
	return identityClient{
		BaseClient: bc,
		Marshaler:  cdc,
	}
}

func (i identityClient) Name() string {
	return ModuleName
}

func (i identityClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
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

	var pukKeyInfo *PubKeyInfo
	if request.PubkeyInfo != nil {
		if len(request.PubkeyInfo.PubKey) > 0 {
			pubkey, e := sdk.HexBytesFrom(request.PubkeyInfo.PubKey)
			if e != nil {
				return sdk.ResultTx{}, sdk.Wrap(e)
			}
			pukKeyInfo = &PubKeyInfo{
				PubKey:    pubkey.String(),
				Algorithm: request.PubkeyInfo.PubKeyAlgo,
			}
		}
	}

	credentials := ""
	if request.Credentials != nil {
		credentials = *request.Credentials
	}
	msg := &MsgCreateIdentity{
		Id:          id.String(),
		PubKey:      pukKeyInfo,
		Certificate: request.Certificate,
		Credentials: credentials,
		Owner:       sender.String(),
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
			pukKeyInfo.PubKey = pubkey.String()
			pukKeyInfo.Algorithm = request.PubkeyInfo.PubKeyAlgo
		}
	}

	credentials := DoNotModifyDesc
	if request.Credentials != nil {
		credentials = *request.Credentials
	}

	msg := &MsgUpdateIdentity{
		Id:          id.String(),
		PubKey:      &pukKeyInfo,
		Certificate: request.Certificate,
		Credentials: credentials,
		Owner:       sender.String(),
	}
	return i.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (i identityClient) QueryIdentity(id string) (QueryIdentityResp, sdk.Error) {
	conn, err := i.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return QueryIdentityResp{}, sdk.Wrap(err)
	}

	identityId, err := sdk.HexBytesFrom(id)
	if err != nil {
		return QueryIdentityResp{}, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Identity(
		context.Background(),
		&QueryIdentityRequest{Id: identityId.String()},
	)
	if err != nil {
		return QueryIdentityResp{}, sdk.Wrap(err)
	}

	return resp.Identity.Convert().(QueryIdentityResp), nil
}
