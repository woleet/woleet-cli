package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/woleet/woleet-cli/internal/app"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Download all proofs for your anchors in a given directory",
	Long: `Download all proofs for your anchors in a given directory
You can specify a date to only get proofs created after this date`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("api.token") || strings.EqualFold(viper.GetString("api.token"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a token")
		}

		if !viper.IsSet("export.directory") || strings.EqualFold(viper.GetString("export.directory"), "") {
			if !viper.GetBool("log.json") {
				cmd.Help()
			}
			log.Fatalln("Please set a directory")
		}

		absExportDirectory, errAbs := filepath.Abs(viper.GetString("export.directory"))
		if errAbs != nil {
			log.Fatalln("Unable to get Absolute directory from --directory")
		}

		info, err := os.Stat(absExportDirectory)
		if err != nil {
			log.Fatalln("The provided directory does not exists")
		} else {
			if !info.IsDir() {
				log.Fatalln("The provided path is not a directory")
			}
		}

		var unixEpochLimit int64 = 0
		if viper.IsSet("export.limitDate") && !strings.EqualFold(viper.GetString("export.limitDate"), "") {
			limitDate := viper.GetString("export.limitDate")
			limitDateArray := strings.Split(limitDate, "-")
			if len(limitDateArray) != 3 {
				log.Fatalln("The provided date is not properly formatted")
			}
			year, errYear := strconv.Atoi(limitDateArray[0])
			if errYear != nil {
				log.Fatalln("Unable to parse the provided year")
			}
			if !((year >= 0) && (year <= 9999)) {
				log.Fatalln("Please set the year between 0 and 9999")
			}

			month, errMonth := strconv.Atoi(limitDateArray[1])
			if errMonth != nil {
				log.Fatalln("Unable to parse the provided month")
			}
			if !((month >= 0) && (month <= 12)) {
				log.Fatalln("Please set the month between 0 and 12")
			}

			day, errDay := strconv.Atoi(limitDateArray[2])
			if errDay != nil {
				log.Fatalln("Unable to parse the provided day")
			}
			if !((day >= 0) && (day <= 31)) {
				log.Fatalln("Please set the day between 0 and 31")
			}
			unixEpochLimit = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).UnixNano()
		}
		app.ExportReceipts(viper.GetString("api.token"), viper.GetString("api.url"), absExportDirectory, unixEpochLimit, viper.GetBool("export.exitonerror"), log)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportDirectory, "directory", "d", "", "directory where to store the proofs (required)")
	exportCmd.Flags().StringVarP(&exportLimitDate, "limitdate", "l", "", "get only proofs created after the provided date (format: yyyy-MM-dd)")
	exportCmd.Flags().BoolVarP(&exportExitonerror, "exitonerror", "e", false, "exit with an error code if anything goes wrong")

	viper.BindPFlag("export.directory", exportCmd.Flags().Lookup("directory"))
	viper.BindPFlag("export.limitdate", exportCmd.Flags().Lookup("limitdate"))
	viper.BindPFlag("export.exitonerror", exportCmd.Flags().Lookup("exitonerror"))

	viper.BindEnv("export.directory")
	viper.BindEnv("export.limitdate")
	viper.BindEnv("export.exitonerror")
}
