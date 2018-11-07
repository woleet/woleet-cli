package api

import (
	"github.com/woleet/woleet-cli/pkg/models/woleetapi"
)

func (client *Client) GetReceipt(anchorID string) (*woleetapi.Receipt, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&woleetapi.Receipt{}).
		Get(client.BaseURL + "/receipt/" + anchorID)

	receiptRet := resp.Result().(*woleetapi.Receipt)
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return receiptRet, err
}

func (client *Client) GetReceiptToFile(anchorID string, outputPath string) error {
	resp, err := client.RestyClient.
		R().
		SetOutput(outputPath).
		Get(client.BaseURL + "/receipt/" + anchorID)

	return restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
}
