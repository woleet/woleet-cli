package app

import (
	"log"
	"os"
	"strings"

	"github.com/woleet/woleet-cli/pkg/api"
	"github.com/woleet/woleet-cli/pkg/helpers"
)

func ExportReceipts(token string, url string, exportDirectory string, unixEpochLimit int64, exitOnError bool) {
	stdLogger := log.New(os.Stdout, "woleet-cli ", log.LstdFlags)
	errLogger := log.New(os.Stderr, "woleet-cli ", log.LstdFlags)
	client := api.GetNewClient(url, token)
	client.SetCustomLogger(errLogger)

	end := false

	for pageIndex := 0; !end; pageIndex++ {
		anchors, errAnchors := client.GetAnchors(pageIndex, pageSize, "DESC", "created")
		if errAnchors != nil {
			errHandlerExitOnError(errAnchors, errLogger, true)
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
				stdLogger.Printf("INFO: Proof for anchor: %s named: %s is already on disk\n", anchor.Id, anchor.Name)
				continue
			}
			if !strings.EqualFold(anchor.Status, "CONFIRMED") {
				stdLogger.Printf("INFO: Proof for anchor: %s named: %s not available yet\n", anchor.Id, anchor.Name)
				continue
			}
			stdLogger.Printf("INFO: Retrieving proof for anchor: %s named: %s\n", anchor.Id, anchor.Name)
			errGetReceipt := client.GetReceiptToFile(anchor.Id, receiptPath)
			if errGetReceipt != nil {
				if _, err := os.Stat(receiptPath); !os.IsNotExist(err) {
					errRemove := os.Remove(receiptPath)
					if errRemove != nil {
						errHandlerExitOnError(errRemove, errLogger, exitOnError)
					}
				}
				errHandlerExitOnError(errAnchors, errLogger, exitOnError)
				continue
			}
			stdLogger.Printf("INFO: Done\n")
		}
	}
}
