package client

import (
	"gitlab.bianjie.ai/irita/irita-sdk-go/keys"
	"gitlab.bianjie.ai/irita/irita-sdk-go/types"
	"gitlab.bianjie.ai/irita/irita-sdk-go/util"
	"testing"
)

var (
	baseUrl     = "http://v2.irisnet-lcd.dev.rainbow.one"
	nodeUrl     = "tcp://35.236.146.181:30657"
	networkType = types.Testnet
	km          keys.KeyManager
)

func TestMain(m *testing.M) {
	if k, err := keys.NewKeyStoreKeyManager("../keys/ks_1234567890.json", "1234567890"); err != nil {
		panic(err)
	} else {
		km = k
	}
	m.Run()
}

func TestNewIRISnetClient(t *testing.T) {
	c, err := NewIRITAClient(baseUrl, nodeUrl, networkType, km)
	if err != nil {
		t.Fatal(err)
	} else {
		if res, err := c.QueryAccount("faa1282eufkw9qgm55symgqqg38nremslvggpylkht"); err != nil {
			t.Fatal(err)
		} else {
			t.Log(util.ToJsonIgnoreErr(res))
		}
	}
}
