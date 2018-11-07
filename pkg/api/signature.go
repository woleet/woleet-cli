package api

import (
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
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return signatureRet, err
}

func (client *Client) CheckBackendkitConnection() error {

	queryMap := map[string]string{
		"hashToSign": "0000000000000000000000000000000000000000000000000000000000000000",
	}

	resp, err := client.RestyClient.
		R().
		SetQueryParams(queryMap).
		SetResult(&backendkit.SignatureResult{}).
		Get(client.BaseURL)

	return restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
}
