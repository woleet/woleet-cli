package api

import (
	"errors"

	"github.com/woleet/woleet-cli/pkg/models"
)

func (client *client) PostAnchor(anchor *models.Anchor) (*models.Anchor, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&models.Anchor{}).
		SetBody(anchor).
		Post(client.BaseURL + "anchor")

	anchorRet := resp.Result().(*models.Anchor)

	if resp.StatusCode() != 200 {
		err = errors.New(string(resp.Body()[:]))
	}
	return anchorRet, err
}

func (client *client) GetAnchor(anchorID string) (*models.Anchor, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&models.Anchor{}).
		Get(client.BaseURL + "anchor/" + anchorID)

	anchorRet := resp.Result().(*models.Anchor)

	if !(resp.StatusCode() == 200 || resp.StatusCode() == 202) {
		err = errors.New(string(resp.Body()[:]))
	}
	return anchorRet, err
}
