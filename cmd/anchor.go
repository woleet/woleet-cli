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

		runParameters.Directory = checkDirectory(cmd)
		runParameters.Token = checkToken(cmd)
		runParameters.Domain = checkDomain(cmd)

		runParameters.BaseURL = viper.GetString("api.url")
		runParameters.InvertPrivate = !viper.GetBool("api.private")

		runParameters.Strict = viper.GetBool("app.strict")
		runParameters.Prune = viper.GetBool("app.prune")
		runParameters.ExitOnError = viper.GetBool("app.exitOnError")
		runParameters.Recursive = viper.GetBool("app.recursive")

		if viper.GetBool("app.dryRun") {
			app.DryRun(runParameters, log)
			os.Exit(0)
		}

		app.BulkAnchor(runParameters, log)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(anchorCmd)

	anchorCmd.Flags().StringVarP(&directory, "directory", "d", "", "source directory containing files to anchor (required)")
	anchorCmd.Flags().BoolVarP(&strict, "strict", "", false, "re-anchor any file that has changed since last anchoring")
	anchorCmd.Flags().BoolVarP(&prune, "prune", "", false, `delete receipts that are not along the original file,
with --strict it checks the hash of the original file and deletes the receipt if they do not match`)
	anchorCmd.Flags().BoolVarP(&exitOnError, "exitOnError", "e", false, "exit with an error code if anything goes wrong")
	anchorCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "explore sub-folders recursively")
	anchorCmd.Flags().BoolVarP(&private, "private", "p", false, "create non discoverable proofs")
	anchorCmd.Flags().BoolVarP(&dryRun, "dryRun", "", false, "print information about files to anchor without anchoring")

	viper.BindPFlag("api.private", anchorCmd.Flags().Lookup("private"))
	viper.BindPFlag("app.directory", anchorCmd.Flags().Lookup("directory"))
	viper.BindPFlag("app.strict", anchorCmd.Flags().Lookup("strict"))
	viper.BindPFlag("app.prune", anchorCmd.Flags().Lookup("prune"))
	viper.BindPFlag("app.exitOnError", anchorCmd.Flags().Lookup("exitOnError"))
	viper.BindPFlag("app.recursive", anchorCmd.Flags().Lookup("recursive"))
	viper.BindPFlag("app.dryRun", anchorCmd.Flags().Lookup("dryRun"))

	viper.BindEnv("api.private")
	viper.BindEnv("app.directory")
	viper.BindEnv("app.strict")
	viper.BindEnv("app.prune")
	viper.BindEnv("app.exitOnError")
	viper.BindEnv("app.recursive")
	viper.BindEnv("app.dryRun")
}
