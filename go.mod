module github.com/bianjieai/irita-sdk-go

go 1.14

require (
	github.com/99designs/keyring v1.1.5
	github.com/bluele/gcache v0.0.0-20190518031135-bc40bd653833
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/dvsekhvalnov/jose2go v0.0.0-20180829124132-7f401d37b68a
	github.com/gibson042/canonicaljson-go v1.0.3
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mtibben/percent v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.4
	github.com/tendermint/tm-db v0.5.1
)

replace github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.33.4-irita-200703
