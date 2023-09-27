package log

import (
	"fmt"
	"os"
)

func Err(msg string) {
	_, err := fmt.Fprint(os.Stderr, fmt.Sprintf("(!!!) %s", msg))
	if err != nil {
		fmt.Printf("Error writing to stderr: %v\n", err)
	}
}

func ErrLn(msg string) {
	Err(fmt.Sprintf("%s\n", msg))
}

func Warn(msg string) {
	_, err := fmt.Fprint(os.Stderr, fmt.Sprintf("(!) %s", msg))
	if err != nil {
		fmt.Printf("Error writing to stderr: %v\n", err)
	}
}

func WarnLn(msg string) {
	Warn(fmt.Sprintf("%s\n", msg))
}

func Out(msg string) {
	_, err := fmt.Fprint(os.Stdout, msg)
	if err != nil {
		fmt.Printf("Error writing to stdout: %v\n", err)
	}
}

func OutLn(msg string) {
	Out(fmt.Sprintf("%s\n", msg))
}
