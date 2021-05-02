console.log("Plain variable: {{.my_variable}}")
console.log("Capitalized variable: {{upcase .my_variable}}")
console.log("variable first character: {{short .my_variable 1}}")
console.log("Capitalized variable first character: {{short .my_variable 1 | upcase}}")
console.log("Titled variable: {{title .my_variable}}")
{{if .print_variable}}console.log("Conditional variable: {{.my_variable}}"){{end}}
