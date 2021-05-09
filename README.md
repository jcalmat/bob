# Bob - The Boilerplate Builder

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jcalmat/bob)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Linter](https://github.com/jcalmat/bob/workflows/golangci-lint/badge.svg)

## Introduction

Bob is a tool used to generate boilerplate code.

It was first made to avoid loosing too much time writing redundant code and allow developers to focus and dedicate more time on interesting parts of their projects.


The main idea is to write boilerplate once for all and let bob do the code generation for you.

Though, even if boilerplate represents redundant code, it could be slightly different from one file to another. Variables can change, blocks of code can exist in one but not in the other, this kind of things.
This is why Bob addresses this issue by using a particular syntax called `Actions` to perform some modifications from your template to your final code.

## Setup

### Binary

Download the latest [release](https://github.com/jcalmat/bob/releases), unzip and put the binary somewhere on your PATH.

### Building

#### Requirements

* golang 1.13.x or later

```bash
$> go get -u github.com/jcalmat/bob
$> bob
```

_OR_

```
$> git clone git@github.com:jcalmat/bob.git
$> make install
$> bob
```


## How to use it

### Configuration


Bob is asking you to do 2 things in order to work:
- Create template(s) of code
- Registering those templates to a config file written in yaml or json in your root directory. This file shall be named .bobconfig.yaml/.bobconfig.yml/.bobconfig.json

To generate your first config file, run `bob` and select `Init`.

### The config file

A bobconfig file looks like this:

```yaml
commands:
  test:
    description: "what this command will do"
    path: "/path/to/example.js"
    #OR
    git: "git@github.com/username/my-template.git"
    vars:
      - name: "my_variable"
        type: "string"
        ask: "Override the question Bob will ask in order to replace this value"

# If you need to clone your template from a private git repository
# please provide a ssh key here 
config:
  git:
    ssh:
      privateKeyFile: # your private ssh key, ex: "/home/user/.ssh/id_rsa" 
      privateKeyPassword: "" # password of this key if needed

```

bobconfig files are splitted into 1 required category + 1 optional one.

#### Commands category

Commands represents all the commands that will register to bob in order to clone a template. It contains the path (local or remote with git) of your template and all the `vars` you want to replace.

Actions can be a variable you want to replace, a boolean value representing the fact that a block of code will be in the final code or not, this kind of things.
Actions are blocks that contain the following fields:

1. name: the name of your action in the template
2. type: the type of your variable. In can be one of `string` or `bool`. Basically a string will be displayed as an input field and a bool will be displayed as a checkbox
3. format: allows you to override the question bob will ask to you in order to replace the action

Exemple:

The config file:
```yaml
vars:
  - name: "use_database"
    type: "bool"
    format: "Do you want to use a database?"
```
Your template:
```js
{{if .use_database}}
// init your db, perform operations, etc
{{end}}
```


##### Config category
Config allows you to input a ssh key (+ password if needed) to clone projects from a private repository.

## Templates

Bob uses the `go templates` syntax to parse and replace variables, here is what go template documentation can say about it:

```
The input text for a template is UTF-8-encoded text in any format.
"Actions"--data evaluations or control structures--are delimited by "{{" and "}}"; all text outside actions is copied to the output unchanged. Except for raw strings, actions may not span newlines, although comments can.
```

For more information about the format, here is a [cheat sheet](https://golang.org/pkg/text/template/#hdr-Actions).

_If it doesn't make a lot of sense at first glance, it will quickly, don't worry._


### The Basics

Given the file `example.js` with the following line of code:

```js
var {{.my_variable}}
```

Given the following `.bobconfig.yaml`:

```yaml
commands:
  test:
    path: "/path/to/example.js"
    actions:
      - name: "my_variable"
        type: "string"
```

When you run bob, it will automatically offer you to replace "my_variable" by any string and ask you where to put your newly created boiler file.

## Special variables

Since bob uses go template to perform the variable replacement, it has some interesting specificities.
You can, for instance, perform conditional operations.

Ex:
```go
{{if .my_variable}}
// do something only if .my_variable is defined
{{end}}
```

Bob also ships with homemade functions to ease string formatting
- `{{short .my_variable [X]}}` will only keep the X first characters of your variable
- `{{upcase .my_variable}}` will capitalize your variable
- `{{title .my_variable}}` will return a copy of the string s with all Unicode letters that begin words mapped to their Unicode title case

Ex:
```go
{{title .my_variable}}
// .my_variable = test -> Test

{{short .my_variable 1}}
// .my_variable = test -> t

{{upcase .my_variable}}
// .my_variable = test -> TEST
```

You can even combine multiple functions.

Ex:
```go
{{short .my_variable 3 | upcase}}
// my_variable = test -> TES
```

Once again, if you want to know what you can do, check out the [cheat sheet](https://golang.org/pkg/text/template/#hdr-Actions).


### Opinionated demo

Here is a small demo with one of my templates: creation of an entire compilable microservice.

![demo](./_examples/demo.gif)