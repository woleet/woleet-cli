package api

import (
	"errors"

	"github.com/woleet/woleet-cli/pkg/models/idserver"
)

func (client *Client) GetSignature(hashToSign string, pubKey string) (*idserver.SignatureResult, error) {
	queryMap := map[string]string{
		"hashToSign": hashToSign,
	}

	if pubKey != "" {
		queryMap["pubKey"] = pubKey
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

func (client *Client) dummySignature() error {
	queryMap := map[string]string{
		"hashToSign": "0000000000000000000000000000000000000000000000000000000000000000",
	}

	resp, err := client.RestyClient.
		R().
		SetQueryParams(queryMap).
		SetResult(&idserver.SignatureResult{}).
		Get(client.BaseURL + "/sign")

	if err == nil {
		_, ok := defaultAllowedCodesMap[resp.StatusCode()]
		if !ok {
			err = errors.New(string(resp.Body()[:]))
		} else {
			if resp.Result().(*idserver.SignatureResult).SignedHash != "0000000000000000000000000000000000000000000000000000000000000000" {
				err = errors.New("Unable to sign, please check your parameters")
			}
		}
	}

	return err
}
