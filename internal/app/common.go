package app

import (
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/api"
	"github.com/woleet/woleet-cli/pkg/models/idserver"
)

var log *logrus.Logger
var returnValue = 0

const pageSize int = 1000

type RunParameters struct {
	Signature           bool
	ExitOnError         bool
	Recursive           bool
	InvertPrivate       bool
	Strict              bool
	Prune               bool
	IDServerUnsecureSSL bool
	Directory           string
	BaseURL             string
	Token               string
	IDServerSignURL     string
	IDServerToken       string
	IDServerPubKey      string
	Include             *regexp.Regexp
}

type commonInfos struct {
	client          *api.Client
	widsClient      *api.Client
	mapPathFileinfo map[string]os.FileInfo
	pending         map[string]os.FileInfo
	pendingToDelete map[string]os.FileInfo
	receipt         map[string]os.FileInfo
	receiptToDelete map[string]os.FileInfo
	runParameters   RunParameters
}

func initCommonInfos(runParameters *RunParameters) *commonInfos {
	infos := new(commonInfos)
	infos.mapPathFileinfo = make(map[string]os.FileInfo)
	infos.pending = make(map[string]os.FileInfo)
	infos.pendingToDelete = make(map[string]os.FileInfo)
	infos.receipt = make(map[string]os.FileInfo)
	infos.receiptToDelete = make(map[string]os.FileInfo)
	infos.runParameters = *runParameters
	return infos
}

func errHandlerExitOnError(err error, exitOnError bool) {
	if err != nil {
		returnValue = 1
		log.Errorln(err)
		if exitOnError {
			os.Exit(1)
		}
	}
}

func checkWIDSConnectionPubKey(commonInfos *commonInfos) {
	user, errUser := commonInfos.widsClient.GetUser()
	if errUser != nil {
		log.Fatalf("Unable to request current userID on Woleet.ID Server: %s\n", errUser)
	}

	if strings.EqualFold(user.Id, "admin") {
		if strings.EqualFold(commonInfos.runParameters.IDServerPubKey, "") {
			return
		}
		user, errUser = commonInfos.widsClient.GetUserDiscoFromPubkey(commonInfos.runParameters.IDServerPubKey)
		if errUser != nil {
			log.Fatalf("This public key does not exists on this Woleet.ID Server: %s\n", errUser)
		}
		if user.Mode == idserver.USERMODEENUM_ESIGN {
			log.Fatalln("You can't sign with a user configured for e-signature with an admin token")
		}
	}

	pubKeys, errPubKeys := commonInfos.widsClient.ListKeysFromUserID(user.Id)

	if errPubKeys != nil {
		log.Fatalf("Unable to get current userID public keys on this Woleet.ID Server: %s\n", errPubKeys)
	}

	for _, pubKey := range *pubKeys {
		if strings.EqualFold(pubKey.PubKey, commonInfos.runParameters.IDServerPubKey) || strings.EqualFold(pubKey.Id, user.DefaultKeyId) {
			if pubKey.Status != idserver.KEYSTATUSENUM_ACTIVE {
				log.Fatalf("The specified pulblic key is not active")
			}
			if pubKey.Device != idserver.KEYDEVICEENUM_SERVER {
				log.Fatalf("The specified public key is not owned by the server")
			}
			return
		}
	}
	log.Fatalf("Unable to find specified publicKey on this Woleet.ID Server with provided token")
}
