package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jcalmat/bob/pkg/config"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type App struct {
}

var _ config.App = (*App)(nil)

// Parse parse a bobconfig
func (a App) Parse() (config.C, error) {

	configPath, err := a.getConfigFile("~")
	if err != nil {
		return config.C{}, err
	}

	// open file
	file, err := os.Open(configPath)
	if err != nil {
		return config.C{}, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return config.C{}, err
	}

	switch filepath.Ext(configPath) {
	case ".yaml", ".yml":
		return parseYaml(content)
	case ".json":
		return parseJSON(content)
	}

	return config.C{}, errors.New("File extension not handled")
}

func parseYaml(content []byte) (config.C, error) {
	var config config.C

	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func parseJSON(content []byte) (config.C, error) {
	var config config.C

	err := json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// getConfigFile retrieves the config file, wether it's a yaml or a json file
func (a App) getConfigFile(dir string) (string, error) {
	absPath, _ := homedir.Expand(dir)
	var configPath string

	if _, err := os.Stat(filepath.Join(absPath, ".bobconfig.yml")); err == nil {
		configPath = filepath.Join(absPath, ".bobconfig.yml")
	} else if _, err := os.Stat(filepath.Join(absPath, ".bobconfig.yaml")); err == nil {
		configPath = filepath.Join(absPath, ".bobconfig.yaml")
	} else if _, err := os.Stat(filepath.Join(absPath, ".bobconfig.json")); err == nil {
		configPath = filepath.Join(absPath, ".bobconfig.json")
	} else {
		return "", config.ErrConfigFileNotFound
	}

	return configPath, nil
}

// ParseSpecs parses only the specs part of a bobconfig file
func (a App) ParseSpecs(dir string) (config.Specs, error) {

	configPath, err := a.getConfigFile(dir)
	if err != nil {
		return config.Specs{}, err
	}

	// open file
	file, err := os.Open(configPath)
	if err != nil {
		return config.Specs{}, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return config.Specs{}, err
	}

	switch filepath.Ext(configPath) {
	case ".yaml", ".yml":
		return parseYamlSpecs(content)
	case ".json":
		return parseJSONSpecs(content)
	}

	return config.Specs{}, errors.New("File extension not handled")
}

func parseYamlSpecs(content []byte) (config.Specs, error) {
	var specs config.Specs

	err := yaml.Unmarshal(content, &specs)
	if err != nil {
		return specs, err
	}

	return specs, nil
}

func parseJSONSpecs(content []byte) (config.Specs, error) {
	var specs config.Specs

	err := json.Unmarshal(content, &specs)
	if err != nil {
		return specs, err
	}

	return specs, nil
}

func (a App) InitYamlConfig() error {
	path, err := a.getConfigFile("~")
	if err == nil {
		return fmt.Errorf("config file already exist at %s", path)
	}

	config := []byte(`# Register your commands here
commands:
  example:
    path: "/path/to/templates_dir"
settings:
git:
  ssh:
    privateKeyFile: "/home/user/.ssh/id_rsa"
    privateKeyPassword: ""
`)

	absPath, _ := homedir.Expand("~")
	err = ioutil.WriteFile(filepath.Join(absPath, ".bobconfig.yml"), config, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (a App) InitJSONConfig() error {
	path, err := a.getConfigFile("~")
	if err == nil {
		return fmt.Errorf("config file already exist at %s", path)
	}

	config := []byte(`{
	"commands": [
		{
			"example": {
				"path": "/path/to/templates_dir"
			}
		}
	],
	"settings": {
		"git": {
			"ssh": {
				"privateKeyFile": "/path/to/ssh/key",
				"privateKeyPassword": "/password/for/this/key"
			}
		}
	}
}`)

	absPath, _ := homedir.Expand("~")
	err = ioutil.WriteFile(filepath.Join(absPath, ".bobconfig.json"), config, 0600)
	if err != nil {
		return err
	}

	return nil
}
