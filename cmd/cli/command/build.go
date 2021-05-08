package command

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/bob/pkg/config"
	"github.com/jcalmat/bob/pkg/file"
	"github.com/jcalmat/termui/v3"
	"github.com/jcalmat/termui/v3/widgets"
)

type Form struct {
	form            *ui.Form
	screen          *ui.Screen
	settings        config.Settings
	command         config.Command
	stringQuestions map[string]*widgets.TextField
	boolQuestions   map[string]*widgets.Checkbox
	skipMap         map[string]struct{}
}

// predefine functions
var funcMap = template.FuncMap{
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

func (c Command) Build(args ...string) {

	var infos strings.Builder

	globalConfig, err := c.ConfigApp.Parse()
	if err != nil {
		infos.WriteString(err.Error())
		return
	}

	cmd := args[0]

	command := globalConfig.Commands[cmd]

	buildForm := &Form{
		form:            ui.NewForm(),
		screen:          c.Screen,
		settings:        globalConfig.Settings,
		skipMap:         make(map[string]struct{}),
		stringQuestions: make(map[string]*widgets.TextField),
		boolQuestions:   make(map[string]*widgets.Checkbox),
	}

	buildForm.command = command

	buildForm.form.SetTitle(cmd)

	if command.Git != "" {
		infos.WriteString(fmt.Sprintf("Cloning template from: %s\n\n", command.Git))
	} else {
		infos.WriteString(fmt.Sprintf("Using template path: %s\n\n", command.Path))
	}
	path := widgets.NewTextField("Where do you want to copy this template? ")

	infos.WriteString(fmt.Sprintf("Current path: %s\n\n", file.GetWorkingDirectory()))

	closeButton := widgets.NewButton("Done", func() {
		err := buildForm.ProcessBuild(path)
		if err != nil {
			modale := ui.NewModale(fmt.Sprintln(err.Error()), ui.ModaleTypeErr)
			modale.Resize()
			modale.Render()
			c.Screen.SetModale(modale)
			return
		}
		modale := ui.NewModale(`
			Done
			Press Enter to quit
			Esc to get back to main menu
			`, ui.ModaleTypeInfo)
		modale.Resize()
		modale.Render()
		c.Screen.SetModale(modale)

		uiEvents := termui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "<C-c>":
				c.Screen.Stop()
				return
			case "<Enter>":
				c.Screen.Stop()
				return
			case "<Escape>":
				c.Screen.Restore()
				c.Screen.Restore()
				return
			}
		}
	})

	nodes := []*widgets.FormNode{
		{
			Item: path,
		},
		{
			Item: widgets.NewLabel(""),
		},
		{
			Item: widgets.NewLabel("==== Variable replacement ===="),
		},
		{
			Item: widgets.NewLabel(""),
		},
	}

	for _, v := range command.Variables {
		nodes = append(nodes, buildForm.parseQuestion(v))
	}

	nodes = append(nodes, &widgets.FormNode{
		Item: widgets.NewLabel(""),
	})

	nodes = append(nodes, &widgets.FormNode{
		Item: closeButton,
	})

	for _, v := range command.Skip {
		buildForm.skipMap[v] = struct{}{}
	}

	buildForm.form.SetNodes(nodes)
	buildForm.form.SetInfos(infos.String())
	c.Screen.SetForm(buildForm.form)
	buildForm.form.Render()
}

func (f *Form) ProcessBuild(path *widgets.TextField) error {
	replacementMap := make(map[string]interface{})

	dir, err := ioutil.TempDir("", "bob")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %s", err.Error())
	}
	defer os.RemoveAll(dir) // clean up

	cloneOpts := &git.CloneOptions{
		URL: f.command.Git,
		// Progress: &infos,
	}

	if _, err := os.Stat(f.settings.Git.SSH.PrivateKeyFile); err == nil {
		auth, err := ssh.NewPublicKeysFromFile("git", f.settings.Git.SSH.PrivateKeyFile, f.settings.Git.SSH.PrivateKeyPassword)
		if err != nil {
			return fmt.Errorf("generate publickeys failed: %s", err.Error())
		}
		cloneOpts.Auth = auth
	}
	if f.command.Git != "" {
		_, err := git.PlainClone(dir, false, cloneOpts)
		if err != nil {
			return fmt.Errorf("failed to clone template: %s", err.Error())
		}
		f.command.Path = dir
	} else {
		err := file.Copy(f.command.Path, dir)
		if err != nil {
			return fmt.Errorf("failed to copy template: %s", err.Error())
		}
	}

	for k, v := range f.stringQuestions {
		replacementMap[k] = v.Answer()
	}
	for k, v := range f.boolQuestions {
		replacementMap[k] = v.Answer()
	}

	// parseFiles replaces folders, files names and file content recursively
	var parseFiles func(string) error
	parseFiles = func(path string) error {
		folderFiles, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		for _, ff := range folderFiles {
			if _, ok := f.skipMap[ff.Name()]; ok {
				continue
			}

			filePath := filepath.Join(path, ff.Name())

			// # File/Directory name modifications
			var sb strings.Builder

			tmpl, err := template.New(filepath.Base(filePath)).Funcs(funcMap).Parse(filePath)
			if err != nil {
				return err
			}
			err = tmpl.Execute(&sb, replacementMap)
			if err != nil {
				return err
			}

			// replace if file doesn't already exist; i.e if we actually modified the file name
			if _, err := os.Stat(sb.String()); err != nil {
				err = os.Rename(filePath, sb.String())
				if err != nil {
					return err
				}
				filePath = filepath.Join(path, filepath.Base(sb.String()))
			}

			if ff.IsDir() {
				// go deeper in recursion
				err = parseFiles(filePath)
				if err != nil {
					return err
				}
			} else {
				// # File content modifications

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

	err = file.Move(dir, path.Answer(), f.command.Skip)
	if err != nil {
		return err
	}

	return nil
}

func (f *Form) parseQuestion(v config.Variable) *widgets.FormNode {
	question := fmt.Sprintf("%s: ", v.Name)
	if v.Desc != nil {
		question = *v.Desc
	}

	var node = &widgets.FormNode{}
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
			node.Nodes = append(node.Nodes, f.parseQuestion(s))
		}
	}

	node.Item = item
	return node
}
