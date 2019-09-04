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

const SuffixAnchorPending string = ".anchor-pending.json"
const SuffixAnchorReceipt string = ".anchor-receipt.json"
const SuffixSignaturePending string = ".signature-pending.json"
const SuffixSignatureReceipt string = ".signature-receipt.json"
const AllSuffixRegexp = SuffixAnchorPending + "|" + SuffixAnchorReceipt + "|" + SuffixSignaturePending + "|" + SuffixSignatureReceipt

var regexpAnchorIDFromName = regexp.MustCompile("(^.*)-(?P<anchor_id>[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12})(" + strings.Replace(AllSuffixRegexp, ".", "\\.", -1) + ")$")
var regexNameSuffixReceipt = regexp.MustCompile("^.*" + "(" + AllSuffixRegexp + ")$")

type RegexExtracted struct {
	Filename         string
	OriginalFilename string
	AnchorID         string
	Suffix           string
}

func checkFilename(fileName string, include *regexp.Regexp) bool {
	if strings.HasPrefix(fileName, ".") {
		return false
	}
	if include != nil {
		tempFileName := fileName
		if regexNameSuffixReceipt.MatchString(fileName) {
			anchorIDF, errAnchorIDF := GetAnchorIDFromName(fileName)
			if errAnchorIDF != nil {
				// Safeguard here, not the best way to handle issues
				return false
			}
			tempFileName = anchorIDF.OriginalFilename
		}
		if !include.MatchString(tempFileName) {
			return false
		}
	}
	return true
}

func checkDirectory(string directoryPath) bool {

}

func ExploreDirectory(directory string, recursive bool, include *regexp.Regexp, log *logrus.Logger) (map[string]string, error) {
	mapPathFileName := make(map[string]string)
	errWalk := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if !recursive && !strings.EqualFold(filepath.Clean(directory), filepath.Clean(path)) {
				return filepath.SkipDir
			}
			pathlenght := len(strings.Replace(strings.TrimPrefix(path, directory), string(os.PathSeparator), "", -1))
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			} else if pathlenght > 128 {
				log.Warnf("The directory: %s will be ignored, as it's path exceed 128 chars\n", path)
				return filepath.SkipDir
			} else if strings.Contains(strings.TrimPrefix(path, directory), " ") {
				log.Warnf("The directory: %s will be ignored, as it's name contains a space\n", path)
				return filepath.SkipDir
			}
		} else if info.Mode().IsRegular() && checkFilename(info.Name(), include) {
			mapPathFileName[path] = info.Name()
		}
		return nil
	})
	return mapPathFileName, errWalk
}

func ExploreS3(S3Client *minio.Client, bucket string, recursive bool, include *regexp.Regexp, log *logrus.Logger) map[string]string {
	mapPathFileName := make(map[string]string)
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := S3Client.ListObjects(bucket, "", true, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			log.Warnln(object.Err)
		} else {
			if !strings.Contains(object.Key, "/") && checkFilename(object.Key, include) {
				mapPathFileName[object.Key] = object.Key
			} else if !strings.HasSuffix(object.Key, "/") && checkFilename(extractFileNameFromPath(object.Key), include) {
				mapPathFileName[extractFileNameFromPath(object.Key)] = object.Key
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
	for path, fineName := range mapPathFileinfo {
		if strings.HasSuffix(fineName, SuffixAnchorPending) {
			anchorPendingFiles[path] = fineName
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fineName, SuffixAnchorReceipt) {
			anchorReceiptedFiles[path] = fineName
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fineName, SuffixSignaturePending) {
			signaturePendingFiles[path] = fineName
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fineName, SuffixSignatureReceipt) {
			signatureReceiptedFiles[path] = fineName
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

func extractFileNameFromPath(path string) string {

}
