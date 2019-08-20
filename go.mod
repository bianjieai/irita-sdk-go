module github.com/irisnet/sdk-go

require (
	github.com/binance-chain/go-sdk v1.0.9 // indirect
	github.com/irisnet/irishub v0.15.0
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pkg/errors v0.8.0
	github.com/stretchr/testify v1.2.2
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/tendermint v0.31.2-rc0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
)

replace (
	github.com/tendermint/iavl => github.com/irisnet/iavl v0.12.2
	github.com/tendermint/tendermint => github.com/irisnet/tendermint v0.31.0
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
)
