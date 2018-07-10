package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/woleet/woleet-cli/internal/app"
)

// anchorCmd represents the anchor command
var anchorCmd = &cobra.Command{
	Use:   "anchor",
	Short: "Recursively anchor all files in a directory to create timestamped proofs of existence",
	Long: `woleet-cli is a command line interface allowing to interact with
woleet API (https://api.woleet.io). For now, this tool
just support folder anchoring.`,
	Run: func(cmd *cobra.Command, args []string) {
		if strings.EqualFold(Token, "") {
			cmd.Help()
			log.Fatalln("Please set a token")
		}
		directory, errDir := cmd.Flags().GetString("directory")
		if errDir != nil {
			log.Fatalln("Unable to parse --directory flag")
		}

		if strings.EqualFold(directory, "") {
			log.Fatalln("Please set a directory")
		}

		absDirectory, errAbs := filepath.Abs(directory)
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

		exitOnErr, errExitOnErr := cmd.Flags().GetBool("exitonerror")
		if errExitOnErr != nil {
			log.Fatalln("Unable to parse --exitonerror flag")
		}
		private, privateErr := cmd.Flags().GetBool("private")
		if privateErr != nil {
			log.Fatalln("Unable to parse --private flag")
		}
		strict, errStrict := cmd.Flags().GetBool("strict")
		if errStrict != nil {
			log.Fatalln("Unable to parse --strict flag")
		}
		strictPrune, errStrictPrune := cmd.Flags().GetBool("strict-prune")
		if errStrictPrune != nil {
			log.Fatalln("Unable to parse --strict-prune flag")
		}
		if strictPrune {
			strict = true
		}
		app.BulkAnchor(BaseURL, Token, absDirectory, exitOnErr, private, strict, strictPrune)
	},
}

func init() {
	rootCmd.AddCommand(anchorCmd)

	anchorCmd.Flags().StringP("directory", "d", "", "source directory containing files to anchor")
	anchorCmd.Flags().BoolP("strict", "", false, "re-anchor any file that has changed since last anchoring")
	anchorCmd.Flags().BoolP("strict-prune", "", false, "same as --strict, plus delete the previous anchoring receipt")
	anchorCmd.Flags().BoolP("exitonerror", "e", false, "exit the app with an error code if something goes wrong")
	anchorCmd.Flags().BoolP("private", "p", false, "create anchors with non-public access")
}
