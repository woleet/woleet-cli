package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/woleet/woleet-cli/internal/app"
)

// sealCmd represents the sign/sign command
var sealCmd = &cobra.Command{
	Use:     "seal",
	Aliases: []string{"sign"},
	Short:   "Recursively seal all files in a given directory and retrieve proofs of seal",
	Long: `Recursively sign all files in a given directory and retrieve proofs of seal
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

		if viper.IsSet("sign.widsUnsecureSSL") {
			runParameters.IDServerUnsecureSSL = viper.GetBool("sign.widsUnsecureSSL")
		}
		if viper.IsSet("seal.widsUnsecureSSL") {
			runParameters.IDServerUnsecureSSL = viper.GetBool("seal.widsUnsecureSSL")
		}

		//TODO
		params, _ := json.MarshalIndent(runParameters, "", " ")
		fmt.Printf("%s\n", string(params))
		os.Exit(0)

		//os.Exit(app.BulkAnchor(runParameters, log))

	},
}

func init() {
	rootCmd.AddCommand(sealCmd)

	sealCmd.Flags().StringVarP(&directory, "directory", "d", "", "source directory containing files to seal (required)")
	sealCmd.Flags().StringVarP(&filter, "filter", "i", "", "seal only files matching this regex")
	sealCmd.Flags().StringVar(&s3Bucket, "s3Bucket", "", "bucket name that contains files to seal")
	sealCmd.Flags().StringVarP(&s3Endpoint, "s3Endpoint", "", "s3.amazonaws.com", `specify an alternative S3 endpoint: ex: storage.googleapis.com,
	don't specify the transport (https://), https will be used by default if you want to use http see --s3NoSSL param`)
	sealCmd.Flags().StringVar(&s3AccessKeyID, "s3AccessKeyID", "", "your AccessKeyID")
	sealCmd.Flags().StringVar(&s3SecretAccessKey, "s3SecretAccessKey", "", "your SecretAccessKey")
	sealCmd.Flags().StringVar(&widsSignURL, "widsSignURL", "", "Woleet.ID Server sign URL ex: \"https://idserver.com:3002\" (required)")
	sealCmd.Flags().StringVar(&widsToken, "widsToken", "", "Woleet.ID Server API token (required)")
	sealCmd.Flags().StringVar(&widsPubKey, "widsPubKey", "", "public key (ie. bitcoin address) to use to seal (required)")
	sealCmd.Flags().BoolVar(&strict, "strict", false, "re-seal any file that has changed since last sealing or if the pubkey was changed")
	sealCmd.Flags().BoolVar(&prune, "prune", false, `delete receipts that are not along the original file,
with --strict it checks the hash of the original file and deletes the receipt if they do not match or if the pubkey has changed`)
	sealCmd.Flags().BoolVar(&fixReceipts, "fixReceipts", false, `Check the format and fix (if necessary) every existing receipts,
 also rename legacy receipts ending by signature-receipt.json to seal-receipt.json`)
	sealCmd.Flags().BoolVarP(&exitOnError, "exitOnError", "e", false, "exit with an error code if anything goes wrong")
	sealCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "explore sub-folders recursively")
	sealCmd.Flags().BoolVar(&dryRun, "dryRun", false, "print information about files to seal without sealing")
	sealCmd.Flags().BoolVarP(&private, "private", "p", false, "create non discoverable proofs")
	sealCmd.Flags().BoolVar(&widsUnsecureSSL, "widsUnsecureSSL", false, "do not check Woleet.ID Server's SSL certificate validity (only for development)")
	sealCmd.Flags().BoolVarP(&s3NoSSL, "s3NoSSL", "", false, "use S3 without SSL (strongly discouraged)")

	viper.BindPFlag("app.directory", sealCmd.Flags().Lookup("directory"))
	viper.BindPFlag("app.filter", sealCmd.Flags().Lookup("filter"))
	viper.BindPFlag("s3.bucket", sealCmd.Flags().Lookup("s3Bucket"))
	viper.BindPFlag("s3.endpoint", sealCmd.Flags().Lookup("s3Endpoint"))
	viper.BindPFlag("s3.accessKeyID", sealCmd.Flags().Lookup("s3AccessKeyID"))
	viper.BindPFlag("s3.secretAccessKey", sealCmd.Flags().Lookup("s3SecretAccessKey"))
	viper.BindPFlag("sign.widsSignURL", sealCmd.Flags().Lookup("widsSignURL"))
	viper.BindPFlag("seal.widsSignURL", sealCmd.Flags().Lookup("widsSignURL"))
	viper.BindPFlag("sign.widsToken", sealCmd.Flags().Lookup("widsToken"))
	viper.BindPFlag("seal.widsToken", sealCmd.Flags().Lookup("widsToken"))
	viper.BindPFlag("sign.widsPubKey", sealCmd.Flags().Lookup("widsPubKey"))
	viper.BindPFlag("seal.widsPubKey", sealCmd.Flags().Lookup("widsPubKey"))
	viper.BindPFlag("app.exitOnError", sealCmd.Flags().Lookup("exitOnError"))
	viper.BindPFlag("app.recursive", sealCmd.Flags().Lookup("recursive"))
	viper.BindPFlag("api.private", sealCmd.Flags().Lookup("private"))
	viper.BindPFlag("app.strict", sealCmd.Flags().Lookup("strict"))
	viper.BindPFlag("app.prune", sealCmd.Flags().Lookup("prune"))
	viper.BindPFlag("app.fixReceipts", sealCmd.Flags().Lookup("fixReceipts"))
	viper.BindPFlag("app.dryRun", sealCmd.Flags().Lookup("dryRun"))
	viper.BindPFlag("sign.widsUnsecureSSL", sealCmd.Flags().Lookup("widsUnsecureSSL"))
	viper.BindPFlag("seal.widsUnsecureSSL", sealCmd.Flags().Lookup("widsUnsecureSSL"))
	viper.BindPFlag("s3.noSSL", sealCmd.Flags().Lookup("s3NoSSL"))
}
