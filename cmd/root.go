package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcalmat/bob/pkg/io"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bob",
	Short: "Code builder",
	Long:  `Bob is a tool for creating flexible pieces of code from templates.`,
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
	io.AsciiBob()

	updateConfig()

	// updateConfig each time the Run command is called
	cobra.OnInitialize(updateConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bobconfig)")
}

// updateConfig reads in config file and ENV variables if set.
func updateConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bobconfig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bobconfig")

		viper.SetConfigFile(filepath.Join(home, ".bobconfig.yml"))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		err := parseConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
