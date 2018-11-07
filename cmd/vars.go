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
var strictPrune bool
var exitonerror bool
var private bool

// Sign
var backendkitSignURL string
var backendkitToken string
var backendkitPubKey string
var unsecureSSL bool

// Export
var exportDirectory string
var exportLimitDate string
var exportExitonerror bool
