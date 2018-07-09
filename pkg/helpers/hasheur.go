package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func HashFile(file string) (string, error) {
	hasheur := sha256.New()
	openedFile, errFile := os.Open(file)
	if errFile != nil {
		return "", errFile
	}

	_, errHash := io.Copy(hasheur, openedFile)
	if errHash != nil {
		return "", errHash
	}

	openedFile.Close()
	return hex.EncodeToString(hasheur.Sum(nil)), nil
}
