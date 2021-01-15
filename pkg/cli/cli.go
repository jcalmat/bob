package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Router struct {
	Logger  zerolog.Logger
	Command *Command
}

func (r *Router) Handle() {
	if r.Command == nil {
		r.Logger.Error().Msg("Router has not been initialized")
	}
	r.Command.handle()
}

type Flag struct {
	Key     string
	Aliases []string
	Handler func()
}

type Command struct {
	Key         string
	Aliases     []string
	Description string
	parent      *Command
	children    []*Command
	Handler     func(...string)
	Flags       []Flag
	dep         int
	logger      zerolog.Logger
}

func Init(logger zerolog.Logger) *Router {
	r := &Router{
		Logger: logger,
		Command: &Command{
			Key:    os.Args[0],
			dep:    0,
			logger: logger,
		},
	}

	r.Command.addHelpFlag()

	return r
}

func (c *Command) addHelpFlag() {
	c.Flags = append(c.Flags, Flag{
		Key:     "--help",
		Aliases: []string{"-h"},
		Handler: c.displayHelp,
	})
}

func (c *Command) AddCommand(cmd *Command) *Command {
	cmd.parent = c
	cmd.dep = c.dep + 1
	c.children = append(c.children, cmd)
	cmd.logger = c.logger
	cmd.addHelpFlag()
	return cmd
}

func (c *Command) handle() {
	if len(c.children) > 0 {
		for _, child := range c.children {
			if child.dep >= len(os.Args) {
				break
			}
			if child.Key == os.Args[child.dep] {
				child.handle()
				return
			}
			for _, alias := range child.Aliases {
				if alias == os.Args[child.dep] {
					child.handle()
					return
				}
			}
		}
		c.displayHelp()
	} else {
		if c.Handler != nil {
			c.Handler(c.Key)
		} else {
			fmt.Printf("no command available for %s", c.Key)
			c.logger.Debug().Msgf("No command available for %s", c.Key)
		}
	}
}

func (c Command) displayHelp() {
	if len(c.children) > 0 {
		if c.Description != "" {
			fmt.Printf("\n%s\n\n", c.Description)
		}
		fmt.Printf("Usage:\n")
		root := &c
		parent := &c
		var path string
		for {
			if parent == nil {
				break
			}
			path = fmt.Sprintf("%s %s", parent.Key, path)
			root = parent
			parent = parent.parent
		}
		path = strings.TrimSpace(path)
		if len(c.children) > 0 {
			fmt.Printf("  %s [command]\n", path)
		}
		if len(c.Flags) > 0 {
			fmt.Printf("  %s [flags]\n", path)
		}
		fmt.Println()

		if len(c.children) > 0 {
			fmt.Printf("Available Commands for \"%s\":\n", path)
			for _, child := range c.children {
				fmt.Printf("  %-12s\t%s\n", child.Key, child.Description)
			}
		}
		fmt.Println()
		if len(c.Flags) > 0 {
			fmt.Print("Flags:\n  ")
			for _, flag := range c.Flags {
				fmt.Printf("%s", flag.Key)
				for _, alias := range flag.Aliases {
					fmt.Printf(", %s", alias)
				}
			}
		}
		fmt.Printf("\n\nUse \"%s [command] --help\" for more information about a command.\n", root.Key)
	}
}
