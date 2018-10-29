package api

import (
	"errors"

	"github.com/woleet/woleet-cli/pkg/models/woleetapi"
)

func (client *Client) PostAnchor(anchor *woleetapi.Anchor) (*woleetapi.Anchor, error) {
	resp, err := client.RestyClient.
		R().
		SetBody(anchor).
		SetResult(&woleetapi.Anchor{}).
		Post(client.BaseURL + "/anchor")

	anchorRet := resp.Result().(*woleetapi.Anchor)

	if resp.StatusCode() != 200 {
		err = errors.New(string(resp.Body()[:]))
	}
	return anchorRet, err
}

func (client *Client) GetAnchor(anchorID string) (*woleetapi.Anchor, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&woleetapi.Anchor{}).
		Get(client.BaseURL + "/anchor/" + anchorID)

	anchorRet := resp.Result().(*woleetapi.Anchor)

	if !(resp.StatusCode() == 200 || resp.StatusCode() == 202) {
		err = errors.New(string(resp.Body()[:]))
	}
	return anchorRet, err
}
