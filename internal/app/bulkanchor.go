package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/api"
	"github.com/woleet/woleet-cli/pkg/helpers"
	"github.com/woleet/woleet-cli/pkg/models/woleetapi"
)

func BulkAnchor(runParameters *RunParameters, logInput *logrus.Logger) {
	commonInfos := new(commonInfos)
	commonInfos.runParameters = *runParameters

	log = logInput

	commonInfos.client = api.GetNewClient(runParameters.BaseURL, runParameters.Token)

	var err error
	commonInfos.mapPathFileinfo, err = helpers.ExploreDirectory(runParameters.Directory, runParameters.Recursive, log)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	if !commonInfos.runParameters.Signature {
		commonInfos.pending, commonInfos.receipt, _, _ = helpers.SeparateAll(commonInfos.mapPathFileinfo)
	} else {
		_, _, commonInfos.pending, commonInfos.receipt = helpers.SeparateAll(commonInfos.mapPathFileinfo)
	}

	if runParameters.Signature {
		//Check backendkit connection
		commonInfos.backendkitClient = api.GetNewClient(commonInfos.runParameters.BackendkitSignURL, commonInfos.runParameters.BackendkitToken)
		if commonInfos.runParameters.UnsecureSSL {
			commonInfos.backendkitClient.DisableSSLVerification()
		}
		errBackendkit := commonInfos.backendkitClient.CheckBackendkitConnection()
		if errBackendkit != nil {
			log.Fatalf("Unable to connect to the backendkit: %s\n", errBackendkit)
		}
	}

	commonInfos.splitPendingReceipt()
	commonInfos.getReceipts(commonInfos.pending)
	if !commonInfos.runParameters.Prune {
		commonInfos.getReceipts(commonInfos.pendingToDelete)
	} else {
		for path := range commonInfos.pendingToDelete {
			log.Infof("Deleting old pending file %s\n", path)
			os.Remove(path)
		}
		for path := range commonInfos.receiptToDelete {
			log.Infof("Deleting old receipt file %s\n", path)
			os.Remove(path)
		}
	}
	commonInfos.checkStandardFiles()
}

func (commonInfos *commonInfos) splitPendingReceipt() {
	for path, fileinfo := range commonInfos.pending {
		errHandlerExitOnError(commonInfos.sortFile(path, fileinfo, true, false), commonInfos.runParameters.ExitOnError)
	}
	for path, fileinfo := range commonInfos.receipt {
		errHandlerExitOnError(commonInfos.sortFile(path, fileinfo, false, true), commonInfos.runParameters.ExitOnError)
	}
}

func (commonInfos *commonInfos) sortFile(path string, fileinfo os.FileInfo, pending bool, receipt bool) error {
	anchorNameInfo, erranchorNameInfo := helpers.GetAnchorIDFromName(fileinfo)
	if erranchorNameInfo != nil {
		return erranchorNameInfo
	}

	// Extracting the file's original path by the name of the pending/receipt
	originalFilePath := strings.TrimSuffix(path, "-"+anchorNameInfo.AnchorID+anchorNameInfo.Suffix)

	// If there is no strict mode, there is nothing to check and the file will not be reanchored
	if !commonInfos.runParameters.Strict {
		delete(commonInfos.mapPathFileinfo, originalFilePath)
		return nil
	}

	// If strict mode is actived, we check that the hash of the file
	// is the same as the one in the pending/receipt
	// If the file does not exists anymore and the prune mode is set the file will be deleted
	// if the prune mode is not set the file will be converted to a proper receipt
	_, exists := commonInfos.mapPathFileinfo[originalFilePath]
	if !exists {
		if commonInfos.runParameters.Prune {
			if pending {
				commonInfos.pendingToDelete[path] = fileinfo
				delete(commonInfos.pending, path)
			}
			if receipt {
				commonInfos.receiptToDelete[path] = fileinfo
				delete(commonInfos.receipt, path)
			}
		}
	}

	receiptJSON, errReceiptJSON := ioutil.ReadFile(path)
	if errReceiptJSON != nil {
		return errReceiptJSON
	}

	var receiptUnmarshalled woleetapi.Receipt
	json.Unmarshal(receiptJSON, &receiptUnmarshalled)
	hash, errhash := helpers.HashFile(originalFilePath)
	if errhash != nil {
		return errhash
	}

	// If the hashes are equal, there is nothing to do
	if (!commonInfos.runParameters.Signature && strings.EqualFold(hash, receiptUnmarshalled.TargetHash)) || (commonInfos.runParameters.Signature && strings.EqualFold(hash, receiptUnmarshalled.Signature.SignedHash)) {
		// File already anchored and valid
		delete(commonInfos.mapPathFileinfo, originalFilePath)
		return nil
	}
	// If they are not and there is a prune flag, the old pending file will be marked for deletion
	if commonInfos.runParameters.Prune {
		if pending {
			commonInfos.pendingToDelete[path] = fileinfo
			delete(commonInfos.pending, path)
		}
		if receipt {
			commonInfos.receiptToDelete[path] = fileinfo
			delete(commonInfos.receipt, path)
		}
	}
	return nil
}

func (commonInfos *commonInfos) checkStandardFiles() {
	// In this loop only the standard files are used (not receipt or pending files)
	for path, fileinfo := range commonInfos.mapPathFileinfo {
		anchor := new(woleetapi.Anchor)
		hash, errHash := helpers.HashFile(path)
		if errHash != nil {
			errHandlerExitOnError(errHash, commonInfos.runParameters.ExitOnError)
			continue
		}
		tagsSlice := make([]string, 0)
		var tags []string
		if !(strings.HasPrefix(path, commonInfos.runParameters.Directory) && strings.HasSuffix(path, fileinfo.Name())) {
			log.Errorf("Unable to extract tags form the path: %s\n", path)
			if commonInfos.runParameters.ExitOnError {
				os.Exit(1)
			}
			continue
		}
		tags = strings.Split(strings.TrimSuffix(strings.TrimPrefix(path, commonInfos.runParameters.Directory), fileinfo.Name()), string(os.PathSeparator))
		for i := range tags {
			if !(strings.Contains(tags[i], " ") || strings.EqualFold(tags[i], "")) {
				tagsSlice = append(tagsSlice, tags[i])
			}
		}

		if !commonInfos.runParameters.Signature {
			anchor.Name = fileinfo.Name()
			anchor.Hash = hash
			anchor.Tags = tagsSlice
			anchor.Public = &commonInfos.runParameters.InvertPrivate
		} else {
			signatureGet, errSignatureGet := commonInfos.backendkitClient.GetSignature(hash, commonInfos.runParameters.BackendkitPubKey)
			if errSignatureGet != nil {
				errHandlerExitOnError(errSignatureGet, commonInfos.runParameters.ExitOnError)
				continue
			}
			anchor.Name = fileinfo.Name()
			anchor.Tags = tagsSlice
			anchor.Public = &commonInfos.runParameters.InvertPrivate
			anchor.PubKey = signatureGet.PubKey
			anchor.SignedHash = signatureGet.SignedHash
			anchor.Signature = signatureGet.Signature
			anchor.IdentityURL = signatureGet.IdentityURL
		}
		commonInfos.postAnchorCreatePendingFile(anchor, path)
	}
}

func (commonInfos *commonInfos) postAnchorCreatePendingFile(anchor *woleetapi.Anchor, path string) {
	anchorPost, errAnchorPost := commonInfos.client.PostAnchor(anchor)
	if errAnchorPost != nil {
		errHandlerExitOnError(errAnchorPost, commonInfos.runParameters.ExitOnError)
		return
	}
	pendingReceipt := new(woleetapi.Receipt)
	if !commonInfos.runParameters.Signature {
		pendingReceipt.TargetHash = anchorPost.Hash
	} else {
		pendingReceipt.Signature.SignedHash = anchorPost.SignedHash
	}
	pendingJSON, errPendingJSON := json.Marshal(pendingReceipt)
	if errPendingJSON != nil {
		errHandlerExitOnError(errPendingJSON, commonInfos.runParameters.ExitOnError)
		return
	}
	currentSuffix := helpers.SuffixAnchorPending
	if commonInfos.runParameters.Signature {
		currentSuffix = helpers.SuffixSignaturePending
	}
	errWriteFile := ioutil.WriteFile(path+"-"+anchorPost.Id+currentSuffix, pendingJSON, 0644)
	if errWriteFile != nil {
		errHandlerExitOnError(errWriteFile, commonInfos.runParameters.ExitOnError)
		return
	}
	if !commonInfos.runParameters.Signature {
		log.Infof("Anchoring file: %s\n", path)
	} else {
		log.Infof("Signing file: %s\n", path)
	}
}

func (commonInfos *commonInfos) getReceipts(mapPending map[string]os.FileInfo) {
	for path, fileinfo := range mapPending {
		anchorNameInfo, erranchorNameInfo := helpers.GetAnchorIDFromName(fileinfo)
		if erranchorNameInfo != nil {
			errHandlerExitOnError(erranchorNameInfo, commonInfos.runParameters.ExitOnError)
			continue
		}

		anchorGet, errAnchorGet := commonInfos.client.GetAnchor(anchorNameInfo.AnchorID)
		if errAnchorGet != nil {
			errHandlerExitOnError(errAnchorGet, commonInfos.runParameters.ExitOnError)
			continue
		}

		originalFilePath := strings.TrimSuffix(path, "-"+anchorNameInfo.AnchorID+anchorNameInfo.Suffix)
		if !strings.EqualFold(anchorGet.Status, "CONFIRMED") {
			log.Infof("Proof for file: %s with anchorID: %s not yet available\n", originalFilePath, anchorNameInfo.AnchorID)
			continue
		}

		// If the anchor is confirmed, we get its receipt and we deletes the old pending file
		log.Infof("Retrieving proof for file %s\n", originalFilePath)
		currentSuffix := helpers.SuffixAnchorReceipt
		if commonInfos.runParameters.Signature {
			currentSuffix = helpers.SuffixSignatureReceipt
		}
		errGetReceipt := commonInfos.client.GetReceiptToFile(anchorNameInfo.AnchorID, strings.TrimSuffix(path, anchorNameInfo.Suffix)+currentSuffix)
		if errGetReceipt != nil {
			errHandlerExitOnError(errGetReceipt, commonInfos.runParameters.ExitOnError)
			continue
		}
		errRemove := os.Remove(path)
		if errRemove != nil {
			errHandlerExitOnError(errRemove, commonInfos.runParameters.ExitOnError)
		}
		log.Infof("Done\n")
	}
}
