package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const SuffixAnchorPending string = ".pending.json"
const SuffixAnchorReceipt string = ".receipt.json"
const SuffixSignaturePending string = ".signature-pending.json"
const SuffixSignatureReceipt string = ".signature-receipt.json"

type RegexExtracted struct {
	Filename string
	AnchorID string
	Suffix   string
}

func checkFilename(fileInfo os.FileInfo) bool {
	if !fileInfo.Mode().IsRegular() {
		return false
	} else if strings.HasPrefix(fileInfo.Name(), ".") {
		return false
	} else {
		return true
	}
}

func ExploreDirectory(directory string) (map[string]os.FileInfo, error) {
	mapPathFileinfo := make(map[string]os.FileInfo)
	errWalk := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			pathlenght := len(strings.Replace(strings.TrimPrefix(path, directory), string(os.PathSeparator), "", -1))
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			} else if pathlenght > 128 {
				fmt.Fprintf(os.Stderr, "The directory: %s will be ignored, as it's path exceed 128 chars\n", path)
				return filepath.SkipDir
			} else if strings.Contains(strings.TrimPrefix(path, directory), " ") {
				fmt.Fprintf(os.Stderr, "The directory: %s will be ignored, as it's name contains a space\n", path)
				return filepath.SkipDir
			}
		}
		if checkFilename(info) {
			mapPathFileinfo[path] = info
		}
		return nil
	})
	return mapPathFileinfo, errWalk
}

func Separate(mapPathFileinfo map[string]os.FileInfo, signature bool, strict bool) (map[string]os.FileInfo, map[string]os.FileInfo, map[string]os.FileInfo) {
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
	return mapPathFileinfo, pendingFiles, receiptedFiles
}

func GetAnchorIDFromName(fileInfo os.FileInfo) (*RegexExtracted, error) {
	re := regexp.MustCompile(".*?-(?P<anchor_id>[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12})(" + strings.Replace(SuffixAnchorPending+"|"+SuffixAnchorReceipt+"|"+SuffixSignaturePending+"|"+SuffixSignatureReceipt, ".", "\\.", -1) + ")")
	match := re.FindStringSubmatch(fileInfo.Name())
	if len(match) != 3 {
		err := errors.New("Unable to extract anchorID form the filename:" + fileInfo.Name())
		return nil, err
	}
	out := new(RegexExtracted)
	out.Filename = match[0]
	out.AnchorID = match[1]
	out.Suffix = match[2]
	return out, nil
}
