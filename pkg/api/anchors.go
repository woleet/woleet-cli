package api

import (
	"errors"
	"strconv"

	"github.com/woleet/woleet-cli/pkg/models"
)

func (client *Client) GetAnchors(page int, size int, direction string, sort string) (*models.Anchors, error) {
	resp, err := client.RestyClient.
		R().
		SetQueryParams(map[string]string{
			"page":      strconv.Itoa(page),
			"size":      strconv.Itoa(size),
			"direction": direction,
			"sort":      sort,
		}).SetResult(&models.Anchors{}).
		Get(client.BaseURL + "/anchors")

	anchorsRet := resp.Result().(*models.Anchors)

	if resp.StatusCode() != 200 {
		err = errors.New(string(resp.Body()[:]))
	}
	return anchorsRet, err
}
