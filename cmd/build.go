package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/jcalmat/bob/pkg/file"
	"github.com/jcalmat/bob/pkg/io"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build a project from a specified template",
	Long: `- select a template and clone it in the desired folder
- replace the variables set in your config file to customize your template`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	updateConfig()
	buildCmd.ResetCommands()
	for _, c := range Commands {
		c := c
		buildCmd.AddCommand(&cobra.Command{
			Use: c.Alias,
			RunE: func(cmd *cobra.Command, args []string) error {

				io.Title(c.Alias)

				// predefine functions
				funcMap := template.FuncMap{
					// truncate the i first chars of s
					"short": func(s string, i int) string {
						runes := []rune(s)
						if len(runes) > i {
							return string(runes[:i])
						}
						return s
					},

					// upcase the string s
					"upcase": func(s string) string {
						return strings.ToUpper(s)
					},

					"title": func(s string) string {
						return strings.Title(s)
					},
				}

				for _, s := range c.Templates {

					replacementMap := make(map[string]interface{})

					skipMap := make(map[string]struct{})

					t, ok := Templates[s]
					if !ok {
						io.Info("Template %s not found, skipping\n\n", s)
						continue
					}

					if t.Git != "" {
						io.Info("Cloning template from: %s\n\n", t.Git)
					} else {
						io.Info("Using template path: %s\n\n", t.Path)
					}

					path := io.Ask("Where do you want to copy the template (path)? ")

					io.Title("Variable replacement")

					for _, v := range t.Variables {
						skip := false
						for _, dependency := range v.Dependencies {
							if _, ok := replacementMap[dependency]; !ok {
								skip = true
								break
							}
						}
						if skip {
							continue
						}

						question := fmt.Sprintf("%s: ", v.Name)
						if v.Desc != nil {
							question = *v.Desc
						}

						switch v.Type {
						case String:
							replacementMap[v.Name] = io.Ask(question)
						case Bool:
							replacementMap[v.Name] = io.AskBool(question)
						case Array:
						//TODO:
						default:
							//default case is string
							replacementMap[v.Name] = io.Ask(fmt.Sprintf("%s: ", v.Name))
						}
					}

					for _, v := range t.Skip {
						skipMap[v] = struct{}{}
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
						err := file.Copy(t.Path, dir)
						if err != nil {
							return err
						}
					}

					for k, v := range replacementMap {
						// remplace the folders and files names recursively
						var parseFiles func(string) error
						parseFiles = func(path string) error {
							folderFiles, err := ioutil.ReadDir(path)
							if err != nil {
								return err
							}

							for _, f := range folderFiles {
								if _, ok := skipMap[f.Name()]; ok {
									continue
								}

								filePath := filepath.Join(path, f.Name())

								if _, ok := v.(string); ok {
									replacedName, err := file.RenameFile(filePath, fmt.Sprintf("{{%s}}", k), v.(string))
									if err != nil {
										return err
									}
									filePath = filepath.Join(path, filepath.Base(replacedName))
								}
								if f.IsDir() {
									// go deeper in recursion
									err = parseFiles(filePath)
									if err != nil {
										return err
									}
								} else {
									// use go templates to replace
									tt, err := template.New(filepath.Base(filePath)).Funcs(funcMap).ParseFiles(filePath)
									if err != nil {
										return err
									}
									fd, err := os.Create(filePath)
									if err != nil {
										return err
									}
									defer fd.Close()

									err = tt.Execute(fd, replacementMap)
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

					err = file.Move(dir, path, t.Skip)
					if err != nil {
						return err
					}
				}

				io.Title("Done")

				return nil
			},
		})
	}
}
