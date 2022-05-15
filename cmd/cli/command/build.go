package command

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/jcalmat/bob/cmd/cli/ui"
	"github.com/jcalmat/bob/pkg/config"
	"github.com/jcalmat/bob/pkg/file"
	"github.com/jcalmat/termui/v3"
	"github.com/jcalmat/termui/v3/widgets"
)

type Builder struct {
	form            *ui.Form
	screen          *ui.Screen
	settings        config.Settings
	command         config.Command
	configApp       config.App
	temporaryPath   string
	stringQuestions map[string]*widgets.TextField
	boolQuestions   map[string]*widgets.Checkbox
	skipMap         map[string]struct{}
	// forcedPath will prevent the user to change the path
	forcedPath *string
}

// predefine functions used in go-template
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

	builder := &Builder{
		form:            ui.NewForm(),
		screen:          c.Screen,
		settings:        globalConfig.Settings,
		skipMap:         make(map[string]struct{}),
		stringQuestions: make(map[string]*widgets.TextField),
		boolQuestions:   make(map[string]*widgets.Checkbox),
		configApp:       c.ConfigApp,
		command:         globalConfig.Commands[cmd],
	}

	builder.form.SetTitle(cmd)

	builder.SetupBuild(c.Screen)
}

func (b Builder) SetupBuild(screen *ui.Screen) {
	var infos strings.Builder

	err := b.ParseSubConfigSpecs()
	if err != nil {
		screen.RenderModale(err.Error(), ui.ModaleTypeErr)
		return
	}

	if b.command.Git != "" {
		infos.WriteString(fmt.Sprintf("Cloning template from: %s\n\n", b.command.Git))
	} else {
		infos.WriteString(fmt.Sprintf("Using template path: %s\n\n", b.command.Path))
	}

	pathWidget := widgets.NewTextField("Where do you want to copy this template? ")

	infos.WriteString(fmt.Sprintf("Current path: %s\n\n", file.GetWorkingDirectory()))

	closeButton := widgets.NewButton("Done", func() {
		err := b.ProcessBuild(pathWidget)
		if err != nil {
			screen.RenderModale(err.Error(), ui.ModaleTypeErr)
			return
		}

		screen.RenderModale(`
			Done
			Press Enter to quit
			Esc to get back to main menu
		`, ui.ModaleTypeInfo)

		uiEvents := termui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "<C-c>":
				screen.Stop()
				return
			case "<Enter>":
				screen.Stop()
				return
			case "<Escape>":
				screen.Restore()
				screen.Restore()
				return
			}
		}
	})

	nodes := make([]*widgets.FormNode, 0)
	// path is not forced, ask the user to choose
	if b.forcedPath == nil {
		nodes = append(nodes, []*widgets.FormNode{
			{
				Item: pathWidget,
			},
			{
				Item: widgets.NewLabel(""),
			},
		}...)
	}

	nodes = append(nodes, []*widgets.FormNode{
		{
			Item: widgets.NewLabel(" ----------------------"),
		},
		{
			Item: widgets.NewLabel(" -  Customize fields  -"),
		},
		{
			Item: widgets.NewLabel(" ----------------------"),
		},
		{
			Item: widgets.NewLabel(""),
		},
	}...)

	for _, v := range b.command.Variables {
		nodes = append(nodes, b.parseQuestion(v))
	}

	nodes = append(nodes, &widgets.FormNode{
		Item: widgets.NewLabel(""),
	})

	nodes = append(nodes, &widgets.FormNode{
		Item: closeButton,
	})

	for _, v := range b.command.Skip {
		b.skipMap[v] = struct{}{}
	}

	b.form.SetNodes(nodes)
	b.form.SetInfos(infos.String())
	screen.SetForm(b.form)
	b.form.Render()
}

// ParseSubconfigSpecs will look for bobconfig file inside the template itself
// and parse it if found
func (b *Builder) ParseSubConfigSpecs() error {
	err := b.Clone()
	if err != nil {
		return err
	}
	defer os.RemoveAll(b.temporaryPath) // clean up

	s, err := b.configApp.ParseSpecs(b.temporaryPath)
	if err != nil {
		// ignore if there is no subconfig file
		if !errors.Is(err, config.ErrConfigFileNotFound) {
			return err
		}
		return nil
	}

	b.command.Specs = s

	return nil

}

// Clone either clone the project from a git repository either copy it from a
// local folder to a temporary folder. It doesn't remove the temporary folder,
// it is the caller's responsibility to remove the directory when no longer needed.
func (b *Builder) Clone() error {
	dir, err := ioutil.TempDir("", "bob")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %s", err.Error())
	}

	gitSplits := strings.Split(b.command.Git, " ")
	var url, branch string
	if len(gitSplits) == 2 {
		url = gitSplits[0]
		branch = gitSplits[1]
	} else {
		url = gitSplits[0]
		branch = "master"
	}

	cloneOpts := &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		// Progress: &infos,
	}

	if _, err := os.Stat(b.settings.Git.SSH.PrivateKeyFile); err == nil {
		auth, err := ssh.NewPublicKeysFromFile("git", b.settings.Git.SSH.PrivateKeyFile, b.settings.Git.SSH.PrivateKeyPassword)
		if err != nil {
			return fmt.Errorf("generate publickeys failed: %s", err.Error())
		}
		cloneOpts.Auth = auth
	}
	if b.command.Git != "" {
		_, err := git.PlainClone(dir, false, cloneOpts)
		if err != nil {
			return fmt.Errorf("failed to clone template: %s", err.Error())
		}
		b.command.Path = dir
	} else {
		err := file.Copy(b.command.Path, dir)
		if err != nil {
			return fmt.Errorf("failed to copy template: %s", err.Error())
		}
	}

	b.temporaryPath = dir

	return nil
}

// ProcessBuild clones the template and replace all the actionable items
// to the files contents, files names and directory names recursively.
func (b *Builder) ProcessBuild(path *widgets.TextField) error {
	replacementMap := make(map[string]interface{})

	err := b.Clone()
	if err != nil {
		return err
	}
	defer os.RemoveAll(b.temporaryPath)

	for k, v := range b.stringQuestions {
		replacementMap[k] = v.Answer()
	}
	for k, v := range b.boolQuestions {
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
			if _, ok := b.skipMap[ff.Name()]; ok {
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

	err = parseFiles(b.temporaryPath)
	if err != nil {
		return err
	}

	finalPath := path.Answer()
	if b.forcedPath != nil {
		finalPath = *b.forcedPath
	}

	err = file.Move(b.temporaryPath, finalPath, b.command.Skip)
	if err != nil {
		return err
	}

	return nil
}

// parseQuestion convert a config.Variable to a termui FormNode
func (b *Builder) parseQuestion(v config.Variable) *widgets.FormNode {
	question := v.Name
	if v.Format != nil {
		question = *v.Format
	}

	var node = &widgets.FormNode{}
	var item widgets.FormItem

	switch v.Type {
	case config.String:
		question = fmt.Sprintf("%s: ", question)
		textfield := widgets.NewTextField(question)
		item = textfield
		b.stringQuestions[v.Name] = textfield
	case config.Bool:
		checkbox := widgets.NewCheckbox(question, false)
		item = checkbox
		b.boolQuestions[v.Name] = checkbox
	case config.Array:
	//TODO:
	default:
		//default case is string
		textfield := widgets.NewTextField(fmt.Sprintf("%s: ", v.Name))
		item = textfield
		b.stringQuestions[v.Name] = textfield
	}

	if v.Dependencies != nil {
		for _, s := range v.Dependencies {
			node.Nodes = append(node.Nodes, b.parseQuestion(s))
		}
	}

	node.Item = item
	return node
}
