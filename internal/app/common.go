package app

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/minio/minio-go/v6"
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
	FixReceipts         bool
	IDServerUnsecureSSL bool
	IsFS                bool
	IsS3                bool
	Directory           string
	S3Bucket            string
	BaseURL             string
	Token               string
	IDServerSignURL     string
	IDServerToken       string
	IDServerPubKey      string
	SignedIdentity      string
	SignedIssuerDomain  string
	Filter              *regexp.Regexp
	S3Client            *minio.Client
}

type commonInfos struct {
	client          *api.Client
	widsClient      *api.Client
	mapPathFileName map[string]string
	pending         map[string]string
	pendingToDelete map[string]string
	receipt         map[string]string
	receiptToDelete map[string]string
	runParameters   RunParameters
}

func initCommonInfos(runParameters *RunParameters) *commonInfos {
	infos := new(commonInfos)
	infos.mapPathFileName = make(map[string]string)
	infos.pending = make(map[string]string)
	infos.pendingToDelete = make(map[string]string)
	infos.receipt = make(map[string]string)
	infos.receiptToDelete = make(map[string]string)
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

func buildSignedIdentityString(user *idserver.UserDisco) string {
	signedIdentity := "CN=" + user.Identity.CommonName
	if user.Identity.Organization != "" {
		signedIdentity = signedIdentity + ",O=" + user.Identity.Organization
	}
	if user.Identity.OrganizationalUnit != "" {
		signedIdentity = signedIdentity + ",OU=" + user.Identity.OrganizationalUnit
	}
	if user.Identity.Locality != "" {
		signedIdentity = signedIdentity + ",L=" + user.Identity.Locality
	}
	if user.Identity.Country != "" {
		signedIdentity = signedIdentity + ",C=" + user.Identity.Country
	}
	if user.Email != "" {
		signedIdentity = signedIdentity + ",EMAILADDRESS=" + user.Email
	}
	return signedIdentity
}

func buildSignedIssuerDomainString(config *idserver.ConfigDisco) string {
	url, errURL := url.Parse(config.IdentityURL)
	if errURL != nil {
		return ""
	}
	domainParts := strings.Split(url.Hostname(), ".")
	if len(domainParts) == 0 {
		return ""
	} else if len(domainParts) == 1 {
		return domainParts[0]
	} else {
		return domainParts[len(domainParts)-2] + "." + domainParts[len(domainParts)-1]
	}
}
