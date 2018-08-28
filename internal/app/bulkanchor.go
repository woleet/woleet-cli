package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/woleet/woleet-cli/pkg/models"

	"github.com/woleet/woleet-cli/pkg/api"
	"github.com/woleet/woleet-cli/pkg/helpers"
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
	stdLogger        *log.Logger
	errLogger        *log.Logger
	client           *api.Client
	backendkitClient *api.Client
	mapPathFileinfo  map[string]os.FileInfo
	pending          map[string]os.FileInfo
	receipt          map[string]os.FileInfo
	runParameters    RunParameters
}

func (commonInfos *commonInfos) errHandlerExitOnError(err error) {
	if err != nil {
		commonInfos.errLogger.Printf("ERROR: %v\n", err)
		if commonInfos.runParameters.ExitOnError {
			os.Exit(1)
		}
	}
}

func BulkAnchor(runParameters *RunParameters) {
	commonInfos := new(commonInfos)
	commonInfos.runParameters = *runParameters

	commonInfos.stdLogger = log.New(os.Stdout, "woleet-cli ", log.LstdFlags)
	commonInfos.errLogger = log.New(os.Stderr, "woleet-cli ", log.LstdFlags)
	commonInfos.client = api.GetNewClient(runParameters.BaseURL, runParameters.Token)
	commonInfos.client.SetCustomLogger(commonInfos.errLogger)

	var err error
	commonInfos.mapPathFileinfo, err = helpers.Explore(runParameters.Directory)
	if err != nil {
		commonInfos.errLogger.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	commonInfos.mapPathFileinfo, commonInfos.pending, commonInfos.receipt = helpers.Separate(commonInfos.mapPathFileinfo, commonInfos.runParameters.Signature, commonInfos.runParameters.Strict)

	if runParameters.Signature {
		//Check backendkit connection
		commonInfos.backendkitClient = api.GetNewClient(commonInfos.runParameters.BackendkitSignURL, commonInfos.runParameters.BackendkitToken)
		if commonInfos.runParameters.UnsecureSSL {
			commonInfos.backendkitClient.DisableSSLVerification()
		}
		commonInfos.backendkitClient.CheckBackendkitConnection(commonInfos.errLogger)
		commonInfos.backendkitClient.SetCustomLogger(commonInfos.errLogger)
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
			commonInfos.errHandlerExitOnError(erranchorNameInfo)
			continue
		}
		anchorGet, errAnchorGet := commonInfos.client.GetAnchor(anchorNameInfo.AnchorID)
		if errAnchorGet != nil {
			commonInfos.errHandlerExitOnError(errAnchorGet)
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
				commonInfos.errHandlerExitOnError(errhash)
				continue
			}
			// If the hashes corresponds, we removes the original file from the filelist
			// Doing so it will not be reanchored
			if (!commonInfos.runParameters.Signature && strings.EqualFold(hash, anchorGet.Hash)) || (commonInfos.runParameters.Signature && strings.EqualFold(hash, anchorGet.SignedHash)) {
				delete(commonInfos.mapPathFileinfo, originalFilePath)
			} else if commonInfos.runParameters.Prune {
				commonInfos.stdLogger.Printf("INFO: Prune enabled, deleting old pending file: %s\n", path)
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
			commonInfos.stdLogger.Printf("INFO: Proof for file: %s with anchorID: %s not yet available\n", originalFilePath, anchorNameInfo.AnchorID)
		} else {
			// If the anchor is confirmed, we get its receipt and we deletes the old pending file
			commonInfos.stdLogger.Printf("INFO: Retrieving proof for file %s\n", originalFilePath)
			currentSuffix := helpers.SuffixAnchorReceipt
			if commonInfos.runParameters.Signature {
				currentSuffix = helpers.SuffixSignatureReceipt
			}
			errGetReceipt := commonInfos.client.GetReceiptToFile(anchorNameInfo.AnchorID, strings.TrimSuffix(path, anchorNameInfo.Suffix)+currentSuffix)
			if errGetReceipt != nil {
				commonInfos.errHandlerExitOnError(errGetReceipt)
				continue
			}
			errRemove := os.Remove(path)
			if errRemove != nil {
				commonInfos.errHandlerExitOnError(errRemove)
			}
			commonInfos.stdLogger.Printf("INFO: Done\n")
		}
	}
}

func (commonInfos *commonInfos) checkReceipts() {
	// In this loop only receipt files are used
	for path, fileinfo := range commonInfos.receipt {
		// Extracting the file's original path by the name of the pending file
		anchorNameInfo, erranchorNameInfo := helpers.GetAnchorIDFromName(fileinfo)
		if erranchorNameInfo != nil {
			commonInfos.errHandlerExitOnError(erranchorNameInfo)
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
			commonInfos.errHandlerExitOnError(errHash)
			continue
		}
		// Get hash from receipt
		receiptJSON, errReceiptJSON := ioutil.ReadFile(path)
		if errReceiptJSON != nil {
			commonInfos.errHandlerExitOnError(errReceiptJSON)
			continue
		}
		var receiptUnmarshalled models.Receipt
		json.Unmarshal(receiptJSON, &receiptUnmarshalled)
		if (!commonInfos.runParameters.Signature && strings.EqualFold(hash, receiptUnmarshalled.TargetHash)) || (commonInfos.runParameters.Signature && strings.EqualFold(hash, receiptUnmarshalled.Signature.SignedHash)) {
			// If the hashes corresponds, we removes the original file from the filelist
			// Doing so it will not be reanchored
			delete(commonInfos.mapPathFileinfo, originalFilePath)
		} else if commonInfos.runParameters.Prune {
			// If prune is defined and the hashes does not correspond
			// we remove the file as it does not correspond with the current hash
			os.Remove(path)
			commonInfos.stdLogger.Printf("INFO: Strict-prune mode enabled, deleting: %s\n", path)
		}
	}
}

func (commonInfos *commonInfos) checkStandardFiles() {
	// In this loop only the standard files are used (not receipt or pending files)
	for path, fileinfo := range commonInfos.mapPathFileinfo {
		anchor := new(models.Anchor)
		hash, errHash := helpers.HashFile(path)
		if errHash != nil {
			commonInfos.errHandlerExitOnError(errHash)
			continue
		}
		tagsSlice := make([]string, 0)
		var tags []string
		if !(strings.HasPrefix(path, commonInfos.runParameters.Directory) && strings.HasSuffix(path, fileinfo.Name())) {
			commonInfos.errLogger.Printf("ERROR: Unable to extract tags form the path: %s\n", path)
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
				commonInfos.errHandlerExitOnError(errSignatureGet)
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

func (commonInfos *commonInfos) postAnchorCreatePendingFile(anchor *models.Anchor, path string) {
	anchorPost, errAnchorPost := commonInfos.client.PostAnchor(anchor)
	if errAnchorPost != nil {
		commonInfos.errHandlerExitOnError(errAnchorPost)
		return
	}
	pendingReceipt := new(models.Receipt)
	if !commonInfos.runParameters.Signature {
		pendingReceipt.TargetHash = anchorPost.Hash
	} else {
		pendingReceipt.Signature.SignedHash = anchorPost.SignedHash
	}
	pendingJSON, errPendingJSON := json.Marshal(pendingReceipt)
	if errPendingJSON != nil {
		commonInfos.errHandlerExitOnError(errPendingJSON)
		return
	}
	currentSuffix := helpers.SuffixAnchorPending
	if commonInfos.runParameters.Signature {
		currentSuffix = helpers.SuffixSignaturePending
	}
	errWriteFile := ioutil.WriteFile(path+"-"+anchorPost.Id+currentSuffix, pendingJSON, 0644)
	if errWriteFile != nil {
		commonInfos.errHandlerExitOnError(errWriteFile)
		return
	}
	if !commonInfos.runParameters.Signature {
		commonInfos.stdLogger.Printf("INFO: Anchoring file: %s\n", path)
	} else {
		commonInfos.stdLogger.Printf("INFO: Signing file: %s\n", path)
	}
}
