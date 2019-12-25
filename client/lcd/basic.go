package lcd

import (
	"gitlab.bianjie.ai/irita/irita-sdk-go/client/basic"
)

type LiteClient interface {
	QueryAccount(address string) (AccountInfo, error)
}

type client struct {
	httpClient basic.HttpClient
}

func NewClient(c basic.HttpClient) LiteClient {
	return &client{
		httpClient: c,
	}
}
