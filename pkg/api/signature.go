package api

import (
	"github.com/woleet/woleet-cli/pkg/models/idserver"
)

func (client *Client) GetSignature(hashToSign string, pubKey string, integratedSignature bool) (*idserver.SignatureResult, error) {
	queryMap := map[string]string{
		"hashToSign": hashToSign,
	}

	if pubKey != "" {
		queryMap["pubKey"] = pubKey
	}

	if integratedSignature {
		queryMap["identityToSign"] = ""
	}

	resp, err := client.RestyClient.
		R().
		SetQueryParams(queryMap).
		SetResult(&idserver.SignatureResult{}).
		Get(client.BaseURL + "/sign")

	signatureRet := resp.Result().(*idserver.SignatureResult)
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return signatureRet, err
}
