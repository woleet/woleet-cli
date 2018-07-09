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

func BulkAnchor(baseURL string, token string, directory string, exitOnError bool, strict bool, prune bool) {
	stdLogger := log.New(os.Stdout, "woleet-cli ", log.LstdFlags)
	errLogger := log.New(os.Stderr, "woleet-cli ", log.LstdFlags)
	client := api.GetNewClient(baseURL, token)

	mapPathFileinfo, errFiles := helpers.Explore(directory)
	if errFiles != nil {
		errLogger.Printf("ERROR :%v\n", errFiles)
		os.Exit(1)
	}

	mapPathFileinfo, pending, receipt := helpers.Separate(mapPathFileinfo, strict)

	for path, fileinfo := range pending {
		anchorID, errAnchorID := helpers.GetAnchorIDFromName(fileinfo)
		if errAnchorID != nil {
			errLogger.Printf("ERROR :%v\n", errAnchorID)
			if exitOnError {
				os.Exit(1)
			}
		} else {
			anchorGet, errAnchorGet := client.GetAnchor(anchorID)
			if errAnchorGet != nil {
				errLogger.Printf("ERROR :%v\n", errAnchorGet)
				if exitOnError {
					os.Exit(1)
				}
			} else {
				originalFilePath := strings.TrimSuffix(path, "-"+anchorID+".pending.json")
				if strict {
					_, exists := mapPathFileinfo[originalFilePath]
					if exists {
						actualHash, erractualHash := helpers.HashFile(originalFilePath)
						if erractualHash != nil {
							errLogger.Printf("ERROR :%v\n", erractualHash)
							if exitOnError {
								os.Exit(1)
							}
						} else {
							if strings.EqualFold(actualHash, anchorGet.Hash) {
								delete(mapPathFileinfo, originalFilePath)
							} else if prune {
								os.Remove(path)
							}
						}
					}
				} else {
					delete(mapPathFileinfo, originalFilePath)
				}
				if !strings.EqualFold(anchorGet.Status, "CONFIRMED") {
					stdLogger.Printf("WARN : anchorID: %s not availaible yet\n", path)
				} else {
					errReceipt := client.GetReceiptToFile(anchorID, strings.TrimSuffix(path, ".pending.json")+".receipt.json")
					if errReceipt != nil {
						errLogger.Printf("ERROR :%v\n", errReceipt)
						if exitOnError {
							os.Exit(1)
						}
					} else {
						errRemoval := os.Remove(path)
						if errRemoval != nil {
							errLogger.Printf("ERROR :%v\n", errRemoval)
							if exitOnError {
								os.Exit(1)
							}
						}
					}
				}
			}
		}
	}

	for path, fileinfo := range receipt {
		anchorID, errAnchorID := helpers.GetAnchorIDFromName(fileinfo)
		if errAnchorID != nil {
			errLogger.Printf("ERROR :%v\n", errAnchorID)
			if exitOnError {
				os.Exit(1)
			}
		} else {
			originalFilePath := strings.TrimSuffix(path, "-"+anchorID+".receipt.json")
			if !strict {
				delete(mapPathFileinfo, originalFilePath)
			} else {
				_, exists := mapPathFileinfo[originalFilePath]
				if exists {
					hash, errHash := helpers.HashFile(originalFilePath)
					if errHash != nil {
						errLogger.Printf("ERROR :%v\n", errHash)
						if exitOnError {
							os.Exit(1)
						}
					}
					// Get hash from receipt
					receiptJSON, errFile := ioutil.ReadFile(path)
					if errFile != nil {
						errLogger.Printf("ERROR :%v\n", errFile)
						if exitOnError {
							os.Exit(1)
						}
					} else {
						var receiptUnmarshalled models.Receipt
						json.Unmarshal(receiptJSON, &receiptUnmarshalled)
						if strings.EqualFold(hash, receiptUnmarshalled.TargetHash) {
							delete(mapPathFileinfo, originalFilePath)
						} else if prune {
							os.Remove(path)
						}
					}
				}
			}
		}
	}

	for path, fileinfo := range mapPathFileinfo {
		anchor := new(models.Anchor)
		hash, errHash := helpers.HashFile(path)
		if errHash != nil {
			errLogger.Printf("ERROR :%v\n", errHash)
			if exitOnError {
				os.Exit(1)
			}
		} else {
			tagsSlice := make([]string, 0)
			var tags []string
			if !(strings.HasPrefix(path, directory) && strings.HasSuffix(path, fileinfo.Name())) {
				errLogger.Printf("ERROR : Unable to extract tags form the path: %s\n", path)
				if exitOnError {
					os.Exit(1)
				}
			} else {
				tags = strings.Split(strings.TrimSuffix(strings.TrimPrefix(path, directory), fileinfo.Name()), "/")
				for i := range tags {
					if !(strings.Contains(tags[i], " ") || strings.EqualFold(tags[i], "")) {
						tagsSlice = append(tagsSlice, tags[i])
					}
				}
			}

			anchor.Name = fileinfo.Name()
			anchor.Hash = hash
			anchor.Tags = tagsSlice

			anchorPost, errAnchorPost := client.PostAnchor(anchor)
			if errAnchorPost != nil {
				errLogger.Printf("ERROR :%v\n", errAnchorPost)
				if exitOnError {
					os.Exit(1)
				}
			} else {
				pendingReceipt := new(models.Receipt)
				pendingReceipt.TargetHash = anchorPost.Hash
				pendingJSON, errPendingJSON := json.Marshal(pendingReceipt)
				if errPendingJSON != nil {
					errLogger.Printf("ERROR :%v\n", errPendingJSON)
					if exitOnError {
						os.Exit(1)
					}
				} else {
					errWrite := ioutil.WriteFile(path+"-"+anchorPost.Id+".pending.json", pendingJSON, 0644)
					if errWrite != nil {
						errLogger.Printf("ERROR :%v\n", errWrite)
						if exitOnError {
							os.Exit(1)
						}
					}
				}
			}
		}
	}
}
