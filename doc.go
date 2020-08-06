// Package sdk is the entrance of the entire SDK function. SDKConfig is used to configure SDK parameters.
//
// The SDK mainly provides the functions of the following modules, including:
// admin, bank, keys, nft, params, record,  service, token, validator
//
// As a quick start:
//
//	options := []types.Option{
//		types.KeyDAOOption(store.NewMemory(nil)),
//	}
//	cfg, err := types.NewClientConfig(nodeURI, chainID, options...)
//	if err != nil {
//		panic(err)
//	}
//
// 	client := sdk.NewIRITAClient(cfg)
package sdk
