package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func checkToken(cmd *cobra.Command) string {
	if !viper.IsSet("api.token") || strings.EqualFold(viper.GetString("api.token"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a token")
	}
	return viper.GetString("api.token")
}

func checkExportDirectory(cmd *cobra.Command) string {
	viper.Set("app.directory", viper.GetString("export.directory"))
	return checkDirectory(cmd)
}

func checkDirectory(cmd *cobra.Command) string {
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
	return absDirectory
}

func checkWidSignURL(cmd *cobra.Command) string {
	if !viper.IsSet("sign.widsSignURL") || strings.EqualFold(viper.GetString("sign.widsSignURL"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a widsSignURL")
	}
	return viper.GetString("sign.widsSignURL")
}

func checkWidToken(cmd *cobra.Command) string {
	if !viper.IsSet("sign.widsToken") || strings.EqualFold(viper.GetString("sign.widsToken"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a widsToken")
	}
	return viper.GetString("sign.widsToken")
}
