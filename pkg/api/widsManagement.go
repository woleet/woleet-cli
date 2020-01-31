package api

import (
	"github.com/woleet/woleet-cli/pkg/models/idserver"
)

func (client *Client) GetServerConfig() (*idserver.ConfigDisco, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&idserver.ConfigDisco{}).
		Get(client.BaseURL + "/discover/config")

	discoRet := resp.Result().(*idserver.ConfigDisco)
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)
	return discoRet, err
}

func (client *Client) GetUser() (*idserver.UserDisco, error) {
	resp, _ := client.RestyClient.
		R().
		SetResult(&idserver.UserDisco{}).
		Get(client.BaseURL + "/discover/user")

	userRet := resp.Result().(*idserver.UserDisco)
	allowedCodesMap := map[int]struct{}{
		200: {},
		204: {},
		404: {},
	}
	err := restyErrHandlerAllowedCodes(resp, nil, allowedCodesMap)

	if resp.StatusCode() == 204 || resp.StatusCode() == 404 {
		respConfig, errConfig := client.RestyClient.
			R().
			SetResult(&idserver.ConfigDisco{}).
			Get(client.BaseURL + "/discover/config")
		errConfig = restyErrHandlerAllowedCodes(respConfig, errConfig, defaultAllowedCodesMap)
		return &idserver.UserDisco{Id: "admin"}, errConfig
	}
	return userRet, err
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

func (client *Client) GetUserDiscoFromPubkey(pubKey string) (*idserver.UserDisco, error) {
	resp, err := client.RestyClient.
		R().
		SetResult(&idserver.UserDisco{}).
		Get(client.BaseURL + "/discover/user/" + pubKey)

	userRet := resp.Result().(*idserver.UserDisco)
	err = restyErrHandlerAllowedCodes(resp, err, defaultAllowedCodesMap)

	return userRet, err
}
