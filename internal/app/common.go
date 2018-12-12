package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/api"
)

var log *logrus.Logger

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
		log.Errorln(err)
		if exitOnError {
			os.Exit(1)
		}
	}
}
