package helpers

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
)

const SuffixAnchorPendingCurrent string = ".timestamp-pending.json"
const SuffixAnchorPendingLegacy string = ".anchor-pending.json"
const SuffixAnchorReceiptCurrent string = ".timestamp-receipt.json"
const SuffixAnchorReceiptLegacy string = ".anchor-receipt.json"
const SuffixSignaturePendingCurrent string = ".seal-pending.json"
const SuffixSignaturePendingLegacy string = ".signature-pending.json"
const SuffixSignatureReceiptCurrent string = ".seal-receipt.json"
const SuffixSignatureReceiptLegacy string = ".signature-receipt.json"

const SuffixRegexpCurrent = SuffixAnchorPendingCurrent + "|" + SuffixAnchorReceiptCurrent + "|" + SuffixSignaturePendingCurrent + "|" + SuffixSignatureReceiptCurrent
const SuffixRegexpLegacy = SuffixAnchorPendingLegacy + "|" + SuffixAnchorReceiptLegacy + "|" + SuffixSignaturePendingLegacy + "|" + SuffixSignatureReceiptLegacy
const SuffixRegexpAll = SuffixRegexpCurrent + "|" + SuffixRegexpLegacy

var regexHasSuffixAnchorPending = regexp.MustCompile("(" + SuffixAnchorPendingCurrent + "|" + SuffixAnchorPendingLegacy + ")$")
var regexHasSuffixAnchorReceipt = regexp.MustCompile("(" + SuffixAnchorReceiptCurrent + "|" + SuffixAnchorReceiptLegacy + ")$")
var regexHasSuffixSignaturePending = regexp.MustCompile("(" + SuffixSignaturePendingCurrent + "|" + SuffixSignaturePendingLegacy + ")$")
var regexHasSuffixSignatureReceipt = regexp.MustCompile("(" + SuffixSignatureReceiptCurrent + "|" + SuffixSignatureReceiptLegacy + ")$")

var regexpAnchorIDFromName = regexp.MustCompile("(^.*)-(?P<anchor_id>[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12})(" + strings.Replace(SuffixRegexpAll, ".", "\\.", -1) + ")$")
var regexNameSuffixReceipt = regexp.MustCompile("^.*" + "(" + SuffixRegexpAll + ")$")

type RegexExtracted struct {
	Filename         string
	OriginalFilename string
	AnchorID         string
	Suffix           string
}

func checkFilename(fileName string, filter *regexp.Regexp) bool {
	if strings.HasPrefix(fileName, ".") {
		return false
	}
	if filter != nil {
		tempFileName := fileName
		if regexNameSuffixReceipt.MatchString(fileName) {
			anchorIDF, errAnchorIDF := GetAnchorIDFromName(fileName)
			if errAnchorIDF != nil {
				// Safeguard here, not the best way to handle issues
				return false
			}
			tempFileName = anchorIDF.OriginalFilename
		}
		if !filter.MatchString(tempFileName) {
			return false
		}
	}
	return true
}

func checkDirectory(path string, directory string) bool {
	if strings.HasPrefix(directory, ".") {
		return false
	}
	return true
}

func checkDirectoryS3(path string, log *logrus.Logger) bool {
	if !strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, extractFileNameFromPathS3(path))
	}

	dirArray := strings.Split(strings.TrimSuffix(path, "/"), "/")
	for _, dir := range dirArray {
		if strings.HasPrefix(dir, ".") {
			return false
		}
	}

	return true
}

func ExploreDirectory(baseDirectory string, recursive bool, filter *regexp.Regexp, log *logrus.Logger) (map[string]string, error) {
	mapPathFileName := make(map[string]string)
	errWalk := filepath.Walk(baseDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if strings.EqualFold(filepath.Clean(baseDirectory), filepath.Clean(path)) {
			} else if !recursive {
				return filepath.SkipDir
			} else if !checkDirectory(path, info.Name()) {
				return filepath.SkipDir
			}
		} else if info.Mode().IsRegular() && checkFilename(info.Name(), filter) {
			mapPathFileName[path] = info.Name()
		}
		return nil
	})
	return mapPathFileName, errWalk
}

func ExploreS3(S3Client *minio.Client, bucket string, recursive bool, filter *regexp.Regexp, log *logrus.Logger) map[string]string {
	mapPathFileName := make(map[string]string)
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := S3Client.ListObjects(bucket, "", true, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			log.Warnln(object.Err)
		} else {
			if !strings.Contains(object.Key, "/") && checkFilename(object.Key, filter) {
				mapPathFileName[object.Key] = object.Key
			} else if strings.HasSuffix(object.Key, "/") {
				checkDirectoryS3(object.Key, log)
			} else if !strings.HasSuffix(object.Key, "/") && checkFilename(extractFileNameFromPathS3(object.Key), filter) && checkDirectoryS3(object.Key, log) {
				mapPathFileName[object.Key] = extractFileNameFromPathS3(object.Key)
			}
		}
	}
	return mapPathFileName
}

func SeparateAll(mapPathFileinfo map[string]string) (map[string]string, map[string]string, map[string]string, map[string]string) {
	anchorPendingFiles := make(map[string]string)
	anchorReceiptedFiles := make(map[string]string)
	signaturePendingFiles := make(map[string]string)
	signatureReceiptedFiles := make(map[string]string)
	for path, fileName := range mapPathFileinfo {
		if regexHasSuffixAnchorPending.MatchString(fileName) {
			anchorPendingFiles[path] = fileName
			delete(mapPathFileinfo, path)
		} else if regexHasSuffixAnchorReceipt.MatchString(fileName) {
			anchorReceiptedFiles[path] = fileName
			delete(mapPathFileinfo, path)
		} else if regexHasSuffixSignaturePending.MatchString(fileName) {
			signaturePendingFiles[path] = fileName
			delete(mapPathFileinfo, path)
		} else if regexHasSuffixSignatureReceipt.MatchString(fileName) {
			signatureReceiptedFiles[path] = fileName
			delete(mapPathFileinfo, path)
		}
	}
	return anchorPendingFiles, anchorReceiptedFiles, signaturePendingFiles, signatureReceiptedFiles
}

func GetAnchorIDFromName(fileName string) (*RegexExtracted, error) {
	match := regexpAnchorIDFromName.FindStringSubmatch(fileName)
	if len(match) != 4 {
		err := errors.New("Unable to extract anchorID form the filename:" + fileName)
		return nil, err
	}
	out := new(RegexExtracted)
	out.Filename = match[0]
	out.OriginalFilename = match[1]
	out.AnchorID = match[2]
	out.Suffix = match[3]
	return out, nil
}

func extractFileNameFromPathS3(path string) string {
	pathArray := strings.Split(path, "/")
	return pathArray[len(pathArray)-1]
}

func RenameLegacyReceipts(mapPathFileinfo map[string]string, signature bool, dryRun bool, log *logrus.Logger) {
	for path, fileName := range mapPathFileinfo {
		newPath := ""
		newFileName := ""
		if !signature && strings.HasSuffix(fileName, SuffixAnchorPendingLegacy) {
			newPath = strings.TrimSuffix(path, SuffixAnchorPendingLegacy) + SuffixAnchorPendingCurrent
			newFileName = strings.TrimSuffix(fileName, SuffixAnchorPendingLegacy) + SuffixAnchorPendingCurrent
		} else if !signature && strings.HasSuffix(fileName, SuffixAnchorReceiptLegacy) {
			newPath = strings.TrimSuffix(path, SuffixAnchorReceiptLegacy) + SuffixAnchorReceiptCurrent
			newFileName = strings.TrimSuffix(fileName, SuffixAnchorReceiptLegacy) + SuffixAnchorReceiptCurrent
		} else if signature && strings.HasSuffix(fileName, SuffixSignaturePendingLegacy) {
			newPath = strings.TrimSuffix(path, SuffixSignaturePendingLegacy) + SuffixSignaturePendingCurrent
			newFileName = strings.TrimSuffix(fileName, SuffixSignaturePendingLegacy) + SuffixSignaturePendingCurrent
		} else if signature && strings.HasSuffix(fileName, SuffixSignatureReceiptLegacy) {
			newPath = strings.TrimSuffix(path, SuffixSignatureReceiptLegacy) + SuffixSignatureReceiptCurrent
			newFileName = strings.TrimSuffix(fileName, SuffixSignatureReceiptLegacy) + SuffixSignatureReceiptCurrent
		}

		if newPath != "" && newFileName != "" {
			if !dryRun {
				delete(mapPathFileinfo, path)
				err := os.Rename(path, newPath)
				if err != nil {
					log.WithFields(logrus.Fields{
						"OriginalFile": path,
						"NewFile":      newPath,
						"Error":        err,
					}).Warnln("Unable to rename the file")
				}
				mapPathFileinfo[newPath] = newFileName
			} else {
				log.WithFields(logrus.Fields{
					"OriginalFile": path,
					"NewFile":      newPath,
				}).Infoln("Faking file renaming (dryRun mode)")
			}
		}
	}
}
