package cmd

import (
	"os"
	"path/filepath"
	"strings"

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

		if !viper.IsSet("app.directory") || strings.EqualFold(viper.GetString("app.directory"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a directory")
		}

		absDirectory, errAbs := filepath.Abs(viper.GetString("app.directory"))

		if errAbs != nil {
			log.Fatalln("Unable to get Absolute directory from --directory")
		}

		info, err := os.Stat(absDirectory)
		if err != nil {
			log.Fatalln("The provided directory does not exists")
		} else {
			if !info.IsDir() {
				log.Fatalln("The provided path is not a directory")
			}
		}
		runParameters.Directory = absDirectory

		runParameters.BaseURL = viper.GetString("api.url")
		runParameters.InvertPrivate = !viper.GetBool("api.private")

		runParameters.Prune = viper.GetBool("app.prune")
		runParameters.ExitOnError = viper.GetBool("app.exitOnError")
		runParameters.Recursive = viper.GetBool("app.recursive")
		if runParameters.Prune || viper.GetBool("app.strict") {
			runParameters.Strict = true
		} else {
			runParameters.Strict = false
		}

		if viper.GetBool("app.dryRun") {
			app.DryRun(runParameters, log)
			os.Exit(0)
		}

		if !viper.IsSet("api.token") || strings.EqualFold(viper.GetString("api.token"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a token")
		}
		runParameters.Token = viper.GetString("api.token")

		if !viper.IsSet("sign.iDServerSignURL") || strings.EqualFold(viper.GetString("sign.iDServerSignURL"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a iDServerSignURL")
		}
		runParameters.IDServerSignURL = viper.GetString("sign.iDServerSignURL")

		if !viper.IsSet("sign.iDServerToken") || strings.EqualFold(viper.GetString("sign.iDServerToken"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a iDServerToken")
		}
		runParameters.IDServerToken = viper.GetString("sign.iDServerToken")

		runParameters.IDServerUnsecureSSL = viper.GetBool("sign.iDServerUnsecureSSL")
		if viper.IsSet("sign.iDServerPubKey") {
			runParameters.IDServerPubKey = viper.GetString("sign.iDServerPubKey")
		}

		app.BulkAnchor(runParameters, log)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&directory, "directory", "d", "", "source directory containing files to sign (required)")
	signCmd.Flags().StringVarP(&iDServerSignURL, "iDServerSignURL", "", "", "Woleet.ID Server sign URL ex: \"https://idserver.com:4443/sign\" (required)")
	signCmd.Flags().StringVarP(&iDServerToken, "iDServerToken", "", "", "Woleet.ID Server API token (required)")
	signCmd.Flags().StringVarP(&iDServerPubKey, "iDServerPubKey", "", "", "public key (ie. bitcoin address) to use to sign")
	signCmd.Flags().BoolVarP(&strict, "strict", "", false, "re-sign any file that has changed since last signature")
	signCmd.Flags().BoolVarP(&prune, "prune", "", false, `delete receipts that are not along the original file,
with --strict it checks the hash of the original file and deletes the receipt if they do not match`)
	signCmd.Flags().BoolVarP(&exitOnError, "exitOnError", "e", false, "exit with an error code if anything goes wrong")
	signCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "explore sub-folders recursively")
	signCmd.Flags().BoolVarP(&dryRun, "dryRun", "", false, "print information about files to sign without signing")
	signCmd.Flags().BoolVarP(&private, "private", "p", false, "create non discoverable proofs")
	signCmd.Flags().BoolVarP(&iDServerUnsecureSSL, "iDServerUnsecureSSL", "", false, "do not check Woleet.ID Server's SSL certificate validity (only for development)")

	viper.BindPFlag("api.private", signCmd.Flags().Lookup("private"))
	viper.BindPFlag("app.strict", signCmd.Flags().Lookup("strict"))
	viper.BindPFlag("app.directory", signCmd.Flags().Lookup("directory"))
	viper.BindPFlag("app.prune", signCmd.Flags().Lookup("prune"))
	viper.BindPFlag("app.exitOnError", signCmd.Flags().Lookup("exitOnError"))
	viper.BindPFlag("app.recursive", signCmd.Flags().Lookup("recursive"))
	viper.BindPFlag("app.dryRun", signCmd.Flags().Lookup("dryRun"))
	viper.BindPFlag("sign.iDServerSignURL", signCmd.Flags().Lookup("iDServerSignURL"))
	viper.BindPFlag("sign.iDServerToken", signCmd.Flags().Lookup("iDServerToken"))
	viper.BindPFlag("sign.iDServerPubKey", signCmd.Flags().Lookup("iDServerPubKey"))
	viper.BindPFlag("sign.iDServerUnsecureSSL", signCmd.Flags().Lookup("iDServerUnsecureSSL"))

	viper.BindEnv("api.private")
	viper.BindEnv("app.directory")
	viper.BindEnv("app.strict")
	viper.BindEnv("app.strict-prune")
	viper.BindEnv("app.exitOnError")
	viper.BindEnv("app.recursive")
	viper.BindEnv("app.dryRun")
	viper.BindEnv("sign.iDServerSignURL")
	viper.BindEnv("sign.iDServerToken")
	viper.BindEnv("sign.iDServerPubKey")
	viper.BindEnv("sign.iDServerUnsecureSSL")
}
