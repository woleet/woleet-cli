package cmd

import (
	"github.com/sirupsen/logrus"
)

// Global
var cfgFile string
var baseURL string
var token string
var logLevel string
var json bool
var log = logrus.New()

// Anchor / Sign
var directory string
var strict bool
var prune bool
var exitOnError bool
var private bool
var recursive bool
var dryRun bool

// Sign
var iDServerSignURL string
var iDServerToken string
var iDServerPubKey string
var iDServerUnsecureSSL bool

// Export
var exportDirectory string
var exportLimitDate string
var exportExitOnError bool
