package command

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/jcalmat/bob/pkg/config"
	"github.com/jcalmat/bob/pkg/file"
	"github.com/jcalmat/bob/pkg/io"
	"github.com/jcalmat/form"
)

func (c Command) Build(args ...string) {
	globalConfig, err := c.ConfigApp.Parse()
	if err != nil {
		c.Logger.Err(err).Msg("")
		return
	}

	cmd := args[0]

	command := globalConfig.Commands[cmd]
	templates := globalConfig.Templates

	io.Title(cmd)

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

	bobform := form.NewForm()

	for _, s := range command.Templates {

		replacementMap := make(map[string]interface{})
		textFieldMap := make(map[string]*form.TextField)
		checkboxMap := make(map[string]*form.Checkbox)

		skipMap := make(map[string]struct{})

		t, ok := templates[s]
		if !ok {
			io.Info("Template %s not found, skipping\n\n", s)
			continue
		}

		if t.Git != "" {
			io.Info("Cloning template from: %s\n\n", t.Git)
		} else {
			io.Info("Using template path: %s\n\n", t.Path)
		}

		path := form.NewTextField("Where do you want to copy this template? ")
		bobform.AddItem(path)

		bobform.AddItem(form.NewLabel(fmt.Sprintf("Current path: %s", file.GetWorkingDirectory())))
		bobform.AddItem(form.NewLabel(""))
		bobform.AddItem(form.NewLabel("==== Variable replacement ===="))

		for _, v := range t.Variables {
			question := fmt.Sprintf("%s: ", v.Name)
			if v.Desc != nil {
				question = *v.Desc
			}

			switch v.Type {
			case config.String:
				textFieldMap[v.Name] = form.NewTextField(question)
				item := bobform.AddItem(textFieldMap[v.Name])
				for _, dep := range v.Sub {
					question := fmt.Sprintf("%s: ", dep.Name)
					if dep.Desc != nil {
						question = *dep.Desc
					}
					item.AddSubItem(form.NewCheckbox(question, false))
				}

			case config.Bool:
				checkboxMap[v.Name] = form.NewCheckbox(question, false)
				item := bobform.AddItem(checkboxMap[v.Name])
				for _, dep := range v.Sub {
					question := fmt.Sprintf("%s: ", dep.Name)
					if dep.Desc != nil {
						question = *dep.Desc
					}
					item.AddSubItem(form.NewCheckbox(question, false))
				}
			case config.Array:
			//TODO:
			default:
				//default case is string
				textFieldMap[v.Name] = form.NewTextField(fmt.Sprintf("%s: ", v.Name))
				bobform.AddItem(textFieldMap[v.Name])
			}
		}

		for _, v := range t.Skip {
			skipMap[v] = struct{}{}
		}

		// create a tmp dir to revert the operation if an error occurs
		dir, err := ioutil.TempDir("", "bob")
		if err != nil {
			c.Logger.Error().Err(err).Msg("")
			return
		}
		defer os.RemoveAll(dir) // clean up

		if t.Git != "" {
			_, err := git.PlainClone(dir, false, &git.CloneOptions{
				URL:      t.Git,
				Progress: os.Stdout,
			})
			if err != nil {
				c.Logger.Error().Err(err).Msg("")
				return
			}
			t.Path = dir
		} else {
			err := file.Copy(t.Path, dir)
			if err != nil {
				c.Logger.Error().Err(err).Msg("")
				return
			}
		}

		bobform.Run()

		for k, v := range textFieldMap {
			replacementMap[k] = v.Answer()
		}
		for k, v := range checkboxMap {
			replacementMap[k] = v.Answer()
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
				c.Logger.Error().Err(err).Msg("")
				return
			}

		}

		err = file.Move(dir, path.Answer(), t.Skip)
		if err != nil {
			c.Logger.Error().Err(err).Msg("")
			return
		}
	}

	io.Title("Done")
}
