package cmd

import (
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
		absDirectory := checkExportDirectory(cmd)
		token := checkToken(cmd)
		domain := checkDomain(cmd)

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
		app.ExportReceipts(token, viper.GetString("api.url"), domain, absDirectory, unixEpochLimit, viper.GetBool("export.exitonerror"), log)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportDirectory, "directory", "d", "", "directory where to store the proofs (required)")
	exportCmd.Flags().StringVarP(&exportLimitDate, "limitDate", "l", "", "get only proofs created after the provided date (format: yyyy-MM-dd)")
	exportCmd.Flags().BoolVarP(&exportExitOnError, "exitOnError", "e", false, "exit with an error code if anything goes wrong")

	viper.BindPFlag("export.directory", exportCmd.Flags().Lookup("directory"))
	viper.BindPFlag("export.limitDate", exportCmd.Flags().Lookup("limitDate"))
	viper.BindPFlag("export.exitOnError", exportCmd.Flags().Lookup("exitOnError"))

	viper.BindEnv("export.directory")
	viper.BindEnv("export.limitDate")
	viper.BindEnv("export.exitOnError")
}
