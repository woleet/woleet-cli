package api

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/resty.v1"
)

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
	return client
}

func (client *Client) SetCustomLogger(customLogger *log.Logger) {
	client.RestyClient.Log = customLogger
}

func (client *Client) DisableSSLVerification() {
	client.RestyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}
