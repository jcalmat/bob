package main

import "fmt"

func main() {
	fmt.Printf("Plain variable: %s\n", {{.my_variable}})
	fmt.Printf("Capitalized variable: %s\n", {{upcase .my_variable}})
	fmt.Printf("variable first character: %s\n", {{short .my_variable 1}})
	fmt.Printf("Capitalized variable first character: %s\n", {{short .my_variable 1 | upcase}})
	fmt.Printf("Titled variable: %s\n", {{title .my_variable}})
	{{if .print_variable}}fmt.Printf("Conditional variable: %s\n", {{.my_variable}}){{end}}
}
