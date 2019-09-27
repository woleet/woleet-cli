package app

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/woleet/woleet-cli/pkg/helpers"
)

func (commonInfos *commonInfos) getHash(path string) (string, error) {
	log.WithFields(logrus.Fields{
		"file": path,
	}).Infoln("Hashing file")

	var hash string
	var errHash error
	if commonInfos.runParameters.IsFS {
		openedFile, errOpenedFile := os.Open(path)
		if errOpenedFile != nil {
			openedFile.Close()
			return "", errOpenedFile
		}
		hash, errHash = helpers.HashFile(openedFile)
		openedFile.Close()
	} else if commonInfos.runParameters.IsS3 {
		openedFile, errOpenedFile := commonInfos.runParameters.S3Client.GetObject(commonInfos.runParameters.S3Bucket, path, minio.GetObjectOptions{})
		if errOpenedFile != nil {
			openedFile.Close()
			errHandlerExitOnError(errOpenedFile, commonInfos.runParameters.ExitOnError)
			return "", errOpenedFile
		}
		hash, errHash = helpers.HashFile(openedFile)
		openedFile.Close()
	}
	return hash, errHash
}

func (commonInfos *commonInfos) removeFile(path string) error {
	var errRemove error
	if commonInfos.runParameters.IsFS {
		errRemove = os.Remove(path)
	} else if commonInfos.runParameters.IsS3 {
		errRemove = commonInfos.runParameters.S3Client.RemoveObject(commonInfos.runParameters.S3Bucket, path)
	}
	return errRemove
}

func (commonInfos *commonInfos) readFile(path string) ([]byte, error) {
	var content []byte
	var contentErr error
	if commonInfos.runParameters.IsFS {
		content, contentErr = ioutil.ReadFile(path)
	} else if commonInfos.runParameters.IsS3 {
		info, infoErr := commonInfos.runParameters.S3Client.StatObject(commonInfos.runParameters.S3Bucket, path, minio.StatObjectOptions{})
		if infoErr != nil {
			return content, infoErr
		}
		reader, errReader := commonInfos.runParameters.S3Client.GetObject(commonInfos.runParameters.S3Bucket, path, minio.GetObjectOptions{})
		if errReader != nil {
			return content, errReader
		}
		content = make([]byte, info.Size)
		_, contentErr = reader.Read(content)
		reader.Close()
		if contentErr == io.EOF {
			contentErr = nil
		}
	}
	return content, contentErr
}

func (commonInfos *commonInfos) writeFile(path string, json []byte) error {
	var errWrite error
	if commonInfos.runParameters.IsFS {
		errWrite = ioutil.WriteFile(path, json, 0644)
	} else if commonInfos.runParameters.IsS3 {
		_, errWrite = commonInfos.runParameters.S3Client.PutObject(commonInfos.runParameters.S3Bucket, path, bytes.NewReader(json), -1, minio.PutObjectOptions{})
	}
	return errWrite
}
