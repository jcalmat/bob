package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jcalmat/bob/pkg/config"
	"gopkg.in/yaml.v2"
)

type App struct {
	ConfigFilePath string
}

var _ config.App = (*App)(nil)

func (a App) Parse() (config.C, error) {
	// open file
	file, err := os.Open(a.ConfigFilePath)
	if err != nil {
		return config.C{}, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return config.C{}, err
	}

	switch filepath.Ext(a.ConfigFilePath) {
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

func (a App) InitConfig() error {
	config := []byte(`
# Register your commands here
commands:

templates:

settings:
`)
	_, err := os.Stat(a.ConfigFilePath)
	if os.IsNotExist(err) {
		err := ioutil.WriteFile(a.ConfigFilePath, config, 0600)
		if err != nil {
			return err
		}
	} else {
		return errors.New("config file already exist")
	}
	return nil
}
