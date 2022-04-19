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

	log.SetOutput(ioutil.Discard)
	var err error

	if runParameters.IsFS {
		commonInfos.mapPathFileName, err = helpers.ExploreDirectory(runParameters.Directory, runParameters.Recursive, runParameters.Filter, log)
	}

	if runParameters.IsS3 {
		commonInfos.mapPathFileName = helpers.ExploreS3(runParameters.S3Client, runParameters.S3Bucket, runParameters.Recursive, runParameters.Filter, log)
	}

	log.SetOutput(os.Stdout)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	helpers.RenameLegacyReceipts(commonInfos.mapPathFileName, runParameters.Signature, true, log)

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
