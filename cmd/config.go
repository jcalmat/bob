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

	"github.com/go-git/go-git/v5"
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
	Git       string
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
			Git:       viper.GetString(fmt.Sprintf("templates.%s.git", k)),
			Variables: viper.GetStringSlice(fmt.Sprintf("templates.%s.variables", k)),
		}
	}

	//TODO: Find a way to move this somewhere else
	buildCmd.ResetCommands()
	for _, c := range Commands {
		c := c
		buildCmd.AddCommand(&cobra.Command{
			Use: c.Alias,
			RunE: func(cmd *cobra.Command, args []string) error {

				PrintTitle(c.Alias)

				replacementMap := make(map[string]string)

				for _, s := range c.Templates {
					t, ok := Templates[s]
					if !ok {
						fmt.Printf("Template %s not found, skipping\n\n", s)
						continue
					}

					if t.Git != "" {
						fmt.Printf("Cloning template from: %s\n\n", t.Git)
					} else {
						fmt.Printf("Using template path: %s\n\n", t.Path)
					}

					path := Ask("Where do you want to copy the template (path)? ")

					PrintTitle("Variable replacement")

					for _, v := range t.Variables {
						// ask only if the key is not already in the map
						if _, ok := replacementMap[v]; !ok {
							replacementMap[v] = Ask(fmt.Sprintf("%s: ", v))
						}
					}

					// create a tmp dir to revert the operation if an error occurs
					dir, err := ioutil.TempDir("", "bob")
					if err != nil {
						return err
					}
					defer os.RemoveAll(dir) // clean up

					if t.Git != "" {
						_, err := git.PlainClone(dir, false, &git.CloneOptions{
							URL:      t.Git,
							Progress: os.Stdout,
						})
						if err != nil {
							return err
						}
						t.Path = dir
					} else {
						//TODO: Replace by actual golang code
						cmd := exec.Command("cp", "-R", t.Path, dir)
						err = cmd.Run()
						if err != nil {
							return err
						}
					}

					for k, v := range replacementMap {
						// remplace the folders and files names recursively
						var parseFiles func(string) error
						parseFiles = func(path string) error {
							files, err := ioutil.ReadDir(path)
							if err != nil {
								return err
							}

							for _, file := range files {
								fileName := file.Name()
								re := regexp.MustCompile(k)
								// replace matching key in file/folder name
								if re.MatchString(filepath.Join(path, fileName)) {
									replacedName := re.ReplaceAllString(fileName, v)
									err := os.Rename(filepath.Join(path, fileName), filepath.Join(path, replacedName))
									if err != nil {
										return err
									}
									fileName = replacedName
								}
								if file.IsDir() {
									// go deeper in recursion
									err = parseFiles(filepath.Join(path, fileName))
									if err != nil {
										return err
									}
								} else {
									// replace matching key in file
									b, err := ioutil.ReadFile(filepath.Join(path, fileName))
									if err != nil {
										return err
									}
									err = ioutil.WriteFile(filepath.Join(path, fileName), re.ReplaceAll(b, []byte(v)), file.Mode().Perm())
									if err != nil {
										return err
									}
								}
							}
							return nil
						}

						err = parseFiles(dir)
						if err != nil {
							return err
						}
					}

					absPath, err := filepath.Abs(path)
					if err != nil {
						return err
					}

					files, err := ioutil.ReadDir(dir)
					if err != nil {
						return err
					}

					for _, f := range files {
						err = os.Rename(filepath.Join(dir, f.Name()), filepath.Join(absPath, f.Name()))
						if err != nil {
							return err
						}
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
	fmt.Printf("\n==== %s ====\n\n", title)
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
