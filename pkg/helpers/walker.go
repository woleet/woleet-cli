package helpers

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

func checkFilename(fileInfo os.FileInfo, include *regexp.Regexp) bool {
	if !fileInfo.Mode().IsRegular() {
		return false
	}
	if strings.HasPrefix(fileInfo.Name(), ".") {
		return false
	}
	if include != nil {
		tempFileName := fileInfo.Name()
		if regexNameSuffixReceipt.MatchString(fileInfo.Name()) {
			anchorIDF, errAnchorIDF := GetAnchorIDFromName(fileInfo)
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

func ExploreDirectory(directory string, recursive bool, include *regexp.Regexp, log *logrus.Logger) (map[string]os.FileInfo, error) {
	mapPathFileinfo := make(map[string]os.FileInfo)
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
		}
		if checkFilename(info, include) {
			mapPathFileinfo[path] = info
		}
		return nil
	})
	return mapPathFileinfo, errWalk
}

func Separate(mapPathFileinfo map[string]os.FileInfo, signature bool) (map[string]os.FileInfo, map[string]os.FileInfo) {
	pendingFiles := make(map[string]os.FileInfo)
	receiptedFiles := make(map[string]os.FileInfo)
	for path, fileinfo := range mapPathFileinfo {
		if strings.HasSuffix(fileinfo.Name(), SuffixAnchorPending) {
			if !signature {
				pendingFiles[path] = fileinfo
			}
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fileinfo.Name(), SuffixAnchorReceipt) {
			if !signature {
				receiptedFiles[path] = fileinfo
			}
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fileinfo.Name(), SuffixSignaturePending) {
			if signature {
				pendingFiles[path] = fileinfo
			}
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fileinfo.Name(), SuffixSignatureReceipt) {
			if signature {
				receiptedFiles[path] = fileinfo
			}
			delete(mapPathFileinfo, path)
		}
	}
	return pendingFiles, receiptedFiles
}

func SeparateAll(mapPathFileinfo map[string]os.FileInfo) (map[string]os.FileInfo, map[string]os.FileInfo, map[string]os.FileInfo, map[string]os.FileInfo) {
	anchorPendingFiles := make(map[string]os.FileInfo)
	anchorReceiptedFiles := make(map[string]os.FileInfo)
	signaturePendingFiles := make(map[string]os.FileInfo)
	signatureReceiptedFiles := make(map[string]os.FileInfo)
	for path, fileinfo := range mapPathFileinfo {
		if strings.HasSuffix(fileinfo.Name(), SuffixAnchorPending) {
			anchorPendingFiles[path] = fileinfo
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fileinfo.Name(), SuffixAnchorReceipt) {
			anchorReceiptedFiles[path] = fileinfo
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fileinfo.Name(), SuffixSignaturePending) {
			signaturePendingFiles[path] = fileinfo
			delete(mapPathFileinfo, path)
		} else if strings.HasSuffix(fileinfo.Name(), SuffixSignatureReceipt) {
			signatureReceiptedFiles[path] = fileinfo
			delete(mapPathFileinfo, path)
		}
	}
	return anchorPendingFiles, anchorReceiptedFiles, signaturePendingFiles, signatureReceiptedFiles
}

func GetAnchorIDFromName(fileInfo os.FileInfo) (*RegexExtracted, error) {
	match := regexpAnchorIDFromName.FindStringSubmatch(fileInfo.Name())
	if len(match) != 4 {
		err := errors.New("Unable to extract anchorID form the filename:" + fileInfo.Name())
		return nil, err
	}
	out := new(RegexExtracted)
	out.Filename = match[0]
	out.OriginalFilename = match[1]
	out.AnchorID = match[2]
	out.Suffix = match[3]
	return out, nil
}
