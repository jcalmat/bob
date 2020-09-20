package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Commands []Command
var Templates map[string]Template

type Command struct {
	Alias     string   `yaml:"alias"`
	Templates []string `yaml:"templates"`
}

type Template struct {
	Path      string
	Variables []string
}

func parseConfig() {
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	Commands = make([]Command, 0)
	Templates = make(map[string]Template)

	vcs := viper.GetStringMap("commands")
	for k := range vcs {
		Commands = append(Commands, Command{
			Alias:     viper.GetString(fmt.Sprintf("commands.%s.alias", k)),
			Templates: viper.GetStringSlice(fmt.Sprintf("commands.%s.templates", k)),
		})
	}

	vts := viper.GetStringMap("templates")
	for k := range vts {
		Templates[k] = Template{
			Path:      viper.GetString(fmt.Sprintf("templates.%s.path", k)),
			Variables: viper.GetStringSlice(fmt.Sprintf("templates.%s.variables", k)),
		}
	}

	//TODO: Find a way to move this somewhere else
	for _, c := range Commands {
		buildCmd.AddCommand(&cobra.Command{
			Use: c.Alias,
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, s := range c.Templates {
					t, ok := Templates[s]
					if !ok {
						fmt.Printf("Template %s not found, skipping", s)
						continue
					}

					PrintTitle("Info")

					fmt.Println("Using template path: ", t.Path)

					path := Ask("Where do you want to build your code? ")

					variables := make(map[string]string)

					PrintTitle("Variable replacement")

					for _, v := range t.Variables {
						variables[v] = Ask(fmt.Sprintf("%s: ", v))
					}

					fmt.Printf("%v\n", variables)
					cmd := exec.Command("ls", path)
					err := cmd.Run()
					if err != nil {
						return err
					}
				}
				return nil
			},
		})
	}
}

func Ask(question string) string {
	fmt.Print(question)
	var answer string
	fmt.Scanln(&answer)
	return answer
}

func PrintTitle(title string) {
	fmt.Printf("==== %s ====\n", title)
}
