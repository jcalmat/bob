package io

import "fmt"

func Title(title string) {
	fmt.Printf("\n==== %s ====\n\n", title)
}

func Info(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func AsciiBob() {
	fmt.Print(`
 ______     ______     ______
/\  == \   /\  __ \   /\  == \
\ \  __<   \ \ \/\ \  \ \  __<
 \ \_____\  \ \_____\  \ \_____\
  \/_____/   \/_____/   \/_____/
																 
`)
}
