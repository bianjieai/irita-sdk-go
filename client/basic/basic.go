package basic

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

const (
	timeout = 10 * time.Second
)

var (
	retryStatus = []int{http.StatusBadRequest, http.StatusInternalServerError}
)

type HttpClient interface {
	Get(path string, params map[string]string) (*http.Response, []byte, error)
	Post(path string, params map[string]string, body interface{}) (*http.Response, []byte, error)
}

type client struct {
	baseUrl string
	apiUrl  string
}

func NewClient(baseUrl string) HttpClient {
	return &client{baseUrl: baseUrl, apiUrl: fmt.Sprintf("%s", baseUrl)}
}

func (c *client) Get(path string, params map[string]string) (*http.Response, []byte, error) {
	var (
		err error
	)
	request := gorequest.New().Timeout(timeout).Get(c.apiUrl + path)

	if params != nil {
		for k, v := range params {
			request = request.Param(k, v)
		}
	}
	response, body, errs := request.Retry(3, 5*time.Second, retryStatus...).EndBytes()
	if len(errs) > 0 || response.StatusCode != http.StatusOK {
		if len(errs) > 0 {
			err = fmt.Errorf("http get fail, err is %s", errs[0].Error())
		} else {
			err = fmt.Errorf("bad request, status code is %d, response body is %s", response.StatusCode, string(body))
		}
	}

	return response, body, err
}

func (c *client) Post(path string, params map[string]string, body interface{}) (*http.Response, []byte, error) {
	var (
		err error
	)

	request := gorequest.New().Timeout(timeout * 2).Post(c.apiUrl + path)
	if params != nil {
		for k, v := range params {
			request = request.Param(k, v)
		}
	}

	response, resBody, errs := request.Send(body).EndBytes()
	if len(errs) > 0 || response.StatusCode != http.StatusOK {
		if len(errs) > 0 {
			err = fmt.Errorf("http get fail, err is %s", errs[0].Error())
		} else {
			err = fmt.Errorf("bad request, status code is %d, response body is %s", response.StatusCode, string(resBody))
		}
	}

	return response, resBody, err
}
