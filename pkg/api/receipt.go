package api

func (client *Client) GetReceipt(anchorID string) ([]byte, error) {
	resp, err := client.RestyClient.
		R().
		Get(client.BaseURL + "/receipt/" + anchorID)

	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return resp.Body(), err
}

func (client *Client) GetReceiptToFile(anchorID string, outputPath string) error {
	resp, err := client.RestyClient.
		R().
		SetOutput(outputPath).
		Get(client.BaseURL + "/receipt/" + anchorID)

	return restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
}
