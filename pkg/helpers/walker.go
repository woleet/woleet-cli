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

func checkDirectory(path string, directory string, pathSeparator string, log *logrus.Logger) bool {
	if strings.HasPrefix(directory, ".") {
		return false
	}
	return true
}

func checkDirectoryS3(path string, log *logrus.Logger) bool {
	isFile := false
	if !strings.HasSuffix(path, "/") {
		isFile = true
		path = strings.TrimSuffix(path, extractFileNameFromPathS3(path))
	}

	dirArray := strings.Split(strings.TrimSuffix(path, "/"), "/")
	for _, dir := range dirArray {
		if strings.HasPrefix(dir, ".") {
			return false
		}
	}

	if len(strings.Replace(path, "/", "", -1)) > 128 {
		if !isFile {
			log.Warnf("The directory: %s will be ignored, as it's path exceed 128 chars\n", path)
		}
		return false
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
			} else if !checkDirectory(path, info.Name(), string(os.PathSeparator), log) {
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

func extractFileNameFromPathS3(path string) string {
	pathArray := strings.Split(path, "/")
	return pathArray[len(pathArray)-1]
}
