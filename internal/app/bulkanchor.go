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

type RunParameters struct {
	Signature         bool
	ExitOnError       bool
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
	receipt          map[string]os.FileInfo
	runParameters    RunParameters
}

func BulkAnchor(runParameters *RunParameters, logInput *logrus.Logger) {
	commonInfos := new(commonInfos)
	commonInfos.runParameters = *runParameters

	log = logInput

	commonInfos.client = api.GetNewClient(runParameters.BaseURL, runParameters.Token)

	var err error
	commonInfos.mapPathFileinfo, err = helpers.ExploreDirectory(runParameters.Directory, log)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	commonInfos.mapPathFileinfo, commonInfos.pending, commonInfos.receipt = helpers.Separate(commonInfos.mapPathFileinfo, commonInfos.runParameters.Signature, commonInfos.runParameters.Strict)

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

	commonInfos.checkPendings()
	commonInfos.checkReceipts()
	commonInfos.checkStandardFiles()
}

func (commonInfos *commonInfos) checkPendings() {
	// In this loop only pending files are used
	for path, fileinfo := range commonInfos.pending {
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

		// Extracting the file's original path by the name of the pending file
		originalFilePath := strings.TrimSuffix(path, "-"+anchorNameInfo.AnchorID+anchorNameInfo.Suffix)

		// If strict mode is actived, we check that the hash of the file
		// is the same as the one in the pending file
		if commonInfos.runParameters.Strict {
			_, exists := commonInfos.mapPathFileinfo[originalFilePath]
			if !exists {
				continue
			}
			hash, errhash := helpers.HashFile(originalFilePath)
			if errhash != nil {
				errHandlerExitOnError(errhash, commonInfos.runParameters.ExitOnError)
				continue
			}
			// If the hashes corresponds, we removes the original file from the filelist
			// Doing so it will not be reanchored
			if (!commonInfos.runParameters.Signature && strings.EqualFold(hash, anchorGet.Hash)) || (commonInfos.runParameters.Signature && strings.EqualFold(hash, anchorGet.SignedHash)) {
				delete(commonInfos.mapPathFileinfo, originalFilePath)
			} else if commonInfos.runParameters.Prune {
				log.Infof("Prune enabled, deleting old pending file: %s\n", path)
				// If prune is specified, we remove the old pending file that
				// does not correspond anymore to the original file
				os.Remove(path)
			}
		} else {
			// if strict mode is not enable we do not want to rescan
			// the file so we remove the original file from the filelist
			delete(commonInfos.mapPathFileinfo, originalFilePath)
		}
		if !strings.EqualFold(anchorGet.Status, "CONFIRMED") {
			log.Infof("Proof for file: %s with anchorID: %s not yet available\n", originalFilePath, anchorNameInfo.AnchorID)
		} else {
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
}

func (commonInfos *commonInfos) checkReceipts() {
	// In this loop only receipt files are used
	for path, fileinfo := range commonInfos.receipt {
		// Extracting the file's original path by the name of the pending file
		anchorNameInfo, erranchorNameInfo := helpers.GetAnchorIDFromName(fileinfo)
		if erranchorNameInfo != nil {
			errHandlerExitOnError(erranchorNameInfo, commonInfos.runParameters.ExitOnError)
			continue
		}
		// Extracting the file's original path by the name of the receipt file
		originalFilePath := strings.TrimSuffix(path, "-"+anchorNameInfo.AnchorID+anchorNameInfo.Suffix)
		// if strict mode is not enable we do not want to rescan
		// the file so we remove the original file from the filelist
		if !commonInfos.runParameters.Strict {
			delete(commonInfos.mapPathFileinfo, originalFilePath)
			continue
		}
		_, exists := commonInfos.mapPathFileinfo[originalFilePath]
		if !exists {
			continue
		}
		// If strict mode is actived, we check that the hash of the file
		// is the same as the one in the receipt file
		hash, errHash := helpers.HashFile(originalFilePath)
		if errHash != nil {
			errHandlerExitOnError(errHash, commonInfos.runParameters.ExitOnError)
			continue
		}
		// Get hash from receipt
		receiptJSON, errReceiptJSON := ioutil.ReadFile(path)
		if errReceiptJSON != nil {
			errHandlerExitOnError(errReceiptJSON, commonInfos.runParameters.ExitOnError)
			continue
		}
		var receiptUnmarshalled woleetapi.Receipt
		json.Unmarshal(receiptJSON, &receiptUnmarshalled)
		if (!commonInfos.runParameters.Signature && strings.EqualFold(hash, receiptUnmarshalled.TargetHash)) || (commonInfos.runParameters.Signature && strings.EqualFold(hash, receiptUnmarshalled.Signature.SignedHash)) {
			// If the hashes corresponds, we removes the original file from the filelist
			// Doing so it will not be reanchored
			delete(commonInfos.mapPathFileinfo, originalFilePath)
		} else if commonInfos.runParameters.Prune {
			// If prune is defined and the hashes does not correspond
			// we remove the file as it does not correspond with the current hash
			os.Remove(path)
			log.Infof("Strict-prune mode enabled, deleting: %s\n", path)
		}
	}
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
