package config

type C struct {
	Commands  map[string]Command
	Templates map[string]Template
}
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
	Name string
	Type Type
	Desc *string
	Sub  []Variable
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
