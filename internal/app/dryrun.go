package app

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/helpers"
)

func DryRun(runParameters *RunParameters, logInput *logrus.Logger) {
	log = logInput

	commonInfos := new(commonInfos)
	commonInfos.runParameters = *runParameters

	log.SetOutput(ioutil.Discard)
	var err error
	commonInfos.mapPathFileinfo, err = helpers.ExploreDirectory(runParameters.Directory, runParameters.Recursive, log)
	log.SetOutput(os.Stdout)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	if !runParameters.Signature {
		commonInfos.pending, commonInfos.receipt, _, _ = helpers.SeparateAll(commonInfos.mapPathFileinfo)
	} else {
		_, _, commonInfos.pending, commonInfos.receipt = helpers.SeparateAll(commonInfos.mapPathFileinfo)
	}

	commonInfos.splitPendingReceipt()

	fields := logrus.Fields{}
	fields["files"] = len(commonInfos.mapPathFileinfo)
	if commonInfos.runParameters.Prune {
		fields["pendings"] = len(commonInfos.pending)
		fields["pendings_to_delete"] = len(commonInfos.pendingToDelete)
		fields["receipts"] = len(commonInfos.receipt)
		fields["receipts_to_delete"] = len(commonInfos.receiptToDelete)
	} else {
		fields["pendings"] = len(commonInfos.pending) + len(commonInfos.pendingToDelete)
		fields["receipts"] = len(commonInfos.receipt) + len(commonInfos.receiptToDelete)
	}

	log.WithFields(fields).Infoln("Number of each category of files")
}
