package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "woleet-cli",
	Version: "0.1.1",
	Short:   "Woleet command line interface",
	Long:    "woleet-cli is a command line interface allowing to interact with Woleet API (https://api.woleet.io). For now, this tool only supports anchoring and signing all files of a given folder.",
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.woleet-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", "https://api.woleet.io/v1", "custom API url")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "JWT token (required)")

	viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("api.token", rootCmd.PersistentFlags().Lookup("token"))

	viper.BindEnv("api.url")
	viper.BindEnv("api.token")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if (os.Getenv("WLT_CONFIG")) != "" {
		// Use config file from env
		viper.SetConfigFile(os.Getenv("WLT_CONFIG"))
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".woleet-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".woleet-cli")
	}

	viper.SetEnvPrefix("WLT")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	viper.ReadInConfig()

	// If a config file is found, read it in.
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	log.Fatalln("Failed to use config file:", viper.ConfigFileUsed())
	// } else {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}
