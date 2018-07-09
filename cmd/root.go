package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// var cfgFile string
var BaseURL string
var Token string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "woleet-cli",
	Short: "Woleet command line interface (for now only anchor is available)",
	Long: `woleet-cli is a command line interface created to interact with
woleet api available at: https://api.woleet.io for now, this tool
just support folder anchoring`,
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

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "Config", "", "config file (default is $HOME/.woleet-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&BaseURL, "url", "u", "https://api.woleet.io/v1", "Custom api url")
	rootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", "JWT token (required)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}

	// 	// Search config in home directory with name ".woleet-cli" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".woleet-cli")
	// }

	// viper.AutomaticEnv() // read in environment variables that match

	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}
