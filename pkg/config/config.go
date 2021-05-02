package config

type C struct {
	Commands  map[string]Command
	Templates map[string]Template
	Settings  Settings `yaml:"settings" json:"settings"`
}

type Command struct {
	Description string   `yaml:"description" json:"description"`
	Templates   []string `yaml:"templates" json:"templates"`
}

type Settings struct {
	Git GitSettings `yaml:"git" json:"git"`
}

type GitSettings struct {
	SSH SSHGitSettings `yaml:"ssh" json:"ssh"`
}

type SSHGitSettings struct {
	PrivateKeyFile     string `yaml:"privateKeyFile" json:"privateKeyFile"`
	PrivateKeyPassword string `yaml:"privateKeyPassword" json:"privateKeyPassword"`
}

type Template struct {
	Path      string     `yaml:"path" json:"path"`
	Git       string     `yaml:"git" json:"git"`
	Variables []Variable `yaml:"variables" json:"variables"`
	Skip      []string   `yaml:"skip" json:"skip"`
}

type Variable struct {
	Name         string
	Type         Type
	Desc         *string
	Dependencies []Variable `yaml:"deps" json:"deps"`
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

type App interface {
	Parse() (C, error)
	InitConfig() error
}
