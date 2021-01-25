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

type Form struct {
	form         *form.Form
	questionsMap map[string]*form.FormItem
}

func (c Command) Build(args ...string) {
	globalConfig, err := c.ConfigApp.Parse()
	if err != nil {
		c.Logger.Err(err).Msg("")
		return
	}

	cmd := args[0]

	command := globalConfig.Commands[cmd]
	templates := globalConfig.Templates

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

	f := &Form{
		form:         form.NewForm(),
		questionsMap: make(map[string]*form.FormItem),
	}

	f.form.AddItem(form.NewLabel(` ______     ______     ______`))
	f.form.AddItem(form.NewLabel(`/\  == \   /\  __ \   /\  == \`))
	f.form.AddItem(form.NewLabel(`\ \  __<   \ \ \/\ \  \ \  __<`))
	f.form.AddItem(form.NewLabel(` \ \_____\  \ \_____\  \ \_____\`))
	f.form.AddItem(form.NewLabel(`  \/_____/   \/_____/   \/_____/`))
	f.form.AddItem(form.NewLabel(""))
	f.form.AddItem(form.NewLabel(fmt.Sprintf("> %s", cmd)))
	f.form.AddItem(form.NewLabel(""))

	for _, s := range command.Templates {

		replacementMap := make(map[string]interface{})

		skipMap := make(map[string]struct{})

		t, ok := templates[s]
		if !ok {
			io.Info("Template %s not found, skipping\n\n", s)
			continue
		}
		f.form.AddItem(form.NewLabel(fmt.Sprintf("Current template: %s", s)))

		if t.Git != "" {
			f.form.AddItem(form.NewLabel(fmt.Sprintf("Cloning template from: %s", t.Git)))
			f.form.AddItem(form.NewLabel(""))
		} else {
			f.form.AddItem(form.NewLabel(fmt.Sprintf("Using template path: %s", t.Path)))
			f.form.AddItem(form.NewLabel(""))
		}
		path := form.NewTextField("Where do you want to copy this template? ")
		f.form.AddItem(path)

		f.form.AddItem(form.NewLabel(fmt.Sprintf("Current path: %s", file.GetWorkingDirectory())))
		f.form.AddItem(form.NewLabel(""))
		f.form.AddItem(form.NewLabel("==== Variable replacement ===="))

		for _, v := range t.Variables {
			f.form.AddItem(f.parseQuestion(v))
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

		f.form.Run()

		for k, v := range f.questionsMap {
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

		err = file.Move(dir, path.Answer().(string), t.Skip)
		if err != nil {
			c.Logger.Error().Err(err).Msg("")
			return
		}
	}

	io.Title("Done")
}

func (f *Form) parseQuestion(v config.Variable) *form.FormItem {
	question := fmt.Sprintf("%s: ", v.Name)
	if v.Desc != nil {
		question = *v.Desc
	}

	// var item form.Item
	var item *form.FormItem

	switch v.Type {
	case config.String:
		item = form.NewTextField(question)
	case config.Bool:
		item = form.NewCheckbox(question, false)
	case config.Array:
	//TODO:
	default:
		//default case is string
		item = form.NewTextField(fmt.Sprintf("%s: ", v.Name))
	}

	if v.Sub != nil {
		for _, s := range v.Sub {
			item.AddItem(f.parseQuestion(s))
		}
	}
	f.questionsMap[v.Name] = item

	return item
}
