package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize bob's config file",
	Long:  `initialize a .bobconfig.yml file with basic examples in your home directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := []byte(`
# Register your commands here
commands:

templates:
`)
		_, err := os.Stat(viper.ConfigFileUsed())
		if os.IsNotExist(err) {
			err := ioutil.WriteFile(viper.ConfigFileUsed(), config, 0600)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("config file already exist")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
