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
	Short: "Pass a directory and a token and start anchoring!",
	Long: `woleet-cli is a command line interface created to interact with
woleet api available at: https://api.woleet.io for now, this tool
just support folder anchoring`,
	Run: func(cmd *cobra.Command, args []string) {
		if strings.EqualFold(Token, "") {
			cmd.Help()
			log.Fatalln("Please set a token")
		}
		directory, errDir := cmd.Flags().GetString("directory")
		if errDir != nil {
			log.Fatalln("Unable to parse --directory flag")
		}
		absDirectory, errAbs := filepath.Abs(directory)
		if errAbs != nil {
			log.Fatalln("Unable to get Absolute directory from --directory")
		}
		exitOnErr, errExitOnErr := cmd.Flags().GetBool("exitonerror")
		if errExitOnErr != nil {
			log.Fatalln("Unable to parse --exitonerror flag")
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
		app.BulkAnchor(BaseURL, Token, absDirectory, exitOnErr, strict, strictPrune)
	},
}

func init() {
	rootCmd.AddCommand(anchorCmd)

	pwd, errPath := os.Getwd()
	if errPath != nil {
		log.Fatalln("Unable to get the path of the current directory")
	}

	anchorCmd.Flags().StringP("directory", "d", pwd, "Source directory to read from")
	anchorCmd.Flags().BoolP("strict", "", false, "Rescan the directory for file changes")
	anchorCmd.Flags().BoolP("strict-prune", "", false, "Rescan the directory for file changes and delete old receipts and pending file that does not have the same hash")
	anchorCmd.Flags().BoolP("exitonerror", "e", false, "Exit the app with an error code if something goes wrong")
}
