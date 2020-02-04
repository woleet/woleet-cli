package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/woleet/woleet-cli/internal/app"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Recursively sign all files in a given directory and retrieve timestamped proofs of signature",
	Long: `Recursively sign all files in a given directory and retrieve timestamped proofs of signature
Proofs being created asynchronously, you need to run the command at least twice with enough internal to retrieve the proofs.`,
	Run: func(cmd *cobra.Command, args []string) {
		runParameters := new(app.RunParameters)
		runParameters.Signature = true

		checkFolderType(cmd, runParameters)
		if runParameters.IsFS {
			runParameters.Directory = checkDirectory(cmd)
		}
		if runParameters.IsS3 {
			runParameters.S3Client = checkS3(cmd)
			runParameters.S3Bucket = viper.GetString("s3.bucket")
		}

		runParameters.Filter = checkFilter(cmd)
		runParameters.Token = checkToken(cmd)

		runParameters.BaseURL = viper.GetString("api.url")
		runParameters.InvertPrivate = !viper.GetBool("api.private")

		runParameters.Strict = viper.GetBool("app.strict")
		runParameters.Prune = viper.GetBool("app.prune")
		runParameters.FixReceipts = viper.GetBool("app.fixReceipts")
		runParameters.ExitOnError = viper.GetBool("app.exitOnError")
		runParameters.Recursive = viper.GetBool("app.recursive")

		if viper.GetBool("app.dryRun") {
			os.Exit(app.DryRun(runParameters, log))
		}

		runParameters.IDServerSignURL = checkWidSignURL(cmd)
		runParameters.IDServerToken = checkWidToken(cmd)
		runParameters.IDServerPubKey = checkWidPubKey(cmd)
		runParameters.IDServerUnsecureSSL = viper.GetBool("sign.widsUnsecureSSL")

		os.Exit(app.BulkAnchor(runParameters, log))
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&directory, "directory", "d", "", "source directory containing files to sign (required)")
	signCmd.Flags().StringVarP(&filter, "filter", "i", "", "sign only files matching this regex")
	signCmd.Flags().StringVarP(&s3Bucket, "s3Bucket", "", "", "bucket name that contains files to sign")
	signCmd.Flags().StringVarP(&s3Endpoint, "s3Endpoint", "", "s3.amazonaws.com", `specify an alternative S3 endpoint: ex: storage.googleapis.com,
	don't specify the transport (https://), https will be used by default if you want to use http see --s3NoSSL param`)
	signCmd.Flags().StringVarP(&s3AccessKeyID, "s3AccessKeyID", "", "", "your AccessKeyID")
	signCmd.Flags().StringVarP(&s3SecretAccessKey, "s3SecretAccessKey", "", "", "your SecretAccessKey")
	signCmd.Flags().StringVarP(&widsSignURL, "widsSignURL", "", "", "Woleet.ID Server sign URL ex: \"https://idserver.com:3002\" (required)")
	signCmd.Flags().StringVarP(&widsToken, "widsToken", "", "", "Woleet.ID Server API token (required)")
	signCmd.Flags().StringVarP(&widsPubKey, "widsPubKey", "", "", "public key (ie. bitcoin address) to use to sign")
	signCmd.Flags().BoolVarP(&strict, "strict", "", false, "re-sign any file that has changed since last signature or if the pubkey was changed")
	signCmd.Flags().BoolVarP(&prune, "prune", "", false, `delete receipts that are not along the original file,
with --strict it checks the hash of the original file and deletes the receipt if they do not match or if the pubkey has changed`)
	signCmd.Flags().BoolVarP(&fixReceipts, "fixReceipts", "", false, "Check the format and fix (if necessary) every existing receipts")
	signCmd.Flags().BoolVarP(&exitOnError, "exitOnError", "e", false, "exit with an error code if anything goes wrong")
	signCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "explore sub-folders recursively")
	signCmd.Flags().BoolVarP(&dryRun, "dryRun", "", false, "print information about files to sign without signing")
	signCmd.Flags().BoolVarP(&private, "private", "p", false, "create non discoverable proofs")
	signCmd.Flags().BoolVarP(&widsUnsecureSSL, "widsUnsecureSSL", "", false, "do not check Woleet.ID Server's SSL certificate validity (only for development)")
	signCmd.Flags().BoolVarP(&s3NoSSL, "s3NoSSL", "", false, "use S3 without SSL (strongly discouraged)")

	viper.BindPFlag("app.directory", signCmd.Flags().Lookup("directory"))
	viper.BindPFlag("app.filter", signCmd.Flags().Lookup("filter"))
	viper.BindPFlag("s3.bucket", signCmd.Flags().Lookup("s3Bucket"))
	viper.BindPFlag("s3.endpoint", signCmd.Flags().Lookup("s3Endpoint"))
	viper.BindPFlag("s3.accessKeyID", signCmd.Flags().Lookup("s3AccessKeyID"))
	viper.BindPFlag("s3.secretAccessKey", signCmd.Flags().Lookup("s3SecretAccessKey"))
	viper.BindPFlag("sign.widsSignURL", signCmd.Flags().Lookup("widsSignURL"))
	viper.BindPFlag("sign.widsToken", signCmd.Flags().Lookup("widsToken"))
	viper.BindPFlag("sign.widsPubKey", signCmd.Flags().Lookup("widsPubKey"))
	viper.BindPFlag("app.strict", signCmd.Flags().Lookup("strict"))
	viper.BindPFlag("app.prune", signCmd.Flags().Lookup("prune"))
	viper.BindPFlag("app.fixReceipts", signCmd.Flags().Lookup("fixReceipts"))
	viper.BindPFlag("app.exitOnError", signCmd.Flags().Lookup("exitOnError"))
	viper.BindPFlag("app.recursive", signCmd.Flags().Lookup("recursive"))
	viper.BindPFlag("app.dryRun", signCmd.Flags().Lookup("dryRun"))
	viper.BindPFlag("api.private", signCmd.Flags().Lookup("private"))
	viper.BindPFlag("sign.widsUnsecureSSL", signCmd.Flags().Lookup("widsUnsecureSSL"))
	viper.BindPFlag("s3.noSSL", signCmd.Flags().Lookup("s3NoSSL"))

	viper.BindEnv("app.directory")
	viper.BindEnv("app.filter")
	viper.BindEnv("s3.bucket")
	viper.BindEnv("s3.endpoint")
	viper.BindEnv("s3.accessKeyID")
	viper.BindEnv("s3.secretAccessKey")
	viper.BindEnv("sign.widsSignURL")
	viper.BindEnv("sign.widsToken")
	viper.BindEnv("sign.widsPubKey")
	viper.BindEnv("app.strict")
	viper.BindEnv("app.prune")
	viper.BindEnv("app.fixReceipts")
	viper.BindEnv("app.exitOnError")
	viper.BindEnv("app.recursive")
	viper.BindEnv("app.dryRun")
	viper.BindEnv("api.private")
	viper.BindEnv("sign.widsUnsecureSSL")
	viper.BindEnv("s3.noSSL")
}
