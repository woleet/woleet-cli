package api

import (
	"github.com/woleet/woleet-cli/pkg/models/woleetapi"
)

func (client *Client) PostAnchor(anchor *woleetapi.Anchor) (*woleetapi.Anchor, error) {
	resp, err := client.RestyClient.
		R().
		SetBody(anchor).
		SetResult(&woleetapi.Anchor{}).
		Post(client.BaseURL + "/anchor")

	anchorRet := resp.Result().(*woleetapi.Anchor)
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return anchorRet, err
}

func (client *Client) GetAnchor(anchorID string) (*woleetapi.Anchor, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&woleetapi.Anchor{}).
		Get(client.BaseURL + "/anchor/" + anchorID)

	anchorRet := resp.Result().(*woleetapi.Anchor)
	allowedCodesMap := map[int]struct{}{
		200: {},
		202: {},
	}
	err = restyErrHandlerAllowedCodes(resp, err, allowedCodesMap)
	return anchorRet, err
}
