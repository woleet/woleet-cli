package api

import (
	"errors"
	"log"
	"os"

	"github.com/woleet/woleet-cli/pkg/models/backendkit"
)

func (client *Client) GetSignature(hashToSign string, pubKey string) (*backendkit.SignatureResult, error) {

	queryMap := map[string]string{
		"hashToSign": hashToSign,
	}

	if pubKey != "" {
		queryMap["pubKey"] = pubKey
	}

	resp, err := client.RestyClient.
		R().
		SetQueryParams(queryMap).
		SetResult(&backendkit.SignatureResult{}).
		Get(client.BaseURL)

	signatureRet := resp.Result().(*backendkit.SignatureResult)

	if !(resp.StatusCode() == 200) {
		err = errors.New(string(resp.Body()[:]))
	}
	return signatureRet, err
}

func (client *Client) CheckBackendkitConnection(errLogger *log.Logger) {

	resp, _ := client.RestyClient.
		R().
		SetResult(&backendkit.SignatureResult{}).
		Get(client.BaseURL)

	if resp.StatusCode() != 400 {
		errLogger.Printf("ERROR: Unable to connect to the backendkit %v\n", string(resp.Body()[:]))
		os.Exit(1)
	}
}
