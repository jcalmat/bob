package cmd

import (
	"github.com/spf13/viper"
)

var Commands map[string]Command
var Templates map[string]Template

type Command struct {
	Alias     string   `yaml:"alias"`
	Templates []string `yaml:"templates"`
}

type Template struct {
	Path      string     `yaml:"path"`
	Git       string     `yaml:"git"`
	Variables []Variable `yaml:"variables"`
	Skip      []string   `yaml:"skip"`
}

type Variable struct {
	Name         string
	Type         Type
	Desc         *string
	Dependencies []string
	// Items []Item
}

type Item struct {
	Name  string
	Value interface{}
}

type Type string

const (
	String Type = "string"
	Array  Type = "array"
	Bool   Type = "bool"
)

func parseConfig() error {
	err := viper.UnmarshalKey("commands", &Commands)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("templates", &Templates)
	if err != nil {
		return err
	}
	return nil
}
