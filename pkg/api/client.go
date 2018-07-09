package api

import (
	"net/http"
	"strings"
	"time"

	"gopkg.in/resty.v1"
)

type client struct {
	BaseURL     string
	RestyClient *resty.Client
}

func GetNewClient(baseURL string, token string) *client {
	client := new(client)

	if strings.HasSuffix(baseURL, "/") {
		client.BaseURL = baseURL
	} else {
		client.BaseURL = baseURL + "/"
	}

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
