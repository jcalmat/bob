package io

import (
	"fmt"
	"strings"
)

func Ask(question string) string {
	fmt.Print(question)
	var answer string
	fmt.Scanln(&answer)
	return answer
}

func AskBool(question string) bool {
	fmt.Print(question)

	var answer string

	_, err := fmt.Scanln(&answer)
	if err != nil {
		return false
	}

	switch strings.ToLower(answer) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
		return AskBool(question)
	}
}
