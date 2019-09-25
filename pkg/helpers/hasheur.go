package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func HashFile(src io.Reader) (string, error) {
	hasheur := sha256.New()

	_, errHash := io.Copy(hasheur, src)
	if errHash != nil {
		return "", errHash
	}
	return hex.EncodeToString(hasheur.Sum(nil)), nil
}
