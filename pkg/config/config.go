package config

type C struct {
	Commands  map[string]Command
	Templates map[string]Template
}
type Command struct {
	// Alias     string   `yaml:"alias" json:"alias"`
	Description string   `yaml:"description" json:"description"`
	Templates   []string `yaml:"templates" json:"templates"`
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
