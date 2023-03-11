package util

import "fmt"

func Info(message string) {
	fmt.Printf("\x1b[32m[INFO]: \x1b[0m%s\n", message)
}

func Warning(message string) {
	fmt.Printf("\x1b[43m[WARNING]: \x1b[0m%s\n", message)
}

func Error(message string) {
	fmt.Printf("\x1b[31m[ERROR]: \x1b[0m%s\n", message)
}
