package cmd

import (
	"github.com/sirupsen/logrus"
)

// Global
var cfgFile string
var baseURL string
var token string
var logLevel string
var jsonOut bool
var log = logrus.New()

// Anchor / Sign
var directory string
var filter string
var strict bool
var prune bool
var fixReceipts bool
var exitOnError bool
var private bool
var recursive bool
var dryRun bool

// Sign
var widsSignURL string
var widsToken string
var widsPubKey string
var widsUnsecureSSL bool

// Export
var exportDirectory string
var exportLimitDate string
var exportExitOnError bool
var exportFixReceipts bool

// S3
var s3Bucket string
var s3Endpoint string
var s3AccessKeyID string
var s3SecretAccessKey string
var s3NoSSL bool
