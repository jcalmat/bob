package command

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/bob/pkg/config"
	"github.com/jcalmat/bob/pkg/file"
	"github.com/jcalmat/bob/pkg/io"
	"github.com/jcalmat/termui/v3"
	"github.com/jcalmat/termui/v3/widgets"
)

type Form struct {
	form *ui.Form
	// questionsMap       map[string]*widgets.FormItem
	stringQuestions map[string]*widgets.TextField
	boolQuestions   map[string]*widgets.Checkbox
}

type ConfigInfo struct {
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
		form: ui.NewForm(),
		// questionsMap: make(map[string]*widgets.FormItem),
		stringQuestions: make(map[string]*widgets.TextField),
		boolQuestions:   make(map[string]*widgets.Checkbox),
	}

	form := ui.NewForm()
	var nodes []*widgets.FormNode
	var infos strings.Builder

	for _, s := range command.Templates {

		replacementMap := make(map[string]interface{})

		skipMap := make(map[string]struct{})

		t, ok := templates[s]
		if !ok {
			io.Info("Template %s not found, skipping\n\n", s)
			continue
		}
		// f.form.AddItem(form.NewLabel(fmt.Sprintf("Current template: %s", s)))

		form.SetTitle(s)

		if t.Git != "" {
			infos.WriteString(fmt.Sprintf("Cloning template from: %s\n\n", t.Git))
			// f.form.AddItem(form.NewLabel(fmt.Sprintf("Cloning template from: %s", t.Git)))
			// f.form.AddItem(form.NewLabel(""))
		} else {
			infos.WriteString(fmt.Sprintf("Using template path: %s\n\n", t.Path))
			// f.form.AddItem(form.NewLabel(fmt.Sprintf("Using template path: %s", t.Path)))
			// f.form.AddItem(form.NewLabel(""))
		}
		path := widgets.NewTextField("Where do you want to copy this template? ")

		infos.WriteString(fmt.Sprintf("Current path: %s\n\n", file.GetWorkingDirectory()))

		nodes = []*widgets.FormNode{
			{
				Item: path,
			},
			{
				Item: widgets.NewLabel(""),
			},
			{
				Item: widgets.NewLabel("==== Variable replacement ===="),
			},
		}

		// path := form.NewTextField("Where do you want to copy this template? ")
		// f.form.AddItem(path)

		// f.form.AddItem(form.NewLabel(fmt.Sprintf("Current path: %s", file.GetWorkingDirectory())))
		// f.form.AddItem(form.NewLabel(""))
		// f.form.AddItem(form.NewLabel("==== Variable replacement ===="))
		for _, v := range t.Variables {
			nodes = append(nodes, f.parseQuestion2(v))
			// f.form.AddItem(f.parseQuestion(v))
		}

		for _, v := range t.Skip {
			skipMap[v] = struct{}{}
		}

		form.SetNodes(nodes)
		form.SetInfos(infos.String())
		c.Screen.SetForm(form)
		form.Render()

		var close bool
		uiEvents := termui.PollEvents()
		for {
			e := <-uiEvents
			form.Content.HandleKeyboard(e)
			switch e.ID {
			case "<C-c>":
				close = true
			case "<Down>":
				form.Content.ScrollDown()
			case "<Up>":
				form.Content.ScrollUp()
			case "<Enter>":
				form.Content.ToggleExpand()
				form.Content.ScrollDown()
			}
			form.Render()
			if close {
				break
			}
		}

		// create a tmp dir to revert the operation if an error occurs
		dir, err := ioutil.TempDir("", "bob")
		if err != nil {
			infos.WriteString(fmt.Sprintf("failed to create temp dir: %s\n", err.Error()))
			form.SetInfos(infos.String())
			// c.Logger.Error().Err(err).Msg("")
			return
		}
		defer os.RemoveAll(dir) // clean up

		if t.Git != "" {
			_, err := git.PlainClone(dir, false, &git.CloneOptions{
				URL: t.Git,
				// Progress: &infos,
			})
			if err != nil {
				infos.WriteString(fmt.Sprintf("failed to clone template: %s\n", err.Error()))
				form.SetInfos(infos.String())
				// c.Logger.Error().Err(err).Msg("")
				return
			}
			t.Path = dir
		} else {
			err := file.Copy(t.Path, dir)
			if err != nil {
				infos.WriteString(fmt.Sprintf("failed to copy template: %s\n", err.Error()))
				form.SetInfos(infos.String())
				// c.Logger.Error().Err(err).Msg("")
				return
			}
		}

		// err = f.form.Run()
		// if err != nil {
		// 	if errors.Is(err, form.ErrUserCancelRequest) {
		// 		fmt.Println("bye")
		// 		return
		// 	}
		// 	c.Logger.Error().Err(err).Msg("failed to run bob")
		// 	return
		// }

		for k, v := range f.stringQuestions {
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

	infos.WriteString("\nDone")
	infos.WriteString("\nPress ESC to get back to main menu")
	form.SetInfos(infos.String())
	form.Render()

	// c.Logger.Info().Msg("Done")
	// io.Title("Done")
}

// func (f *Form) parseQuestion(v config.Variable) *form.FormItem {
// 	question := fmt.Sprintf("%s: ", v.Name)
// 	if v.Desc != nil {
// 		question = *v.Desc
// 	}

// 	// var item form.Item
// 	var item *form.FormItem

// 	switch v.Type {
// 	case config.String:
// 		item = form.NewTextField(question)
// 	case config.Bool:
// 		item = form.NewCheckbox(question, false)
// 	case config.Array:
// 	//TODO:
// 	default:
// 		//default case is string
// 		item = form.NewTextField(fmt.Sprintf("%s: ", v.Name))
// 	}

// 	if v.Dependencies != nil {
// 		for _, s := range v.Dependencies {
// 			item.AddItem(f.parseQuestion(s))
// 		}
// 	}
// 	f.questionsMap[v.Name] = item

// 	return item
// }

func (f *Form) parseQuestion2(v config.Variable) *widgets.FormNode {
	question := fmt.Sprintf("%s: ", v.Name)
	if v.Desc != nil {
		question = *v.Desc
	}

	var node = &widgets.FormNode{}
	// var item form.Item
	var item widgets.FormItem

	switch v.Type {
	case config.String:
		textfield := widgets.NewTextField(question)
		item = textfield
		f.stringQuestions[v.Name] = textfield
	case config.Bool:
		checkbox := widgets.NewCheckbox(question, false)
		item = checkbox
		f.boolQuestions[v.Name] = checkbox
	case config.Array:
	//TODO:
	default:
		//default case is string
		textfield := widgets.NewTextField(fmt.Sprintf("%s: ", v.Name))
		item = textfield
		f.stringQuestions[v.Name] = textfield
	}

	if v.Dependencies != nil {
		for _, s := range v.Dependencies {
			node.Nodes = append(node.Nodes, f.parseQuestion2(s))
			// item.AddItem(f.parseQuestion(s))
		}
	}
	// f.questionsMap[v.Name] = item

	node.Item = item
	return node
}
