package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

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
		c := c
		buildCmd.AddCommand(&cobra.Command{
			Use: c.Alias,
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, s := range c.Templates {
					t, ok := Templates[s]
					if !ok {
						fmt.Printf("Template %s not found, skipping", s)
						continue
					}

					PrintTitle(c.Alias)

					fmt.Printf("Using template path: %s\n\n", t.Path)

					path := Ask("Where do you want to build your code (path)? ")

					replacementMap := make(map[string]string)

					PrintTitle("Variable replacement")

					for _, v := range t.Variables {
						replacementMap[v] = Ask(fmt.Sprintf("%s: ", v))
					}

					// create a tmp dir to revert the operation if an error occures
					dir, err := ioutil.TempDir("", "bob")
					if err != nil {
						return err
					}
					defer os.RemoveAll(dir) // clean up

					//TODO: Replace by actual golang code
					cmd := exec.Command("cp", "-R", t.Path, dir)
					err = cmd.Run()
					if err != nil {
						return err
					}

					files, err := ioutil.ReadDir(dir)
					if err != nil {
						return err
					}

					for k, v := range replacementMap {
						var parseFiles func(string, []os.FileInfo) error
						parseFiles = func(path string, files []os.FileInfo) error {
							for _, file := range files {
								fileName := file.Name()
								re := regexp.MustCompile(k)
								if re.MatchString(filepath.Join(path, fileName)) {
									replacedName := re.ReplaceAllString(fileName, v)
									err := os.Rename(filepath.Join(path, fileName), filepath.Join(path, replacedName))
									if err != nil {
										return err
									}
									fileName = replacedName
								}
								if file.IsDir() {
									filestmp, err := ioutil.ReadDir(filepath.Join(path, fileName))
									if err != nil {
										return err
									}
									err = parseFiles(filepath.Join(path, fileName), filestmp)
									if err != nil {
										return err
									}
								}
							}
							return nil
						}

						err = parseFiles(dir, files)
						if err != nil {
							return err
						}
						// ret, err := Exec("find", ". -type f -print0 | xargs -0 sed -i 's/", k, "/", v, "/g'")
						// if err != nil {
						// 	return err
						// }
						// fmt.Println(ret)

						// os.Rename()
						// ret, err := Exec(rename -n 's/hello/hi/g' $(find /home/devel/stuff/static/ -type f))
					}

					absPath, err := filepath.Abs(path)
					if err != nil {
						return err
					}
					err = os.Rename(filepath.Join(dir, files[0].Name()), filepath.Join(absPath, filepath.Base(files[0].Name())))
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
	fmt.Printf("\n\n==== %s ====\n\n", title)
}

func Exec(name string, args ...string) (string, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	fmt.Print(stdout.String())
	return stdout.String(), nil
}
