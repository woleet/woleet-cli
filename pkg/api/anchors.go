package api

import (
	"errors"
	"strconv"

	"github.com/woleet/woleet-cli/pkg/models/woleetapi"
)

func (client *Client) GetAnchors(page int, size int, direction string, sort string) (*woleetapi.Anchors, error) {
	resp, err := client.RestyClient.
		R().
		SetQueryParams(map[string]string{
			"page":      strconv.Itoa(page),
			"size":      strconv.Itoa(size),
			"direction": direction,
			"sort":      sort,
		}).SetResult(&woleetapi.Anchors{}).
		Get(client.BaseURL + "/anchors")

	anchorsRet := resp.Result().(*woleetapi.Anchors)

	if resp.StatusCode() != 200 {
		err = errors.New(string(resp.Body()[:]))
	}
	return anchorsRet, err
}
