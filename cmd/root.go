package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "woleet-cli",
	Version: "0.5.4",
	Short:   "Woleet command line interface",
	Long: `woleet-cli is a command line interface allowing to interact with the Woleet API (https://api.woleet.io) and the Woleet.ID Server API (https://app.swaggerhub.com/apis/Woleet/WoleetIDServerAPI).
For now, this tool only supports timestamping and sealing the files of a given folder as well as exporting all created proof receipts to a remote.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetOutput(os.Stdout)

		if strings.EqualFold(viper.GetString("log.level"), "info") {
			log.SetLevel(logrus.InfoLevel)
		} else if strings.EqualFold(viper.GetString("log.level"), "warn") {
			log.SetLevel(logrus.WarnLevel)
		} else if strings.EqualFold(viper.GetString("log.level"), "error") {
			log.SetLevel(logrus.ErrorLevel)
		} else if strings.EqualFold(viper.GetString("log.level"), "fatal") {
			log.SetLevel(logrus.PanicLevel)
		} else {
			fmt.Println("Unable to parse provided log level")
			os.Exit(1)
		}

		if viper.GetBool("log.json") {
			log.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})
		} else {
			log.SetFormatter(&logrus.TextFormatter{DisableLevelTruncation: true, DisableTimestamp: true})
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.woleet-cli.yaml)")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", "https://api.woleet.io/v1", "Woleet API URL")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Woleet API token (required)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "logLevel", "", "info", "select log level info|warn|error|fatal")
	rootCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "", false, "use JSON as log output format")

	viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("api.token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("logLevel"))
	viper.BindPFlag("log.json", rootCmd.PersistentFlags().Lookup("json"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if (strings.EqualFold(cfgFile, "disabled")) || (strings.EqualFold(os.Getenv("WCLI_CONFIG"), "disabled")) {
		return
	} else if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if os.Getenv("WCLI_CONFIG") != "" {
		// Use config file from env
		viper.SetConfigFile(os.Getenv("WCLI_CONFIG"))
	} else {
		// Find home directory.
		home, errHome := homedir.Dir()
		if errHome != nil {
			fmt.Println(errHome)
			os.Exit(1)
		}

		// Search config in home directory with name ".woleet-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".woleet-cli")
	}

	viper.SetEnvPrefix("WCLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	viper.ReadInConfig()
}
