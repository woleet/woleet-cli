package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

const pageSize int = 1000

func errHandlerExitOnError(err error, exitOnError bool) {
	if err != nil {
		log.Errorf("%s\n", err)
		if exitOnError {
			os.Exit(1)
		}
	}
}
