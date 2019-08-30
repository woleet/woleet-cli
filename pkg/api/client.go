package api

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

var defaultAllowedCodesMap = map[int]struct{}{
	200: {},
}

type Client struct {
	BaseURL     string
	RestyClient *resty.Client
}

func GetNewClient(baseURL string, token string) *Client {
	client := new(Client)

	client.BaseURL = strings.TrimSuffix(baseURL, "/")

	client.RestyClient = resty.New()
	client.RestyClient.SetAuthToken(token)
	client.RestyClient.SetRetryCount(3)
	client.RestyClient.SetRetryWaitTime(1 * time.Second)
	client.RestyClient.SetRetryMaxWaitTime(3 * time.Second)
	client.RestyClient.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	if !strings.EqualFold(os.Getenv("WCLI_RESTY_DEBUG"), "true") {
		client.RestyClient.SetLogger(createRestyLogger(ioutil.Discard))
	} else {
		client.RestyClient.SetLogger(createRestyLogger(os.Stdout))

		client.RestyClient.SetDebug(true)
	}
	return client
}

func (client *Client) DisableSSLVerification() {
	client.RestyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

func restyErrHandlerAllowedCodes(resp *resty.Response, err error, allowedCodes map[int]struct{}) error {
	if err == nil {
		_, ok := allowedCodes[resp.StatusCode()]
		if !ok {
			err = errors.New(string(resp.Body()[:]))
		}
	}
	return err
}
