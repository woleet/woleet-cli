package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/woleet/woleet-cli/internal/app"
)

// anchorCmd represents the anchor command
var anchorCmd = &cobra.Command{
	Use:   "anchor",
	Short: "Recursively anchor all files in a given directory and retrieve timestamped proofs of existence",
	Long: `Recursively anchor all files in a given directory and retrieve timestamped proofs of existence
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

		runParameters.Include = checkInclude(cmd)
		runParameters.Token = checkToken(cmd)

		runParameters.BaseURL = viper.GetString("api.url")
		runParameters.InvertPrivate = !viper.GetBool("api.private")

		runParameters.Strict = viper.GetBool("app.strict")
		runParameters.Prune = viper.GetBool("app.prune")
		runParameters.ExitOnError = viper.GetBool("app.exitOnError")
		runParameters.Recursive = viper.GetBool("app.recursive")

		if viper.GetBool("app.dryRun") {
			os.Exit(app.DryRun(runParameters, log))
		}
		os.Exit(app.BulkAnchor(runParameters, log))
	},
}

func init() {
	rootCmd.AddCommand(anchorCmd)

	anchorCmd.Flags().StringVarP(&directory, "directory", "d", "", "source directory containing files to anchor (required)")
	anchorCmd.Flags().StringVarP(&include, "include", "i", "", "only files taht match that regex will be anchored")
	anchorCmd.Flags().StringVarP(&s3Bucket, "s3Bucket", "", "", "bucket name that contains files to anchor")
	anchorCmd.Flags().StringVarP(&s3Endpoint, "s3Endpoint", "", "s3.amazonaws.com", `Specify an alternative S3 endpoint: ex: storage.googleapis.com,
don't specify the transport (https://), https will be used by default if you want to use http see --s3NoSSL param`)
	anchorCmd.Flags().StringVarP(&s3AccessKeyID, "s3AccessKeyID", "", "", "your AccessKeyID")
	anchorCmd.Flags().StringVarP(&s3SecretAccessKey, "s3SecretAccessKey", "", "", "your SecretAccessKey")
	anchorCmd.Flags().BoolVarP(&strict, "strict", "", false, "re-anchor any file that has changed since last anchoring")
	anchorCmd.Flags().BoolVarP(&prune, "prune", "", false, `delete receipts that are not along the original file,
with --strict it checks the hash of the original file and deletes the receipt if they do not match`)
	anchorCmd.Flags().BoolVarP(&exitOnError, "exitOnError", "e", false, "exit with an error code if anything goes wrong")
	anchorCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "explore sub-folders recursively")
	anchorCmd.Flags().BoolVarP(&private, "private", "p", false, "create non discoverable proofs")
	anchorCmd.Flags().BoolVarP(&dryRun, "dryRun", "", false, "print information about files to anchor without anchoring")
	anchorCmd.Flags().BoolVarP(&s3NoSSL, "s3NoSSL", "", false, "Use S3 without SSL (Strongly discouraged)")

	viper.BindPFlag("app.directory", anchorCmd.Flags().Lookup("directory"))
	viper.BindPFlag("app.include", anchorCmd.Flags().Lookup("include"))
	viper.BindPFlag("s3.bucket", anchorCmd.Flags().Lookup("s3Bucket"))
	viper.BindPFlag("s3.endpoint", anchorCmd.Flags().Lookup("s3Endpoint"))
	viper.BindPFlag("s3.accessKeyID", anchorCmd.Flags().Lookup("s3AccessKeyID"))
	viper.BindPFlag("s3.secretAccessKey", anchorCmd.Flags().Lookup("s3SecretAccessKey"))
	viper.BindPFlag("app.strict", anchorCmd.Flags().Lookup("strict"))
	viper.BindPFlag("app.prune", anchorCmd.Flags().Lookup("prune"))
	viper.BindPFlag("app.exitOnError", anchorCmd.Flags().Lookup("exitOnError"))
	viper.BindPFlag("app.recursive", anchorCmd.Flags().Lookup("recursive"))
	viper.BindPFlag("api.private", anchorCmd.Flags().Lookup("private"))
	viper.BindPFlag("app.dryRun", anchorCmd.Flags().Lookup("dryRun"))
	viper.BindPFlag("s3.noSSL", anchorCmd.Flags().Lookup("s3NoSSL"))

	viper.BindEnv("app.directory")
	viper.BindEnv("app.include")
	viper.BindEnv("s3.bucket")
	viper.BindEnv("s3.endpoint")
	viper.BindEnv("s3.accessKeyID")
	viper.BindEnv("s3.secretAccessKey")
	viper.BindEnv("app.strict")
	viper.BindEnv("app.prune")
	viper.BindEnv("app.exitOnError")
	viper.BindEnv("app.recursive")
	viper.BindEnv("api.private")
	viper.BindEnv("app.dryRun")
	viper.BindEnv("s3.noSSL")
}
