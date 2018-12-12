package api

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	gologger "log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/resty.v1"
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
		func(r *resty.Response) (bool, error) {
			return r.StatusCode() == http.StatusTooManyRequests, nil
		},
	)
	if !strings.EqualFold(os.Getenv("WCLI_RESTY_DEBUG"), "true") {
		client.RestyClient.Log = gologger.New(ioutil.Discard, "RESTY - ", gologger.LstdFlags)
	} else {
		client.RestyClient.Log = gologger.New(os.Stdout, "RESTY - ", gologger.LstdFlags)
		client.RestyClient.SetDebug(true)
	}
	return client
}

func (client *Client) SetDomain(domain string) {
	if !strings.EqualFold(domain, "") {
		client.RestyClient.SetHeader("Domain", domain)
	}
}

func (client *Client) SetCustomLogger(customLogger *log.Logger) {
	client.RestyClient.Log = customLogger
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
