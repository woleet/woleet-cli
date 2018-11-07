package app

import (
	"os"
	"strings"

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
			receiptPath := exportDirectory + string(os.PathSeparator) + anchor.Name + "-" + anchor.Id + currentSuffix
			if _, err := os.Stat(receiptPath); !os.IsNotExist(err) {
				log.Infof("Proof for anchor: %s named: %s is already on disk\n", anchor.Id, anchor.Name)
				continue
			}
			if !strings.EqualFold(anchor.Status, "CONFIRMED") {
				log.Infof("Proof for anchor: %s named: %s not available yet\n", anchor.Id, anchor.Name)
				continue
			}
			log.Infof("Retrieving proof for anchor: %s named: %s\n", anchor.Id, anchor.Name)
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
			log.Infof("Done\n")
		}
	}
}
