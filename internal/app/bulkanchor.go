package app

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"

	"github.com/clarketm/json"
	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/api"
	"github.com/woleet/woleet-cli/pkg/helpers"
	"github.com/woleet/woleet-cli/pkg/models/woleetapi"
)

func BulkAnchor(runParameters *RunParameters, logInput *logrus.Logger) int {
	commonInfos := initCommonInfos(runParameters)

	log = logInput

	commonInfos.client = api.GetNewClient(runParameters.BaseURL, runParameters.Token)

	var err error
	if runParameters.IsFS {
		commonInfos.mapPathFileName, err = helpers.ExploreDirectory(runParameters.Directory, runParameters.Recursive, runParameters.Filter, log)
	}

	if runParameters.IsS3 {
		commonInfos.mapPathFileName = helpers.ExploreS3(runParameters.S3Client, runParameters.S3Bucket, runParameters.Recursive, runParameters.Filter, log)
	}

	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	if !commonInfos.runParameters.Signature {
		commonInfos.pending, commonInfos.receipt, _, _ = helpers.SeparateAll(commonInfos.mapPathFileName)
	} else {
		_, _, commonInfos.pending, commonInfos.receipt = helpers.SeparateAll(commonInfos.mapPathFileName)
	}

	if runParameters.Signature {
		// Check Woleet.ID Server connection
		commonInfos.widsClient = api.GetNewClient(commonInfos.runParameters.IDServerSignURL, commonInfos.runParameters.IDServerToken)
		if commonInfos.runParameters.IDServerUnsecureSSL {
			commonInfos.widsClient.DisableSSLVerification()
		}
		checkWIDSConnectionPubKey(commonInfos)

		user, errUser := commonInfos.widsClient.GetUser()
		errHandlerExitOnError(errUser, commonInfos.runParameters.ExitOnError)
		commonInfos.runParameters.SignedIdentity = buildSignedIdentityString(user)

		if commonInfos.runParameters.SignedIdentity != "" {
			commonInfos.widsClient.GetServerConfig()
			config, errConfig := commonInfos.widsClient.GetServerConfig()
			errHandlerExitOnError(errConfig, commonInfos.runParameters.ExitOnError)
			commonInfos.runParameters.SignedIssuerDomain = buildSignedIssuerDomainString(config)
		}
	}

	commonInfos.splitPendingReceipt()
	commonInfos.getReceipts(commonInfos.pending)
	if !commonInfos.runParameters.Prune {
		commonInfos.getReceipts(commonInfos.pendingToDelete)
	} else {
		for path := range commonInfos.pendingToDelete {
			log.WithFields(logrus.Fields{
				"file": path,
			}).Infoln("Deleting old pending file")
			errRemove := commonInfos.removeFile(path)
			errHandlerExitOnError(errRemove, commonInfos.runParameters.ExitOnError)
		}
		for path := range commonInfos.receiptToDelete {
			log.WithFields(logrus.Fields{
				"file": path,
			}).Infoln("Deleting old receipt file")
			errRemove := commonInfos.removeFile(path)
			errHandlerExitOnError(errRemove, commonInfos.runParameters.ExitOnError)
		}
	}
	commonInfos.checkStandardFiles()
	return returnValue
}

func (commonInfos *commonInfos) splitPendingReceipt() {
	for path, fileinfo := range commonInfos.pending {
		errHandlerExitOnError(commonInfos.sortFile(path, fileinfo, true, false), commonInfos.runParameters.ExitOnError)
	}
	for path, fileinfo := range commonInfos.receipt {
		errHandlerExitOnError(commonInfos.sortFile(path, fileinfo, false, true), commonInfos.runParameters.ExitOnError)
	}
}

func (commonInfos *commonInfos) sortFile(path string, fileName string, pending bool, receipt bool) error {
	if receipt && commonInfos.runParameters.FixReceipts {
		errFix := commonInfos.fixReceipt(path)
		if errFix != nil {
			return errFix
		}
	}

	anchorNameInfo, erranchorNameInfo := helpers.GetAnchorIDFromName(fileName)
	if erranchorNameInfo != nil {
		return erranchorNameInfo
	}

	// Extracting the file's original path by the name of the pending/receipt
	originalFilePath := strings.TrimSuffix(path, "-"+anchorNameInfo.AnchorID+anchorNameInfo.Suffix)

	_, exists := commonInfos.mapPathFileName[originalFilePath]
	if !exists {
		if commonInfos.runParameters.Prune {
			if pending {
				commonInfos.pendingToDelete[path] = fileName
				delete(commonInfos.pending, path)
			}
			if receipt {
				commonInfos.receiptToDelete[path] = fileName
				delete(commonInfos.receipt, path)
			}
		}
		return nil
	}

	// If there is no strict mode, there is nothing to check and the file will not be reanchored
	if !commonInfos.runParameters.Strict {
		delete(commonInfos.mapPathFileName, originalFilePath)
		return nil
	}

	// If strict mode is actived, we check that the hash of the file
	// is the same as the one in the pending/receipt
	// If the file does not exists anymore and the prune mode is set the file will be deleted
	// if the prune mode is not set the file will be converted to a proper receipt

	receiptJSON, errReceiptJSON := commonInfos.readFile(path)
	if errReceiptJSON != nil {
		return errReceiptJSON
	}

	var receiptUnmarshalled woleetapi.Receipt
	errUnmarshal := json.Unmarshal(receiptJSON, &receiptUnmarshalled)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	hash, errHash := commonInfos.getHash(originalFilePath)
	if errHash != nil {
		return errHash
	}

	// In case of simple anchoring:
	//   If the hashes are equal, there is nothing to do
	// In case of signature:
	//   If the signedhashs and pubkeys are equals, there is nothing to do
	if !commonInfos.runParameters.Signature {
		if strings.EqualFold(hash, receiptUnmarshalled.TargetHash) {
			// File already anchored and valid
			delete(commonInfos.mapPathFileName, originalFilePath)
			return nil
		}
	} else {
		if strings.EqualFold(hash, receiptUnmarshalled.Signature.SignedHash) && strings.EqualFold(commonInfos.runParameters.IDServerPubKey, receiptUnmarshalled.Signature.PubKey) {
			// File signed and signature is up-to-date with current PubKey anchored and valid
			delete(commonInfos.mapPathFileName, originalFilePath)
			return nil
		}
	}

	// If they are not and there is a prune flag, the old pending file will be marked for deletion
	if commonInfos.runParameters.Prune {
		if pending {
			commonInfos.pendingToDelete[path] = fileName
			delete(commonInfos.pending, path)
		}
		if receipt {
			commonInfos.receiptToDelete[path] = fileName
			delete(commonInfos.receipt, path)
		}
	}
	return nil
}

func (commonInfos *commonInfos) checkStandardFiles() {
	// In this loop only the standard files are used (not receipt or pending files)
	for path, fileName := range commonInfos.mapPathFileName {
		anchor := new(woleetapi.Anchor)

		hash, errHash := commonInfos.getHash(path)
		if errHash != nil {
			errHandlerExitOnError(errHash, commonInfos.runParameters.ExitOnError)
			continue
		}

		if !commonInfos.runParameters.Signature {
			anchor.Name = fileName
			anchor.Hash = hash
			anchor.Public = &commonInfos.runParameters.InvertPrivate
		} else {
			messageToSign := hash

			if commonInfos.runParameters.SignedIdentity+commonInfos.runParameters.SignedIssuerDomain != "" {
				signatureHash := sha256.Sum256([]byte(hash + commonInfos.runParameters.SignedIdentity + commonInfos.runParameters.SignedIssuerDomain))
				messageToSign = hex.EncodeToString(signatureHash[:])
			}

			signatureGet, errSignatureGet := commonInfos.widsClient.GetSignature(messageToSign, commonInfos.runParameters.IDServerPubKey)
			if errSignatureGet != nil {
				errHandlerExitOnError(errSignatureGet, commonInfos.runParameters.ExitOnError)
				continue
			}
			anchor.Name = fileName
			anchor.Public = &commonInfos.runParameters.InvertPrivate
			anchor.PubKey = signatureGet.PubKey
			anchor.SignedHash = hash
			anchor.Signature = signatureGet.Signature
			anchor.IdentityURL = signatureGet.IdentityURL
			if commonInfos.runParameters.SignedIdentity != "" {
				anchor.SignedIdentity = commonInfos.runParameters.SignedIdentity
				if commonInfos.runParameters.SignedIssuerDomain != "" {
					anchor.SignedIssuerDomain = commonInfos.runParameters.SignedIssuerDomain
				}
			}
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
		pendingReceipt.Signature.PubKey = anchorPost.PubKey
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
	errWrite := commonInfos.writeFile(path+"-"+anchorPost.Id+currentSuffix, pendingJSON)
	if errWrite != nil {
		errHandlerExitOnError(errWrite, commonInfos.runParameters.ExitOnError)
		return
	}
	if !commonInfos.runParameters.Signature {
		log.WithFields(logrus.Fields{
			"file": path,
		}).Infoln("Anchoring file")
	} else {
		log.WithFields(logrus.Fields{
			"file": path,
		}).Infoln("Signing file")
	}
}

func (commonInfos *commonInfos) getReceipts(mapPending map[string]string) {
	for path, fileName := range mapPending {
		anchorNameInfo, erranchorNameInfo := helpers.GetAnchorIDFromName(fileName)
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
			log.WithFields(logrus.Fields{
				"anchorID":     anchorNameInfo.AnchorID,
				"originalFile": originalFilePath,
			}).Infoln("Proof not available yet")
			continue
		}

		// If the anchor is confirmed, we get its receipt and we deletes the old pending file
		currentSuffix := helpers.SuffixAnchorReceipt
		if commonInfos.runParameters.Signature {
			currentSuffix = helpers.SuffixSignatureReceipt
		}
		receiptPath := strings.TrimSuffix(path, anchorNameInfo.Suffix) + currentSuffix
		receipt, errReceipt := commonInfos.client.GetReceipt(anchorNameInfo.AnchorID)
		if errReceipt != nil {
			errHandlerExitOnError(errReceipt, commonInfos.runParameters.ExitOnError)
			continue
		}
		receiptJSON, errReceiptJSON := json.Marshal(receipt)
		if errReceiptJSON != nil {
			errHandlerExitOnError(errReceiptJSON, commonInfos.runParameters.ExitOnError)
			return
		}
		commonInfos.writeFile(receiptPath, receiptJSON)
		errRemove := commonInfos.removeFile(path)
		if errRemove != nil {
			errHandlerExitOnError(errRemove, commonInfos.runParameters.ExitOnError)
		}
		log.WithFields(logrus.Fields{
			"originalFile": originalFilePath,
			"proofFile":    receiptPath,
		}).Infoln("Proof retrieved")
	}
}

func (commonInfos *commonInfos) fixReceipt(path string) error {
	receiptJSON, errReceiptJSON := commonInfos.readFile(path)
	if errReceiptJSON != nil {
		return errReceiptJSON
	}

	var receiptUnmarshalled woleetapi.Receipt
	errUnmarshal := json.Unmarshal(receiptJSON, &receiptUnmarshalled)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	receiptMarshalled, errReceiptMarshalled := json.Marshal(receiptUnmarshalled)
	if errReceiptMarshalled != nil {
		return errReceiptMarshalled
	}

	if !bytes.Equal(receiptJSON, receiptMarshalled) {
		log.WithFields(logrus.Fields{
			"proofFile": path,
		}).Infoln("Fixing receipt")

		errWrite := commonInfos.writeFile(path, receiptMarshalled)
		if errWrite != nil {
			return errWrite
		}
	}
	return nil
}
