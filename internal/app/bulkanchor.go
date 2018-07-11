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

func BulkAnchor(baseURL string, token string, directory string, exitOnError bool, private bool, strict bool, prune bool) {
	invertPrivate := !private
	stdLogger := log.New(os.Stdout, "woleet-cli ", log.LstdFlags)
	errLogger := log.New(os.Stderr, "woleet-cli ", log.LstdFlags)
	client := api.GetNewClient(baseURL, token)

	mapPathFileinfo, errFiles := helpers.Explore(directory)
	if errFiles != nil {
		errLogger.Printf("ERROR: %v\n", errFiles)
		os.Exit(1)
	}

	mapPathFileinfo, pending, receipt := helpers.Separate(mapPathFileinfo, strict)

	// In this loop only pending files are used
	for path, fileinfo := range pending {
		anchorID, errAnchorID := helpers.GetAnchorIDFromName(fileinfo)
		if errAnchorID != nil {
			errLogger.Printf("ERROR: %v\n", errAnchorID)
			if exitOnError {
				os.Exit(1)
			}
		} else {
			anchorGet, errAnchorGet := client.GetAnchor(anchorID)
			if errAnchorGet != nil {
				errLogger.Printf("ERROR: %v\n", errAnchorGet)
				if exitOnError {
					os.Exit(1)
				}
			} else {
				// Extracting the file's original path by the name of the pending file
				originalFilePath := strings.TrimSuffix(path, "-"+anchorID+".pending.json")
				// If strict mode is actived, we check that the hash of the file
				// is the same as the one in the pending file
				if strict {
					_, exists := mapPathFileinfo[originalFilePath]
					if exists {
						actualHash, erractualHash := helpers.HashFile(originalFilePath)
						if erractualHash != nil {
							errLogger.Printf("ERROR: %v\n", erractualHash)
							if exitOnError {
								os.Exit(1)
							}
						} else {
							// If the hashes corresponds, we removes the original file from the filelist
							// Doing so it will not be reanchored
							if strings.EqualFold(actualHash, anchorGet.Hash) {
								delete(mapPathFileinfo, originalFilePath)
							} else if prune {
								// If prune is specified, we remove the old pending file that
								// does not correspond anymore to the original file
								os.Remove(path)
							}
						}
					}
				} else {
					// if strict mode is not enable we do not want to rescan
					// the file so we remove the original file from the filelist
					delete(mapPathFileinfo, originalFilePath)
				}
				if !strings.EqualFold(anchorGet.Status, "CONFIRMED") {
					stdLogger.Printf("INFO: %s not yet available\n", path)
				} else {
					// If the anchor is confirmed, we get its receipt and we delets the old pending file
					errReceipt := client.GetReceiptToFile(anchorID, strings.TrimSuffix(path, ".pending.json")+".receipt.json")
					if errReceipt != nil {
						errLogger.Printf("ERROR: %v\n", errReceipt)
						if exitOnError {
							os.Exit(1)
						}
					} else {
						errRemoval := os.Remove(path)
						if errRemoval != nil {
							errLogger.Printf("ERROR: %v\n", errRemoval)
							if exitOnError {
								os.Exit(1)
							}
						}
					}
				}
			}
		}
	}

	// In this loop only receipt files are used
	for path, fileinfo := range receipt {
		// Extracting the file's original path by the name of the pending file
		anchorID, errAnchorID := helpers.GetAnchorIDFromName(fileinfo)
		if errAnchorID != nil {
			errLogger.Printf("ERROR: %v\n", errAnchorID)
			if exitOnError {
				os.Exit(1)
			}
		} else {
			// Extracting the file's original path by the name of the receipt file
			originalFilePath := strings.TrimSuffix(path, "-"+anchorID+".receipt.json")
			// if strict mode is not enable we do not want to rescan
			// the file so we remove the original file from the filelist
			if !strict {
				delete(mapPathFileinfo, originalFilePath)
			} else {
				_, exists := mapPathFileinfo[originalFilePath]
				if exists {
					// If strict mode is actived, we check that the hash of the file
					// is the same as the one in the receipt file
					hash, errHash := helpers.HashFile(originalFilePath)
					if errHash != nil {
						errLogger.Printf("ERROR: %v\n", errHash)
						if exitOnError {
							os.Exit(1)
						}
					}
					// Get hash from receipt
					receiptJSON, errFile := ioutil.ReadFile(path)
					if errFile != nil {
						errLogger.Printf("ERROR: %v\n", errFile)
						if exitOnError {
							os.Exit(1)
						}
					} else {
						var receiptUnmarshalled models.Receipt
						json.Unmarshal(receiptJSON, &receiptUnmarshalled)
						if strings.EqualFold(hash, receiptUnmarshalled.TargetHash) {
							// If the hashes corresponds, we removes the original file from the filelist
							// Doing so it will not be reanchored
							delete(mapPathFileinfo, originalFilePath)
						} else if prune {
							// If prune is defined and the hashes does not correspond
							// we remove the file as it does not correspond with the current hash
							os.Remove(path)
						}
					}
				}
			}
		}
	}

	// In this loop only the standard files are used (not receipt or pending files)
	for path, fileinfo := range mapPathFileinfo {
		anchor := new(models.Anchor)
		hash, errHash := helpers.HashFile(path)
		if errHash != nil {
			errLogger.Printf("ERROR: %v\n", errHash)
			if exitOnError {
				os.Exit(1)
			}
		} else {
			tagsSlice := make([]string, 0)
			var tags []string
			if !(strings.HasPrefix(path, directory) && strings.HasSuffix(path, fileinfo.Name())) {
				errLogger.Printf("ERROR: Unable to extract tags form the path: %s\n", path)
				if exitOnError {
					os.Exit(1)
				}
			} else {
				tags = strings.Split(strings.TrimSuffix(strings.TrimPrefix(path, directory), fileinfo.Name()), string(os.PathSeparator))
				for i := range tags {
					if !(strings.Contains(tags[i], " ") || strings.EqualFold(tags[i], "")) {
						tagsSlice = append(tagsSlice, tags[i])
					}
				}
			}

			anchor.Name = fileinfo.Name()
			anchor.Hash = hash
			anchor.Tags = tagsSlice
			anchor.Public = &invertPrivate

			anchorPost, errAnchorPost := client.PostAnchor(anchor)
			if errAnchorPost != nil {
				errLogger.Printf("ERROR: %v\n", errAnchorPost)
				if exitOnError {
					os.Exit(1)
				}
			} else {
				pendingReceipt := new(models.Receipt)
				pendingReceipt.TargetHash = anchorPost.Hash
				pendingJSON, errPendingJSON := json.Marshal(pendingReceipt)
				if errPendingJSON != nil {
					errLogger.Printf("ERROR: %v\n", errPendingJSON)
					if exitOnError {
						os.Exit(1)
					}
				} else {
					errWrite := ioutil.WriteFile(path+"-"+anchorPost.Id+".pending.json", pendingJSON, 0644)
					if errWrite != nil {
						errLogger.Printf("ERROR: %v\n", errWrite)
						if exitOnError {
							os.Exit(1)
						}
					}
				}
			}
		}
	}
}
