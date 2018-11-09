package cmd

import (
	"os"
	"path/filepath"
	"strings"

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

		runParameters.Prune = viper.GetBool("app.strict-prune")
		runParameters.ExitOnError = viper.GetBool("app.exitonerror")
		runParameters.Recursive = viper.GetBool("app.recursive")
		if runParameters.Prune || viper.GetBool("app.strict") {
			runParameters.Strict = true
		} else {
			runParameters.Strict = false
		}

		if viper.GetBool("app.dryrun") {
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

		app.BulkAnchor(runParameters, log)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(anchorCmd)

	anchorCmd.Flags().StringVarP(&directory, "directory", "d", "", "source directory containing files to anchor (required)")
	anchorCmd.Flags().BoolVarP(&strict, "strict", "", false, "re-anchor any file that has changed since last anchoring")
	anchorCmd.Flags().BoolVarP(&strictPrune, "strict-prune", "", false, "same as --strict, plus delete the previous anchoring proof")
	anchorCmd.Flags().BoolVarP(&exitonerror, "exitonerror", "e", false, "exit with an error code if anything goes wrong")
	anchorCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "explore sub-folders recursively")
	anchorCmd.Flags().BoolVarP(&private, "private", "p", false, "create non discoverable proofs")
	anchorCmd.Flags().BoolVarP(&dryRun, "dryrun", "", false, "print information about files to anchor without anchoring")

	viper.BindPFlag("api.private", anchorCmd.Flags().Lookup("private"))
	viper.BindPFlag("app.directory", anchorCmd.Flags().Lookup("directory"))
	viper.BindPFlag("app.strict", anchorCmd.Flags().Lookup("strict"))
	viper.BindPFlag("app.strict-prune", anchorCmd.Flags().Lookup("strict-prune"))
	viper.BindPFlag("app.exitonerror", anchorCmd.Flags().Lookup("exitonerror"))
	viper.BindPFlag("app.recursive", anchorCmd.Flags().Lookup("recursive"))
	viper.BindPFlag("app.dryrun", anchorCmd.Flags().Lookup("dryrun"))

	viper.BindEnv("api.private")
	viper.BindEnv("app.directory")
	viper.BindEnv("app.strict")
	viper.BindEnv("app.strict-prune")
	viper.BindEnv("app.exitonerror")
	viper.BindEnv("app.recursive")
	viper.BindEnv("app.dryrun")
}
