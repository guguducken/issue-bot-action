package util

import "fmt"

var turn bool = true

func Info(message string) {
	if turn {
		fmt.Printf("\x1b[32m[INFO]: \x1b[0m%s\n", message)
	}
}

func Warning(message string) {
	if turn {
		fmt.Printf("\x1b[33m[WARNING]: \x1b[0m%s\n", message)
	}
}

func Error(message string) {
	if turn {
		fmt.Printf("\x1b[31m[ERROR]: \x1b[0m%s\n", message)
	}
}
