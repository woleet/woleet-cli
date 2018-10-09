package app

import (
	"log"
	"os"
)

const pageSize int = 1000

func errHandlerExitOnError(err error, errLogger *log.Logger, exitOnError bool) {
	if err != nil {
		errLogger.Printf("ERROR: %v\n", err)
		if exitOnError {
			os.Exit(1)
		}
	}
}
