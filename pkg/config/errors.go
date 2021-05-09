package config

import "errors"

var (
	ErrConfigFileNotFound error = errors.New("config file not found")
)
