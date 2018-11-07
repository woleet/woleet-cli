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
		if !viper.IsSet("api.token") || strings.EqualFold(viper.GetString("api.token"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a token")
		}

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

		if !viper.IsSet("sign.backendkitSignURL") || strings.EqualFold(viper.GetString("sign.backendkitSignURL"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a backendkitSignURL")
		}

		if !viper.IsSet("sign.backendkitToken") || strings.EqualFold(viper.GetString("sign.backendkitToken"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a backendkitToken")
		}

		runParameters := new(app.RunParameters)
		runParameters.Signature = true

		runParameters.BaseURL = viper.GetString("api.url")
		runParameters.Token = viper.GetString("api.token")
		runParameters.InvertPrivate = !viper.GetBool("api.private")

		runParameters.Directory = absDirectory
		runParameters.Prune = viper.GetBool("app.strict-prune")
		runParameters.ExitOnError = viper.GetBool("app.exitonerror")
		if runParameters.Prune || viper.GetBool("app.strict") {
			runParameters.Strict = true
		} else {
			runParameters.Strict = false
		}

		runParameters.BackendkitSignURL = viper.GetString("sign.backendkitSignURL")
		runParameters.BackendkitToken = viper.GetString("sign.backendkitToken")
		runParameters.UnsecureSSL = viper.GetBool("sign.unsecureSSL")
		if !viper.IsSet("sign.backendkitPubKey") {
			runParameters.BackendkitPubKey = viper.GetString("sign.backendkitPubKey")
		}

		app.BulkAnchor(runParameters, log)
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&directory, "directory", "d", "", "source directory containing files to sign (required)")
	signCmd.Flags().StringVarP(&backendkitSignURL, "backendkitSignURL", "", "", "backendkit sign url ex: \"https://backendkit.com:4443/signature\" (required)")
	signCmd.Flags().StringVarP(&backendkitToken, "backendkitToken", "", "", "backendkit token (required)")
	signCmd.Flags().StringVarP(&backendkitPubKey, "backendkitPubKey", "", "", "backendkit pubkey")
	signCmd.Flags().BoolVarP(&strict, "strict", "", false, "re-sign any file that has changed since last signature")
	signCmd.Flags().BoolVarP(&strictPrune, "strict-prune", "", false, "same as --strict, plus delete the previous signature receipt")
	signCmd.Flags().BoolVarP(&exitonerror, "exitonerror", "e", false, "exit the app with an error code if anything goes wrong")
	signCmd.Flags().BoolVarP(&private, "private", "p", false, "create signatues with non-public access")
	signCmd.Flags().BoolVarP(&unsecureSSL, "unsecureSSL", "", false, "Do not check the ssl certificate validity for the backendkit (only use in developpement)")

	viper.BindPFlag("api.private", signCmd.Flags().Lookup("private"))
	viper.BindPFlag("app.strict", signCmd.Flags().Lookup("strict"))
	viper.BindPFlag("app.directory", signCmd.Flags().Lookup("directory"))
	viper.BindPFlag("app.strict-prune", signCmd.Flags().Lookup("strict-prune"))
	viper.BindPFlag("app.exitonerror", signCmd.Flags().Lookup("exitonerror"))
	viper.BindPFlag("sign.backendkitSignURL", signCmd.Flags().Lookup("backendkitSignURL"))
	viper.BindPFlag("sign.backendkitToken", signCmd.Flags().Lookup("backendkitToken"))
	viper.BindPFlag("sign.backendkitPubKey", signCmd.Flags().Lookup("backendkitPubKey"))
	viper.BindPFlag("sign.unsecureSSL", signCmd.Flags().Lookup("unsecureSSL"))

	viper.BindEnv("api.private")
	viper.BindEnv("app.directory")
	viper.BindEnv("app.strict")
	viper.BindEnv("app.strict-prune")
	viper.BindEnv("app.exitonerror")
	viper.BindEnv("sign.backendkitSignURL")
	viper.BindEnv("sign.backendkitToken")
	viper.BindEnv("sign.backendkitPubKey")
	viper.BindEnv("sign.unsecureSSL")
}
