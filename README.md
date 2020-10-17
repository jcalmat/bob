# Bob - The Boilerplate Builder

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jcalmat/routine)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Linter](https://github.com/jcalmat/bob/workflows/golangci-lint/badge.svg)

Bob is a CLI tool used to generate boilerplate code from templates made beforehand.

## Setup

### Binary

Download the latest [release](https://github.com/jcalmat/bob/releases), unzip and put the binary somewhere on your PATH.

### Building

#### Requirements

* golang 1.13.x or later

```bash
$> go get github.com/jcalmat/bob
$> bob init
```

_OR_

```
$> git clone git@github.com:jcalmat/bob.git
$> cd bob
$> go install
$> bob init
```

### Configuration

Bob settings should be stored in your home directory in a `.bobconfig.yml` file and contain the following fields

```yaml
# Register your commands here
commands:
  microservice:
    alias: "microservice"   #alias is the key provided to bob to match this command
    templates:
      - my_microservice     #templates is an array containing one or multiple templates used during this command
 
  
templates:
  my_microservice:
    #you can either provide a git link to be cloned or a template already in your local environment
    #git: "github.com/USERNAME/templates/microservice.git"
    path: "/path/to/your/template"
    variables:  #variables are the variables to replace in your template
      - "service"
    skip: # files or folders to ignore, usefull when cloning a git project for example
      - ".git"
    

```

## Templates

The templates are prebuilt pieces of code with variables to replace.

### **Variables format**

Bob uses go templates to parse and replace the variables, thus these variables must be surrounded with double brackets `{{VARIABLE}}`.

For more information about the format, here is a [cheat sheet](https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet).

**Custom functions**

To ease your development and avoid variables duplication, bob has custom formatting methods:

- **`short`** will truncate the x first characters of your variable

```go
// template:
type {{short .package 1}} struct {
}

// input:
// package = user
// output
type u struct {
}
```

- **`upcase`** will upcase your variable
- **`title`** will return a copy of the string s with all Unicode letters that begin words mapped to their Unicode title case

You can also pile the functions up

```go
//template:

// {{short .package 1 | upcase}} represents an event's {{.package}}.
type {{short .package 1 | upcase}} struct {
}

// input:
// package = user
// output:
// U represents an event's user.
type U struct {
}
```

### Conditional templates

To be more flexible, you can add conditions to your template to ask the user if he wants to add a particular piece of code.

```go
type Store interface {
	{{if .insert}}
	Insert()
	{{end}}
}
```

This will be in final code only if the variable `insert` is defined.

`insert` can be a boolean or a string. If the string is empty, the block won't be included.

yaml example:

```go
templates:
  microservice_pkg:
    git: "https://github.com/jcalmat/bob/example/templates/main.go"
    variables:
      - name: "insert"
        type: "bool"
        desc: "Add Insert method to the store? [y/n] " // desc will override the replacement question asked to the user
    skip:
      - ".git"
```

### Help

```bash
$> bob --help

 ______     ______     ______
/\  == \   /\  __ \   /\  == \
\ \  __<   \ \ \/\ \  \ \  __<
 \ \_____\  \ \_____\  \ \_____\
  \/_____/   \/_____/   \/_____/

Bob is a tool for creating flexible pieces of code from templates.

Usage:
  bob [command]

Available Commands:
  build       build a project from a specified template
  help        Help about any command
  init        initialize bob's config file

Flags:
      --config string   config file (default is $HOME/.bobconfig)
  -h, --help            help for bob

Use "bob [command] --help" for more information about a command.
```
