package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/woleet/woleet-cli/internal/app"
)

// timestampCmd represents the timestamp/anchor command
var timestampCmd = &cobra.Command{
	Use:     "timestamp",
	Aliases: []string{"anchor"},
	Short:   "Recursively anchor all files in a given directory and retrieve proofs of timestamp",
	Long: `Recursively anchor all files in a given directory and retrieve proofs of timestamp
Proofs being created asynchronously, you need to run the command at least twice with enough internal to retrieve the proofs.`,
	Run: func(cmd *cobra.Command, args []string) {
		runParameters := new(app.RunParameters)
		runParameters.Signature = false

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
		runParameters.RenameReceipts = viper.GetBool("app.renameReceipts")
		runParameters.ExitOnError = viper.GetBool("app.exitOnError")
		runParameters.Recursive = viper.GetBool("app.recursive")

		if viper.GetBool("app.dryRun") {
			os.Exit(app.DryRun(runParameters, log))
		}
		os.Exit(app.BulkAnchor(runParameters, log))
	},
}

func init() {
	rootCmd.AddCommand(timestampCmd)

	timestampCmd.Flags().StringVarP(&directory, "directory", "d", "", "source directory containing files to anchor (required)")
	timestampCmd.Flags().StringVarP(&filter, "filter", "f", "", "anchor only files matching this regex")
	timestampCmd.Flags().StringVar(&s3Bucket, "s3Bucket", "", "bucket name that contains files to anchor")
	timestampCmd.Flags().StringVarP(&s3Endpoint, "s3Endpoint", "", "s3.amazonaws.com", `specify an alternative S3 endpoint: ex: storage.googleapis.com,
don't specify the transport (https://), https will be used by default if you want to use http see --s3NoSSL param`)
	timestampCmd.Flags().StringVar(&s3AccessKeyID, "s3AccessKeyID", "", "your AccessKeyID")
	timestampCmd.Flags().StringVar(&s3SecretAccessKey, "s3SecretAccessKey", "", "your SecretAccessKey")
	timestampCmd.Flags().BoolVar(&strict, "strict", false, "re-anchor any file that has changed since last anchoring")
	timestampCmd.Flags().BoolVar(&prune, "prune", false, `delete receipts that are not along the original file,
with --strict it checks the hash of the original file and deletes the receipt if they do not match`)
	timestampCmd.Flags().BoolVar(&fixReceipts, "fixReceipts", false, "Check the format and fix (if necessary) every existing receipts")
	timestampCmd.Flags().BoolVar(&renameReceipts, "renameReceipts", false, "Rename legacy receipts ending by anchor-receipt.json to timestamp-receipt.json")
	timestampCmd.Flags().BoolVarP(&exitOnError, "exitOnError", "e", false, "exit with an error code if anything goes wrong")
	timestampCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "explore sub-folders recursively")
	timestampCmd.Flags().BoolVarP(&private, "private", "p", false, "create non discoverable proofs")
	timestampCmd.Flags().BoolVar(&dryRun, "dryRun", false, "print information about files to anchor without anchoring")
	timestampCmd.Flags().BoolVar(&s3NoSSL, "s3NoSSL", false, "use S3 without SSL (strongly discouraged)")

	viper.BindPFlag("app.directory", timestampCmd.Flags().Lookup("directory"))
	viper.BindPFlag("app.filter", timestampCmd.Flags().Lookup("filter"))
	viper.BindPFlag("s3.bucket", timestampCmd.Flags().Lookup("s3Bucket"))
	viper.BindPFlag("s3.endpoint", timestampCmd.Flags().Lookup("s3Endpoint"))
	viper.BindPFlag("s3.accessKeyID", timestampCmd.Flags().Lookup("s3AccessKeyID"))
	viper.BindPFlag("s3.secretAccessKey", timestampCmd.Flags().Lookup("s3SecretAccessKey"))
	viper.BindPFlag("app.strict", timestampCmd.Flags().Lookup("strict"))
	viper.BindPFlag("app.prune", timestampCmd.Flags().Lookup("prune"))
	viper.BindPFlag("app.fixReceipts", timestampCmd.Flags().Lookup("fixReceipts"))
	viper.BindPFlag("app.renameReceipts", timestampCmd.Flags().Lookup("renameReceipts"))
	viper.BindPFlag("app.exitOnError", timestampCmd.Flags().Lookup("exitOnError"))
	viper.BindPFlag("app.recursive", timestampCmd.Flags().Lookup("recursive"))
	viper.BindPFlag("api.private", timestampCmd.Flags().Lookup("private"))
	viper.BindPFlag("app.dryRun", timestampCmd.Flags().Lookup("dryRun"))
	viper.BindPFlag("s3.noSSL", timestampCmd.Flags().Lookup("s3NoSSL"))
}
