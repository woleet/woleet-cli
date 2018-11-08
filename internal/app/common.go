package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/api"
)

var log *logrus.Logger

const pageSize int = 1000

type RunParameters struct {
	Signature         bool
	ExitOnError       bool
	Recursive         bool
	InvertPrivate     bool
	Strict            bool
	Prune             bool
	UnsecureSSL       bool
	Directory         string
	BaseURL           string
	Token             string
	BackendkitSignURL string
	BackendkitToken   string
	BackendkitPubKey  string
}

type commonInfos struct {
	client           *api.Client
	backendkitClient *api.Client
	mapPathFileinfo  map[string]os.FileInfo
	pending          map[string]os.FileInfo
	pendingToDelete  map[string]os.FileInfo
	receipt          map[string]os.FileInfo
	receiptToDelete  map[string]os.FileInfo
	runParameters    RunParameters
}

func errHandlerExitOnError(err error, exitOnError bool) {
	if err != nil {
		log.Errorf("%s\n", err)
		if exitOnError {
			os.Exit(1)
		}
	}
}
