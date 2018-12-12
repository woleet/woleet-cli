package app

import (
	"os"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/api"
	"github.com/woleet/woleet-cli/pkg/helpers"
)

func ExportReceipts(token string, url string, exportDirectory string, unixEpochLimit int64, exitOnError bool, logInput *logrus.Logger) {
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
			if anchor.Created*1000000 < unixEpochLimit {
				end = true
				continue
			}
			currentSuffix := helpers.SuffixAnchorReceipt
			if anchor.Signature != "" {
				currentSuffix = helpers.SuffixSignatureReceipt
			}

			fields := logrus.Fields{}
			fields["anchorID"] = anchor.Id
			fields["anchor_Name"] = anchor.Name
			fields["File_Name"] = sanitize.BaseName(anchor.Name) + "-" + anchor.Id + currentSuffix

			receiptPath := exportDirectory + string(os.PathSeparator) + sanitize.BaseName(anchor.Name) + "-" + anchor.Id + currentSuffix
			if _, err := os.Stat(receiptPath); !os.IsNotExist(err) {
				log.WithFields(fields).Infoln("Proof already on disk")
				continue
			}
			if !strings.EqualFold(anchor.Status, "CONFIRMED") {
				log.WithFields(fields).Infoln("Proof not available yet")
				continue
			}
			errGetReceipt := client.GetReceiptToFile(anchor.Id, receiptPath)
			if errGetReceipt != nil {
				if _, err := os.Stat(receiptPath); !os.IsNotExist(err) {
					errRemove := os.Remove(receiptPath)
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
}
