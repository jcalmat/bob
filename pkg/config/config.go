package config

type C struct {
	Commands map[string]Command
	Settings Settings `yaml:"settings" json:"settings"`
}

type Command struct {
	Description string `yaml:"description" json:"description"`
	Path        string `yaml:"path" json:"path"`
	Git         string `yaml:"git" json:"git"`
	Specs       `yaml:",inline" json:",inline"`
}

type Specs struct {
	Variables []Variable `yaml:"vars" json:"vars"`
	Skip      []string   `yaml:"skip" json:"skip"`
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

type Variable struct {
	Name         string     `yaml:"name" json:"name"`
	Type         Type       `yaml:"type" json:"type"`
	Format       *string    `yaml:"format,omitempty" json:"format,omitempty"`
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
	ParseSpecs(string) (Specs, error)
}
