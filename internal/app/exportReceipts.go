package app

import (
	"os"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/api"
	"github.com/woleet/woleet-cli/pkg/helpers"
)

func ExportReceipts(token string, url string, exportDirectory string, unixEpochLimit int64, fixReceipts bool, exitOnError bool, logInput *logrus.Logger) int {
	log = logInput
	client := api.GetNewClient(url, token)

	end := false

	for pageIndex := 0; !end; pageIndex++ {
		anchors, errAnchors := client.GetAnchors(pageIndex, pageSize, "DESC", "created")
		if errAnchors != nil {
			errHandlerExitOnError(errAnchors, true)
		}
		if *anchors.Last == true {
			end = true
		}
		for _, anchor := range anchors.Content {
			if anchor.GetCreated()*1000000 < unixEpochLimit {
				end = true
				continue
			}
			currentSuffix := helpers.SuffixAnchorReceiptCurrent
			legacySuffix := helpers.SuffixAnchorReceiptLegacy
			if anchor.HasSignature() {
				currentSuffix = helpers.SuffixSignatureReceiptCurrent
				legacySuffix = helpers.SuffixSignatureReceiptLegacy
			}

			fields := logrus.Fields{}
			fields["anchorID"] = anchor.Id
			fields["anchor_Name"] = anchor.Name
			fields["File_Name"] = sanitize.BaseName(anchor.GetName()) + "-" + anchor.GetId() + currentSuffix

			currentReceiptPath := exportDirectory + string(os.PathSeparator) + sanitize.BaseName(anchor.GetName()) + "-" + anchor.GetId() + currentSuffix

			if fixReceipts {
				legacyReceiptPath := exportDirectory + string(os.PathSeparator) + sanitize.BaseName(anchor.GetName()) + "-" + anchor.GetId() + legacySuffix
				if _, err := os.Stat(legacyReceiptPath); !os.IsNotExist(err) {
					if _, err := os.Stat(currentReceiptPath); !os.IsNotExist(err) {
						log.WithFields(fields).Warnln("Renaming legacy file aborted, new file already present")
						continue
					}
					err := os.Rename(legacyReceiptPath, currentReceiptPath)
					if err != nil {
						log.WithFields(fields).Warnln("Renaming legacy file failed")
						continue
					}
				}
			}
			if _, err := os.Stat(currentReceiptPath); !os.IsNotExist(err) {
				log.WithFields(fields).Infoln("Proof already on disk")
				continue
			}
			if !strings.EqualFold(anchor.GetStatus(), "CONFIRMED") {
				log.WithFields(fields).Infoln("Proof not available yet")
				continue
			}
			errGetReceipt := client.GetReceiptToFile(anchor.GetId(), currentReceiptPath)
			if errGetReceipt != nil {
				if _, err := os.Stat(currentReceiptPath); !os.IsNotExist(err) {
					errRemove := os.Remove(currentReceiptPath)
					if errRemove != nil {
						errHandlerExitOnError(errRemove, exitOnError)
					}
				}
				errHandlerExitOnError(errAnchors, exitOnError)
				continue
			}
			log.WithFields(fields).Infoln("Proof retrived")
		}
	}
	return returnValue
}
