package api

import (
	"errors"

	"github.com/woleet/woleet-cli/pkg/models"
)

func (client *client) GetReceipt(anchorID string) (*models.Receipt, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&models.Receipt{}).
		Get(client.BaseURL + "receipt/" + anchorID)

	receiptRet := resp.Result().(*models.Receipt)

	if resp.StatusCode() != 200 {
		err = errors.New(string(resp.Body()[:]))
	}
	return receiptRet, err
}

func (client *client) GetReceiptToFile(anchorID string, path string) error {
	resp, err := client.RestyClient.
		R().
		SetOutput(path).
		Get(client.BaseURL + "receipt/" + anchorID)

	if resp.StatusCode() != 200 {
		err = errors.New(string(resp.Body()[:]))
	}
	return err
}
