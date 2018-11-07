package api

import (
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
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return anchorsRet, err
}
