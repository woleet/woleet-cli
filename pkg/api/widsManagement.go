package api

import (
	"github.com/woleet/woleet-cli/pkg/models/idserver"
)

func (client *Client) GetUserID(pubKey string) (string, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&idserver.UserDisco{}).
		Get(client.BaseURL + "/discover/user")

	userRet := resp.Result().(*idserver.UserDisco)
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return userRet.Id, err
}

func (client *Client) ListKeysFromUserID(userID string) (*[]idserver.KeyGet, error) {
	resp, err := client.RestyClient.
		R().
		SetResult([]idserver.KeyGet{}).
		Get(client.BaseURL + "/discover/keys/" + userID)

	keysRet := resp.Result().(*[]idserver.KeyGet)
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return keysRet, err
}
