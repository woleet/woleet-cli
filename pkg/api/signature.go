package api

import (
	"errors"
	"log"
	"os"

	"github.com/woleet/woleet-cli/pkg/modelsBackendkit"
)

func (client *Client) GetSignature(hashToSign string, pubKey string) (*modelsBackendkit.SignatureResult, error) {

	queryMap := map[string]string{
		"hashToSign": hashToSign,
	}

	if pubKey != "" {
		queryMap["pubKey"] = pubKey
	}

	resp, err := client.RestyClient.
		R().
		SetQueryParams(queryMap).
		SetResult(&modelsBackendkit.SignatureResult{}).
		Get(client.BaseURL)

	signatureRet := resp.Result().(*modelsBackendkit.SignatureResult)

	if !(resp.StatusCode() == 200) {
		err = errors.New(string(resp.Body()[:]))
	}
	return signatureRet, err
}

func (client *Client) CheckBackendkitConnection(errLogger *log.Logger) {

	resp, _ := client.RestyClient.
		R().
		SetResult(&modelsBackendkit.SignatureResult{}).
		Get(client.BaseURL)

	if resp.StatusCode() != 400 {
		errLogger.Printf("ERROR: Unable to connect to the backendkit %v\n", string(resp.Body()[:]))
		os.Exit(1)
	}
}
