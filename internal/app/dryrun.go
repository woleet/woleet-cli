package app

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/helpers"
)

func DryRun(runParameters *RunParameters, logInput *logrus.Logger) int {
	log = logInput

	commonInfos := initCommonInfos(runParameters)

	commonInfos.mapPathFileName = make(map[string]string)
	commonInfos.pending = make(map[string]string)
	commonInfos.pendingToDelete = make(map[string]string)
	commonInfos.receipt = make(map[string]string)
	commonInfos.receiptToDelete = make(map[string]string)

	log.SetOutput(ioutil.Discard)
	var err error
	commonInfos.mapPathFileName, err = helpers.ExploreDirectory(runParameters.Directory, runParameters.Recursive, runParameters.Filter, log)
	log.SetOutput(os.Stdout)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	if !runParameters.Signature {
		commonInfos.pending, commonInfos.receipt, _, _ = helpers.SeparateAll(commonInfos.mapPathFileName)
	} else {
		_, _, commonInfos.pending, commonInfos.receipt = helpers.SeparateAll(commonInfos.mapPathFileName)
	}

	commonInfos.splitPendingReceipt()

	fields := logrus.Fields{}
	fields["files"] = len(commonInfos.mapPathFileName)
	if commonInfos.runParameters.Prune {
		fields["pendings"] = len(commonInfos.pending)
		fields["pendingsToDelete"] = len(commonInfos.pendingToDelete)
		fields["receipts"] = len(commonInfos.receipt)
		fields["receiptsToDelete"] = len(commonInfos.receiptToDelete)
	} else {
		fields["pendings"] = len(commonInfos.pending) + len(commonInfos.pendingToDelete)
		fields["receipts"] = len(commonInfos.receipt) + len(commonInfos.receiptToDelete)
	}

	log.WithFields(fields).Infoln("Number of each category of files")
	return returnValue
}
